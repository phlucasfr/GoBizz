package repository

import (
	"auth-service/internal/domain"
	"auth-service/internal/logger"
	"auth-service/utils"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.uber.org/zap"
)

type CustomerRepository struct {
	db      *pgxpool.Pool
	redis   *redis.Client
	queries *Queries
}

// NewCustomerRepository creates a new instance of CustomerRepository.
// It initializes the repository with a PostgreSQL connection pool and a Redis client.
//
// Parameters:
//   - db: A pointer to a pgxpool.Pool instance representing the PostgreSQL connection pool.
//   - redis: A pointer to a redis.Client instance representing the Redis client.
//
// Returns:
//   - A pointer to a newly created CustomerRepository instance.
func NewCustomerRepository(db *pgxpool.Pool, redis *redis.Client) *CustomerRepository {
	return &CustomerRepository{
		db:      db,
		redis:   redis,
		queries: New(db),
	}
}

// Create creates a new customer in the repository. It performs the following steps:
// 1. Verifies if there is already an active customer with the same email or phone.
// 2. Hashes the customer's password for secure storage.
// 3. Encrypts the customer's CPF/CNPJ using the configured master key.
// 4. Constructs a new customer object and serializes it to JSON.
// 5. Stores the serialized customer data in Redis with a 10-minute expiration.
// 6. Sends a verification email to the customer's email address.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, deadlines, and cancellations.
//   - req: The CreateCustomerRequest containing the customer's details.
//
// Returns:
//   - A pointer to CreateCustomerResponse containing the created customer's ID and name.
//   - An error if any step in the process fails, including validation, hashing, encryption,
//     serialization, caching, or email sending.
func (r *CustomerRepository) Create(ctx context.Context, req domain.CreateCustomerRequest) (*domain.CreateCustomerResponse, error) {
	hasActiveCustomer, err := r.queries.HasActiveCustomer(ctx, HasActiveCustomerParams{
		Email: req.Email,
		Phone: req.Phone,
	})

	if err != nil {
		logger.Log.Error("error during verify customer data", zap.Error(err))
		return nil, fmt.Errorf("error during verify customer data: %v", err)
	}

	if hasActiveCustomer {
		logger.Log.Error("already has a customer with this email or phone")
		return nil, fmt.Errorf("already has a customer with this email or phone")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		logger.Log.Error("error during hash password", zap.Error(err))
		return nil, fmt.Errorf("error during hash password: %v", err)
	}

	encryptedCpfCnpj, err := utils.Encrypt(req.CPFCNPJ, utils.ConfigInstance.MasterKey)
	if err != nil {
		logger.Log.Error("error during encrypt cpfcnpj", zap.Error(err))
		return nil, fmt.Errorf("error during encrypt cpfcnpj: %v", err)
	}

	customer := domain.Customer{
		ID:             uuid.New(),
		Name:           req.Name,
		Email:          req.Email,
		Phone:          req.Phone,
		CpfCnpj:        encryptedCpfCnpj,
		HashedPassword: hashedPassword,
	}

	customerJSON, err := json.Marshal(&customer)
	if err != nil {
		logger.Log.Error("error during serialize customer data", zap.Error(err))
		return nil, fmt.Errorf("error during serialize customer data: %v", err)
	}

	cacheKey := fmt.Sprintf("customer:%x", customer.Email)
	err = r.redis.Set(ctx, cacheKey, customerJSON, 10*time.Minute).Err()
	if err != nil {
		logger.Log.Error("error during store customer in redis", zap.Error(err))
		return nil, fmt.Errorf("error during store customer in redis: %v", err)
	}

	err = r.SendVerificationEmail(ctx, customer.Email)
	if err != nil {
		logger.Log.Error("error during send verification email", zap.Error(err))
		return nil, fmt.Errorf("error during send verification email: %w", err)
	}

	logger.Log.Info("customer created successfully", zap.String("email", customer.Email))
	return &domain.CreateCustomerResponse{
		ID:   customer.ID,
		Name: customer.Name,
	}, nil
}

// ValidateEmailVerificationToken validates the email verification token for a given email address.
// It retrieves the token from Redis using the email as a key and compares it with the provided token.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, deadlines, and cancellations.
//   - email: The email address associated with the verification token.
//   - token: The token to be validated.
//
// Returns:
//   - An error if the token is invalid, expired, or if there is an issue accessing Redis.
//     Possible error scenarios include:
//   - Redis key not found (invalid or expired token).
//   - Redis access error.
//   - Deserialization error when parsing the Redis value.
//   - Token mismatch with the provided email.
func (r *CustomerRepository) ValidateEmailVerificationToken(ctx context.Context, email string, token string) error {
	cacheKey := fmt.Sprintf("email-verification:%s", email)

	cacheValueJSON, err := r.redis.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		logger.Log.Error("invalid or expired token", zap.String("email", email))
		return fmt.Errorf("invalid or expired token")
	} else if err != nil {
		logger.Log.Error("error during access redis", zap.Error(err))
		return fmt.Errorf("error during access redis: %v", err)
	}

	var cacheValue map[string]string
	err = json.Unmarshal([]byte(cacheValueJSON), &cacheValue)
	if err != nil {
		logger.Log.Error("error during deserialize data from redis", zap.Error(err))
		return fmt.Errorf("error during deserialize data from redis: %v", err)
	}

	storedToken, ok := cacheValue["token"]
	if !ok || storedToken != token {
		logger.Log.Error("invalid token", zap.String("email", email))
		return fmt.Errorf("token does not match the provided email")
	}

	logger.Log.Info("token validated successfully", zap.String("email", email))
	return nil
}

// ActivateCustomerByEmail activates a customer account based on their email address.
// It performs the following steps:
// 1. Attempts to retrieve the customer data from Redis cache using the email as a key.
// 2. If the customer data is found in the cache, it deserializes the data into a Customer object.
// 3. Persists the customer data into the database.
// 4. Activates the customer in the database using the provided email.
// 5. Deletes the email verification token from Redis.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, deadlines, and cancellations.
//   - email: The email address of the customer to be activated.
//
// Returns:
//   - error: An error if any step in the activation process fails, otherwise nil.
func (r *CustomerRepository) ActivateCustomerByEmail(ctx context.Context, email string) error {
	var customer *domain.Customer
	cacheKey := fmt.Sprintf("customer:%x", email)

	cachedCustomer, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil && cachedCustomer != "" {
		err := json.Unmarshal([]byte(cachedCustomer), &customer)
		if err != nil {
			logger.Log.Error("error during deserialize data from redis", zap.Error(err))
			return fmt.Errorf("error during deserialize data from redis: %v", err)
		}
	}

	err = r.persistCustomerInDB(ctx, customer)
	if err != nil {
		logger.Log.Error("error during save customer in database", zap.Error(err))
		return fmt.Errorf("error during save customer in database: %v", err)
	}

	_, err = r.queries.ActivateCustomerByEmail(ctx, email)
	if err != nil {
		logger.Log.Error("error during activate customer in database", zap.Error(err))
		return fmt.Errorf("error during activate customer in database: %v", err)
	}

	tokenKey := fmt.Sprintf("email-verification:%s", email)
	err = r.redis.Del(ctx, tokenKey).Err()
	if err != nil {
		logger.Log.Error("error during remove email token in redis", zap.Error(err))
		return fmt.Errorf("error during remove email token in redis: %v", err)
	}

	logger.Log.Info("customer activated successfully", zap.String("email", email))
	return nil
}

// SendVerificationEmail sends a verification email to the specified email address.
// It generates a unique token, stores it in Redis with a 10-minute expiration, and sends
// an email containing a verification link to the user.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, deadlines, and cancellation signals.
//   - email: The email address to which the verification email will be sent.
//
// Returns:
//   - error: An error if any issue occurs during the process, including Redis operations,
//     token generation, or email sending.
//
// The function performs the following steps:
//  1. Checks if a verification token already exists for the email in Redis.
//  2. Deletes the old token if it exists.
//  3. Generates a new verification token and stores it in Redis.
//  4. Constructs a verification link using the frontend source URL and the token.
//  5. Sends the verification email using the SendGrid API.
//
// Note:
//   - The frontend source URL must be configured in the application settings.
//   - The SendGrid API key must be set in the application configuration.
func (r *CustomerRepository) SendVerificationEmail(ctx context.Context, email string) error {
	emailKey := fmt.Sprintf("email-verification:%s", email)

	existingToken, err := r.redis.Get(ctx, emailKey).Result()
	if err != nil && err != redis.Nil {
		logger.Log.Error("error while checking email in redis", zap.Error(err))
		return fmt.Errorf("error while checking email in redis: %v", err)
	}

	if existingToken != "" {
		err = r.redis.Del(ctx, emailKey).Err()
		if err != nil {
			logger.Log.Error("error while removing old email token from redis", zap.Error(err))
			return fmt.Errorf("error while removing old email token from redis: %v", err)
		}
	}

	token := utils.GenerateResetToken(email)
	cacheValue := map[string]string{"email": email, "token": token}

	cacheValueJSON, err := json.Marshal(cacheValue)
	if err != nil {
		logger.Log.Error("error while serializing data to redis", zap.Error(err))
		return fmt.Errorf("error while serializing data to redis: %v", err)
	}

	err = r.redis.Set(ctx, emailKey, cacheValueJSON, 10*time.Minute).Err()
	if err != nil {
		logger.Log.Error("error while storing new token in redis", zap.Error(err))
		return fmt.Errorf("error while storing new token in redis: %v", err)
	}

	frontendSource := utils.ConfigInstance.FrontendSource
	if frontendSource == "" {
		logger.Log.Error("frontend source is not set")
		return fmt.Errorf("frontend source is not set")
	}

	verificationLink := fmt.Sprintf("%s/email-verification?token=%s", frontendSource, token)

	from := mail.NewEmail("GoBizz", "gobizz.comercial@gmail.com")
	subject := "Email Verification"
	to := mail.NewEmail("User", email)
	plainTextContent := fmt.Sprintf("Click the link to verify your email: %s", verificationLink)
	htmlContent := fmt.Sprintf(`<p>Click the link to verify your email: <a href="%s">%s</a></p>`, verificationLink, verificationLink)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(utils.ConfigInstance.SendGridApiKey)
	_, err = client.Send(message)
	if err != nil {
		logger.Log.Error("error during sending verification email", zap.Error(err))
		return fmt.Errorf("error during sending verification email: %v", err)
	}

	logger.Log.Info("verification email sent successfully", zap.String("email", email))
	return nil
}

// persistCustomerInDB persists a customer entity into the database and clears the corresponding cache entry in Redis.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, deadlines, and cancellation signals.
//   - customer: A pointer to the domain.Customer object containing the customer details to be persisted.
//
// Returns:
//   - error: An error if the operation fails, otherwise nil.
//
// This function performs the following steps:
//  1. Maps the customer details to the CreateCustomerParams struct and saves the customer in the database.
//  2. Deletes the customer's cache entry in Redis using the customer's email as the cache key.
//
// Errors:
//   - Returns an error if the database operation fails.
//   - Returns an error if the Redis cache deletion fails.
func (r *CustomerRepository) persistCustomerInDB(ctx context.Context, customer *domain.Customer) error {
	params := CreateCustomerParams{
		Name:           customer.Name,
		Email:          customer.Email,
		Phone:          customer.Phone,
		CpfCnpj:        customer.CpfCnpj,
		IsActive:       true,
		HashedPassword: customer.HashedPassword,
	}

	_, err := r.queries.CreateCustomer(ctx, params)
	if err != nil {
		logger.Log.Error("error during save customer in database", zap.Error(err))
		return fmt.Errorf("error during save customer in database: %w", err)
	}

	cacheKey := fmt.Sprintf("customer:%x", customer.Email)
	err = r.redis.Del(ctx, cacheKey).Err()
	if err != nil {
		logger.Log.Error("error during remove customer from redis after persistence", zap.Error(err))
		return fmt.Errorf("error during remove customer from redis after persistence: %v", err)
	}

	logger.Log.Info("customer persisted in database successfully", zap.String("email", customer.Email))
	return nil
}

// GetCustomerByEmail retrieves a customer from the database based on the provided email address.
// It queries the database using the provided context and email, and returns a domain.Customer object
// if a matching record is found. If an error occurs during the query, it returns nil and the error.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, deadlines, and cancelation signals.
//   - email: The email address of the customer to retrieve.
//
// Returns:
//   - *domain.Customer: A pointer to the customer object containing the retrieved customer details.
//   - error: An error object if the query fails or no customer is found.
func (r *CustomerRepository) GetCustomerByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	customer, err := r.queries.GetCustomerByEmail(ctx, email)

	if err != nil {
		if err == pgx.ErrNoRows {
			logger.Log.Error("customer not found", zap.String("email", email))
			return nil, fmt.Errorf("customer not found")
		}
		return nil, err
	}

	logger.Log.Info("customer retrieved successfully", zap.String("email", email))
	return &domain.Customer{
		ID:             customer.ID.Bytes,
		Name:           customer.Name,
		Email:          customer.Email,
		Phone:          customer.Phone,
		CpfCnpj:        customer.CpfCnpj,
		IsActive:       customer.IsActive,
		UpdatedAt:      customer.UpdatedAt,
		CreatedAt:      customer.CreatedAt,
		HashedPassword: customer.HashedPassword,
	}, nil
}

// SendRecoveryEmail sends a password recovery email to the specified email address.
// It generates a unique token, stores it in Redis with a 10-minute expiration, and sends
// an email containing a reset link to the user.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, deadlines, and cancellations.
//   - email: The email address of the user requesting password recovery.
//
// Returns:
//   - error: An error if any issue occurs during the process, such as Redis operations,
//     token generation, or email sending.
//
// Process:
//  1. Checks if an existing recovery token for the email exists in Redis.
//     If found, it deletes the old token.
//  2. Generates a new recovery token and stores it in Redis with a 10-minute expiration.
//  3. Constructs a password reset link using the frontend source URL and the generated token.
//  4. Sends an email to the user with the reset link using SendGrid.
//
// Errors:
//   - Returns an error if Redis operations fail (e.g., getting, deleting, or setting keys).
//   - Returns an error if the frontend source URL or SendGrid API key is not configured.
//   - Returns an error if the email fails to send via SendGrid.
func (r *CustomerRepository) SendRecoveryEmail(ctx context.Context, email string) error {
	existingTokenKey := fmt.Sprintf("password-recovery:%s", email)

	existingToken, err := r.redis.Get(ctx, existingTokenKey).Result()
	if err != nil && err != redis.Nil {
		logger.Log.Error("error while checking existing token in Redis", zap.Error(err))
		return fmt.Errorf("error while checking existing token in Redis: %v", err)
	}

	if existingToken != "" {
		err = r.redis.Del(ctx, existingTokenKey).Err()
		if err != nil {
			logger.Log.Error("error while removing old token from Redis", zap.Error(err))
			return fmt.Errorf("error while removing old token from Redis: %v", err)
		}
	}

	token := utils.GenerateResetToken(email)
	tokenKey := fmt.Sprintf("password-recovery:%s", email)
	cacheValue := map[string]string{"email": email, "token": token}

	cacheValueJSON, err := json.Marshal(cacheValue)
	if err != nil {
		logger.Log.Error("error while serializing data to Redis", zap.Error(err))
		return fmt.Errorf("error while serializing data to Redis: %v", err)
	}

	err = r.redis.Set(ctx, tokenKey, cacheValueJSON, 10*time.Minute).Err()
	if err != nil {
		logger.Log.Error("error while storing token in Redis", zap.Error(err))
		return fmt.Errorf("error while storing token in Redis: %v", err)
	}

	frontendSource := utils.ConfigInstance.FrontendSource
	if frontendSource == "" {
		logger.Log.Error("frontend source is not set")
		return fmt.Errorf("frontend source is not set")
	}

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", frontendSource, token)

	from := mail.NewEmail("GoBizz", "gobizz.comercial@gmail.com")
	subject := "Password Recovery"
	to := mail.NewEmail("User", email)
	plainTextContent := fmt.Sprintf("Click the link to reset your password: %s", resetLink)
	htmlContent := fmt.Sprintf(`<p>Click the link to reset your password: <a href="%s">%s</a></p>`, resetLink, resetLink)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	sendGridApiKey := utils.ConfigInstance.SendGridApiKey
	if sendGridApiKey == "" {
		logger.Log.Error("sendgrid api key is not set")
		return fmt.Errorf("sendgrid api key is not set")
	}

	client := sendgrid.NewSendClient(sendGridApiKey)
	_, err = client.Send(message)
	if err != nil {
		logger.Log.Error("error during sending recovery email", zap.Error(err))
		return fmt.Errorf("error during sending recovery email: %v", err)
	}

	logger.Log.Info("password recovery email sent successfully", zap.String("email", email))
	return nil
}

// ValidateResetToken validates a password reset token by checking its format,
// verifying it against stored data in Redis, and ensuring it matches the expected token.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, deadlines, and cancellations.
//   - token: The reset token to be validated, expected in the format "prefix:email".
//
// Returns:
//   - string: The email associated with the valid token, if validation succeeds.
//   - error: An error if the token is invalid, expired, or if there is an issue accessing Redis.
//
// Validation Steps:
//  1. Splits the token into parts and checks its format.
//  2. Retrieves the token data from Redis using the email as part of the cache key.
//  3. Deserializes the cached data and verifies the token and email match the stored values.
//
// Errors:
//   - Returns an error if the token format is invalid.
//   - Returns an error if the token is not found or has expired in Redis.
//   - Returns an error if there is an issue accessing or deserializing data from Redis.
//   - Returns an error if the token or email does not match the stored values.
func (r *CustomerRepository) ValidateResetToken(ctx context.Context, token string) (string, error) {
	tokenParts := strings.Split(token, ":")
	if len(tokenParts) != 2 {
		logger.Log.Error("invalid token format", zap.String("token", token))
		return "", fmt.Errorf("invalid token format")
	}
	email := tokenParts[1]
	cacheKey := fmt.Sprintf("password-recovery:%s", email)

	cacheValueJSON, err := r.redis.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		logger.Log.Error("invalid or expired token", zap.String("email", email))
		return "", fmt.Errorf("invalid or expired token")
	} else if err != nil {
		logger.Log.Error("error during access redis", zap.Error(err))
		return "", fmt.Errorf("error during access redis: %v", err)
	}

	var cacheValue map[string]string
	err = json.Unmarshal([]byte(cacheValueJSON), &cacheValue)
	if err != nil {
		logger.Log.Error("error during deserialize data from redis", zap.Error(err))
		return "", fmt.Errorf("error during deserialize data from redis: %v", err)
	}

	storedToken, ok := cacheValue["token"]
	if !ok || storedToken != token {
		logger.Log.Error("invalid token", zap.String("token", token))
		return "", fmt.Errorf("invalid token")
	}

	email, ok = cacheValue["email"]
	if !ok {
		logger.Log.Error("email not found for the provided token", zap.String("token", token))
		return "", fmt.Errorf("email not found for the provided token")
	}

	logger.Log.Info("token validated successfully", zap.String("email", email))
	return email, nil
}

// UpdatePasswordByEmail updates the password of a customer in the database
// based on their email address. It also removes the associated password
// recovery token from Redis.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, deadlines, and cancellations.
//   - email: The email address of the customer whose password is being updated.
//   - hashedPassword: The new hashed password to be stored in the database.
//
// Returns:
//   - error: An error if the update operation in the database or the token removal
//     from Redis fails; otherwise, nil.
func (r *CustomerRepository) UpdatePasswordByEmail(ctx context.Context, email string, hashedPassword string) error {
	_, err := r.queries.UpdatePasswordByEmail(ctx, UpdatePasswordByEmailParams{
		Email:          email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		logger.Log.Error("error during update password in database", zap.Error(err))
		return fmt.Errorf("error during update password in database: %v", err)
	}

	tokenKey := fmt.Sprintf("password-recovery:%s", email)
	err = r.redis.Del(ctx, tokenKey).Err()
	if err != nil {
		logger.Log.Error("error during remove token from redis", zap.Error(err))
		return fmt.Errorf("error during remove token from redis: %v", err)
	}

	logger.Log.Info("password updated successfully", zap.String("email", email))
	return nil
}

// BlacklistToken adds a token to the blacklist in the Redis store with a specified time-to-live (TTL).
// This prevents the token from being used for authentication or other purposes.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, deadlines, and cancellations.
//   - tokenKey: The key representing the token to be blacklisted.
//   - ttl: The duration for which the token should remain blacklisted.
//
// Returns:
//   - error: An error if the operation fails, or nil if the token is successfully blacklisted.
func (r *CustomerRepository) BlacklistToken(ctx context.Context, tokenKey string, ttl time.Duration) error {
	logger.Log.Info("blacklisting token", zap.String("tokenKey", tokenKey))
	return r.redis.Set(ctx, tokenKey, "blacklisted", ttl).Err()
}
