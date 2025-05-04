package handlers

import (
	"auth-service/internal/domain"
	"auth-service/internal/infra/repository"
	"auth-service/internal/logger"
	"auth-service/utils"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type CustomerHandler struct {
	repo *repository.CustomerRepository
}

// NewCustomerHandler creates a new instance of CustomerHandler with the provided
// CustomerRepository. It initializes the handler with the given repository to
// manage customer-related operations.
//
// Parameters:
//   - repo: A pointer to CustomerRepository that provides access to customer data.
//
// Returns:
//   - A pointer to a newly created CustomerHandler.
func NewCustomerHandler(repo *repository.CustomerRepository) *CustomerHandler {
	return &CustomerHandler{
		repo: repo,
	}
}

// Create handles the creation of a new customer.
// It parses the incoming request body into a CreateCustomerRequest struct,
// validates the payload, and invokes the repository layer to create the customer.
// If the request payload is invalid, it returns a 400 Bad Request response.
// If an error occurs during the creation process, it returns a 500 Internal Server Error response.
// On success, it returns a 201 Created response with the created customer data.
//
// @Summary Create a new customer
// @Description Handles the creation of a new customer
// @Tags Customers
// @Accept json
// @Produce json
// @Param request body domain.CreateCustomerRequest true "Customer creation payload"
// @Success 201 {object} domain.CustomerResponse
// @Failure 400 {object} fiber.Map "Invalid request payload"
// @Failure 500 {object} fiber.Map "Internal server error"
// @Router /customers [post]
func (h *CustomerHandler) Create(c *fiber.Ctx) error {
	var req domain.CreateCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Log.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	response, err := h.repo.Create(c.Context(), req)
	if err != nil {
		logger.Log.Error("Failed to create customer", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	logger.Log.Info("Customer created successfully", zap.String("customer_id", response.ID.String()))
	return c.Status(fiber.StatusCreated).JSON(response)
}

// RecoverPassword handles the password recovery process for a customer.
//
// This function parses the incoming request to extract the email address,
// checks if a customer with the provided email exists in the database,
// and sends a password recovery email if the customer is found.
//
// Parameters:
//   - c (*fiber.Ctx): The Fiber context containing the HTTP request and response.
//
// Returns:
//   - error: An error if the operation fails, or nil if successful.
//
// Possible Responses:
//   - 400 Bad Request: If the request payload is invalid or if sending the recovery email fails.
//   - 404 Not Found: If no customer is found with the provided email address.
//   - 500 Internal Server Error: If an unexpected error occurs during database operations.
//   - 200 OK: If the recovery email is sent successfully.
func (h *CustomerHandler) RecoverPassword(c *fiber.Ctx) error {
	var req *domain.PasswordRecoveryRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Log.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	_, err := h.repo.GetCustomerByEmail(c.Context(), req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Error("Customer not found", zap.String("email", req.Email))
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "customer not found",
			})
		}

		logger.Log.Error("Failed to get customer by email", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := h.repo.SendRecoveryEmail(c.Context(), req.Email); err != nil {
		logger.Log.Error("Failed to send recovery email", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	logger.Log.Info("Recovery email sent successfully", zap.String("email", req.Email))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "recovery email sent successfully!",
	})
}

// VerifyCustomerByEmail handles the verification of a customer's email address.
// It expects a request payload containing a token in the format "prefix:email".
// The function performs the following steps:
// 1. Parses the request body into a VerifyCustomerByEmailRequest structure.
// 2. Validates the token format to ensure it contains two parts separated by a colon.
// 3. Extracts the email from the token and validates the email verification token using the repository.
// 4. Activates the customer associated with the email if the token is valid.
// 5. Returns appropriate HTTP responses based on the success or failure of the operations.
//
// Possible HTTP responses:
// - 400 Bad Request: If the request payload is invalid or the token format is incorrect.
// - 401 Unauthorized: If the email verification token is invalid.
// - 500 Internal Server Error: If there is an error activating the customer.
// - 200 OK: If the customer is successfully activated.
//
// Parameters:
// - c: The Fiber context containing the HTTP request and response.
//
// Returns:
// - An HTTP response with the appropriate status code and message.
func (h *CustomerHandler) VerifyCustomerByEmail(c *fiber.Ctx) error {
	var req *domain.VerifyCustomerByEmailRequest

	if err := c.BodyParser(&req); err != nil {
		logger.Log.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	tokenParts := strings.Split(req.Token, ":")
	if len(tokenParts) != 2 {
		logger.Log.Error("Invalid token format", zap.String("token", req.Token))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid token format",
		})
	}

	email := tokenParts[1]

	err := h.repo.ValidateEmailVerificationToken(c.Context(), email, req.Token)
	if err != nil {
		logger.Log.Error("Failed to validate email verification token", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := h.repo.ActivateCustomerByEmail(c.Context(), email); err != nil {
		logger.Log.Error("Failed to activate customer", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	logger.Log.Info("Customer activated successfully", zap.String("email", email))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Customer activated successfully.",
	})
}

// ResetPassword handles the password reset process for a customer.
//
// This function performs the following steps:
// 1. Parses the incoming request body to extract the ResetPasswordRequest payload.
// 2. Validates the reset token provided in the request by calling the repository's ValidateResetToken method.
// 3. Hashes the new password provided in the request.
// 4. Updates the customer's password in the database using the repository's UpdatePasswordByEmail method.
// 5. Returns appropriate HTTP responses based on the success or failure of the operations.
//
// Parameters:
// - c: The Fiber context containing the HTTP request and response.
//
// Returns:
// - An error if any step in the process fails, along with an appropriate HTTP status code and error message.
// - A success message with HTTP status 200 if the password reset is completed successfully.
func (h *CustomerHandler) ResetPassword(c *fiber.Ctx) error {
	var req *domain.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Log.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if err := c.BodyParser(&req); err != nil {
		logger.Log.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	email, err := h.repo.ValidateResetToken(c.Context(), req.Token)
	if err != nil {
		logger.Log.Error("Failed to validate reset token", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	hashedPassword, _ := utils.HashPassword(req.Password)
	if err := h.repo.UpdatePasswordByEmail(c.Context(), email, hashedPassword); err != nil {
		logger.Log.Error("Failed to update password", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	logger.Log.Info("Password reset successfully", zap.String("email", email))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password reset successfully.",
	})
}

// Login handles the customer login process.
//
// @Summary      Customer Login
// @Description  Authenticates a customer using their email and password, and returns a JWT token upon successful login.
// @Tags         Customer
// @Accept       json
// @Produce      json
// @Param        request body domain.LoginRequest true "Login Request"
// @Success      200 {object} map[string]interface{} "Returns user details and JWT token"
// @Failure      400 {object} map[string]interface{} "Invalid request payload"
// @Failure      401 {object} map[string]interface{} "Invalid credentials or account not activated"
// @Failure      500 {object} map[string]interface{} "Failed to generate token"
// @Router       /login [post]
func (h *CustomerHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Log.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	customer, err := h.repo.GetCustomerByEmail(c.Context(), req.Email)
	if err != nil {
		logger.Log.Error("Failed to get customer by email", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	if !customer.IsActive {
		logger.Log.Error("Account not activated", zap.String("email", req.Email))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Account not activated",
		})
	}

	if err := utils.CheckPassword(req.Password, customer.HashedPassword); err != nil {
		logger.Log.Error("Invalid credentials", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	token, err := utils.GenerateJWT(customer.ID.String(), customer.Email, c.IP(), c.Get("User-Agent"))
	if err != nil {
		logger.Log.Error("Failed to generate token", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	logger.Log.Info("Customer logged in successfully", zap.String("customer_id", customer.ID.String()))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": fiber.Map{
			"id":    customer.ID,
			"name":  customer.Name,
			"email": customer.Email,
		},
		"token": token,
	})
}

// RefreshToken handles the token refresh process for a customer.
// It parses the incoming request to extract the refresh token payload,
// validates the provided JWT token, and generates a new token if the
// validation is successful.
//
// @param c *fiber.Ctx - The Fiber context containing the HTTP request and response.
//
// @return error - Returns an error if the request payload is invalid, the token
//
//	is missing or invalid, or if there is an issue generating a new token.
//
// Response Codes:
// - 400 Bad Request: If the request payload is invalid.
// - 401 Unauthorized: If the token is missing or invalid.
// - 500 Internal Server Error: If there is an issue generating a new token.
// - 200 OK: If the token is successfully refreshed, returns the new token in the response body.
func (h *CustomerHandler) RefreshToken(c *fiber.Ctx) error {
	var req domain.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Log.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	token := c.Get("Authorization")
	if token == "" {
		logger.Log.Error("No token provided")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No token provided",
		})
	}

	token = strings.TrimPrefix(token, "Bearer ")

	_, err := utils.ValidateJWT(token)
	if err != nil {
		logger.Log.Error("Invalid token", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	newToken, err := utils.GenerateJWT(req.ID.String(), req.Email, c.IP(), c.Get("User-Agent"))
	if err != nil {
		logger.Log.Error("Failed to generate new token", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to refresh token",
		})
	}

	logger.Log.Info("Token refreshed successfully", zap.String("customer_id", req.ID.String()))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Token refreshed successfully",
		"token":   newToken,
	})
}

// Logout handles the logout process for a customer by blacklisting their token.
// It retrieves the "Authorization" header from the request, validates its presence,
// and removes the "Bearer " prefix if present. The token is then blacklisted by
// storing it in the repository with a specified expiration time.
//
// Parameters:
//   - c: The Fiber context containing the HTTP request and response.
//
// Returns:
//   - An error if the token is not provided or if there is an issue blacklisting the token.
//   - A JSON response with a success message if the logout process is completed successfully.
//
// Response Codes:
//   - 401 Unauthorized: If the "Authorization" header is missing.
//   - 200 OK: If the token is successfully blacklisted and the user is logged out.
func (h *CustomerHandler) Logout(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		logger.Log.Error("No token provided")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No token provided",
		})
	}

	token = strings.TrimPrefix(token, "Bearer ")

	tokenKey := fmt.Sprintf("blacklisted-token:%s", token)
	err := h.repo.BlacklistToken(c.Context(), tokenKey, 24*time.Hour*3)
	if err != nil {
		logger.Log.Error("Failed to blacklist token", zap.Error(err))
	}

	logger.Log.Info("Customer logged out successfully", zap.String("token", token))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Logged out successfully",
	})
}

// ValidateSession validates the user's session by checking the Authorization token
// provided in the request header. It ensures the token is present, trims the "Bearer"
// prefix, and validates the token using the utility function ValidateJWT. If the token
// is invalid or missing, it returns an unauthorized status with an appropriate error
// message. If the token is valid, it responds with a success status and the user's
// details extracted from the token claims.
//
// Parameters:
//   - c: The Fiber context containing the HTTP request and response.
//
// Returns:
//   - An error if the token is invalid or missing, or nil if the session is valid.
func (h *CustomerHandler) ValidateSession(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		logger.Log.Error("No token provided")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No token provided",
		})
	}

	token = strings.TrimPrefix(token, "Bearer ")

	claims, err := utils.ValidateJWT(token)
	if err != nil {
		logger.Log.Error("Invalid token", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	logger.Log.Info("Session validated successfully", zap.String("user_id", claims.UserID))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"isValid": true,
		"user": fiber.Map{
			"id":    claims.UserID,
			"email": claims.Email,
		},
	})
}
