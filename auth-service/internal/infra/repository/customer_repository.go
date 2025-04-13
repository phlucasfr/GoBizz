package repository

import (
	"auth-service/internal/domain"
	"auth-service/utils"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type CustomerRepository struct {
	db      *pgxpool.Pool
	redis   *redis.Client
	queries *Queries
}

func NewCustomerRepository(db *pgxpool.Pool, redis *redis.Client) *CustomerRepository {
	return &CustomerRepository{
		db:      db,
		redis:   redis,
		queries: New(db),
	}
}

func (r *CustomerRepository) Create(ctx context.Context, req domain.CreateCustomerRequest) (*domain.CreateCustomerResponse, error) {
	hasActiveCustomer, err := r.queries.HasActiveCustomer(ctx, HasActiveCustomerParams{
		Email: req.Email,
		Phone: req.Phone,
	})

	if err != nil {
		return nil, fmt.Errorf("error during verify customer data: %v", err)
	}

	if hasActiveCustomer {
		return nil, fmt.Errorf("already has a customer with this email or phone")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("error during hash password: %v", err)
	}

	encryptedCpfCnpj, err := utils.Encrypt(req.CPFCNPJ, utils.ConfigInstance.MasterKey)
	if err != nil {
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
		return nil, fmt.Errorf("error during serialize customer data: %v", err)
	}

	cacheKey := fmt.Sprintf("customer:%x", customer.Email)
	err = r.redis.Set(ctx, cacheKey, customerJSON, 10*time.Minute).Err()
	if err != nil {
		return nil, fmt.Errorf("error during store customer in redis: %v", err)
	}

	err = r.SendVerificationEmail(ctx, customer.Email)
	if err != nil {
		return nil, fmt.Errorf("error during send verification email: %w", err)
	}

	return &domain.CreateCustomerResponse{
		ID:   customer.ID,
		Name: customer.Name,
	}, nil
}

func (r *CustomerRepository) ValidateEmailVerificationToken(ctx context.Context, email string, token string) error {
	cacheKey := fmt.Sprintf("email-verification:%s", email)

	cacheValueJSON, err := r.redis.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return fmt.Errorf("invalid or expired token")
	} else if err != nil {
		return fmt.Errorf("error during access redis: %v", err)
	}

	var cacheValue map[string]string
	err = json.Unmarshal([]byte(cacheValueJSON), &cacheValue)
	if err != nil {
		return fmt.Errorf("error during deserialize data from redis: %v", err)
	}

	storedToken, ok := cacheValue["token"]
	if !ok || storedToken != token {
		return fmt.Errorf("token does not match the provided email")
	}

	return nil
}

func (r *CustomerRepository) ActivateCustomerByEmail(ctx context.Context, email string) error {
	var customer *domain.Customer
	cacheKey := fmt.Sprintf("customer:%x", email)

	cachedCustomer, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil && cachedCustomer != "" {
		err := json.Unmarshal([]byte(cachedCustomer), &customer)
		if err != nil {
			return fmt.Errorf("error during deserialize data from redis: %v", err)
		}
	}

	err = r.persistCustomerInDB(ctx, customer)
	if err != nil {
		return fmt.Errorf("error during save customer in database: %v", err)
	}

	_, err = r.queries.ActivateCustomerByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("error during activate customer in database: %v", err)
	}

	tokenKey := fmt.Sprintf("email-verification:%s", email)
	err = r.redis.Del(ctx, tokenKey).Err()
	if err != nil {
		return fmt.Errorf("error during remove email token in redis: %v", err)
	}

	return nil
}

func (r *CustomerRepository) SendVerificationEmail(ctx context.Context, email string) error {
	emailKey := fmt.Sprintf("email-verification:%s", email)

	existingToken, err := r.redis.Get(ctx, emailKey).Result()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("error while checking email in redis: %v", err)
	}

	if existingToken != "" {
		err = r.redis.Del(ctx, emailKey).Err()
		if err != nil {
			return fmt.Errorf("error while removing old email token from redis: %v", err)
		}
	}

	token := utils.GenerateResetToken(email)
	cacheValue := map[string]string{"email": email, "token": token}

	cacheValueJSON, err := json.Marshal(cacheValue)
	if err != nil {
		return fmt.Errorf("error while serializing data to redis: %v", err)
	}

	err = r.redis.Set(ctx, emailKey, cacheValueJSON, 10*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("error while storing new token in redis: %v", err)
	}

	frontendSource := utils.ConfigInstance.FrontendSource
	if frontendSource == "" {
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
		return fmt.Errorf("error during sending verification email: %v", err)
	}

	return nil
}

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
		return fmt.Errorf("error during save customer in database: %w", err)
	}

	cacheKey := fmt.Sprintf("customer:%x", customer.Email)
	err = r.redis.Del(ctx, cacheKey).Err()
	if err != nil {
		return fmt.Errorf("error during remove customer from redis after persistence: %v", err)
	}

	return nil
}

func (r *CustomerRepository) GetCustomerByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	customer, err := r.queries.GetCustomerByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

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

func (r *CustomerRepository) SendRecoveryEmail(ctx context.Context, email string) error {
	existingTokenKey := fmt.Sprintf("password-recovery:%s", email)

	existingToken, err := r.redis.Get(ctx, existingTokenKey).Result()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("error while checking existing token in Redis: %v", err)
	}

	if existingToken != "" {
		err = r.redis.Del(ctx, existingTokenKey).Err()
		if err != nil {
			return fmt.Errorf("error while removing old token from Redis: %v", err)
		}
	}

	token := utils.GenerateResetToken(email)
	tokenKey := fmt.Sprintf("password-recovery:%s", email)
	cacheValue := map[string]string{"email": email, "token": token}

	cacheValueJSON, err := json.Marshal(cacheValue)
	if err != nil {
		return fmt.Errorf("error while serializing data to Redis: %v", err)
	}

	err = r.redis.Set(ctx, tokenKey, cacheValueJSON, 10*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("error while storing token in Redis: %v", err)
	}

	frontendSource := utils.ConfigInstance.FrontendSource
	if frontendSource == "" {
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
		return fmt.Errorf("sendgrid api key is not set")
	}

	client := sendgrid.NewSendClient(sendGridApiKey)
	_, err = client.Send(message)
	if err != nil {
		return fmt.Errorf("error during sending recovery email: %v", err)
	}

	return nil
}

func (r *CustomerRepository) ValidateResetToken(ctx context.Context, token string) (string, error) {
	tokenParts := strings.Split(token, ":")
	if len(tokenParts) != 2 {
		return "", fmt.Errorf("invalid token format")
	}
	email := tokenParts[1]
	cacheKey := fmt.Sprintf("password-recovery:%s", email)

	cacheValueJSON, err := r.redis.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("invalid or expired token")
	} else if err != nil {
		return "", fmt.Errorf("error during access redis: %v", err)
	}

	var cacheValue map[string]string
	err = json.Unmarshal([]byte(cacheValueJSON), &cacheValue)
	if err != nil {
		return "", fmt.Errorf("error during deserialize data from redis: %v", err)
	}

	storedToken, ok := cacheValue["token"]
	if !ok || storedToken != token {
		return "", fmt.Errorf("invalid token")
	}

	email, ok = cacheValue["email"]
	if !ok {
		return "", fmt.Errorf("email not found for the provided token")
	}

	return email, nil
}

func (r *CustomerRepository) UpdatePasswordByEmail(ctx context.Context, email string, hashedPassword string) error {
	_, err := r.queries.UpdatePasswordByEmail(ctx, UpdatePasswordByEmailParams{
		Email:          email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		return fmt.Errorf("error during update password in database: %v", err)
	}

	tokenKey := fmt.Sprintf("password-recovery:%s", email)
	err = r.redis.Del(ctx, tokenKey).Err()
	if err != nil {
		return fmt.Errorf("error during remove token from redis: %v", err)
	}

	return nil
}

func (r *CustomerRepository) BlacklistToken(ctx context.Context, tokenKey string, ttl time.Duration) error {
	return r.redis.Set(ctx, tokenKey, "blacklisted", ttl).Err()
}
