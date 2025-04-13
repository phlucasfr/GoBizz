package handlers

import (
	"auth-service/internal/domain"
	"auth-service/internal/infra/repository"
	"auth-service/utils"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type CustomerHandler struct {
	repo *repository.CustomerRepository
}

func NewCustomerHandler(repo *repository.CustomerRepository) *CustomerHandler {
	return &CustomerHandler{
		repo: repo,
	}
}

func (h *CustomerHandler) Create(c *fiber.Ctx) error {
	var req domain.CreateCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	response, err := h.repo.Create(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *CustomerHandler) RecoverPassword(c *fiber.Ctx) error {
	var req *domain.PasswordRecoveryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	_, err := h.repo.GetCustomerByEmail(c.Context(), req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "customer not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := h.repo.SendRecoveryEmail(c.Context(), req.Email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "recovery email sent successfully!",
	})
}

func (h *CustomerHandler) VerifyCustomerByEmail(c *fiber.Ctx) error {
	var req *domain.VerifyCustomerByEmailRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	tokenParts := strings.Split(req.Token, ":")
	if len(tokenParts) != 2 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid token format",
		})
	}

	email := tokenParts[1]

	err := h.repo.ValidateEmailVerificationToken(c.Context(), email, req.Token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := h.repo.ActivateCustomerByEmail(c.Context(), email); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Customer activated successfully.",
	})
}

func (h *CustomerHandler) ResetPassword(c *fiber.Ctx) error {
	var req *domain.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	email, err := h.repo.ValidateResetToken(c.Context(), req.Token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	hashedPassword, _ := utils.HashPassword(req.Password)
	if err := h.repo.UpdatePasswordByEmail(c.Context(), email, hashedPassword); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password reset successfully.",
	})
}

func (h *CustomerHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	customer, err := h.repo.GetCustomerByEmail(c.Context(), req.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	if !customer.IsActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Account not activated",
		})
	}

	if err := utils.CheckPassword(req.Password, customer.HashedPassword); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	token, err := utils.GenerateJWT(customer.ID.String(), customer.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Set cookie in the response
	c.Cookie(&fiber.Cookie{
		Name:     "auth-token",
		Value:    token,
		Expires:  time.Now().Add(30 * time.Minute),
		HTTPOnly: false,
		Secure:   true,
		SameSite: "None",
		Path:     "/",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": fiber.Map{
			"id":    customer.ID,
			"name":  customer.Name,
			"email": customer.Email,
		},
	})
}

func (h *CustomerHandler) RefreshToken(c *fiber.Ctx) error {
	var req domain.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No token provided",
		})
	}

	token = strings.TrimPrefix(token, "Bearer ")

	_, err := utils.ValidateJWT(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	newToken, err := utils.GenerateJWT(req.ID.String(), req.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to refresh token",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "auth-token",
		Value:    newToken,
		Expires:  time.Now().Add(30 * time.Minute),
		HTTPOnly: false,
		Secure:   true,
		SameSite: "None",
		Path:     "/",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Token refreshed successfully",
	})
}

func (h *CustomerHandler) Logout(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No token provided",
		})
	}

	token = strings.TrimPrefix(token, "Bearer ")

	tokenKey := fmt.Sprintf("blacklisted-token:%s", token)
	err := h.repo.BlacklistToken(c.Context(), tokenKey, 24*time.Hour)
	if err != nil {
		log.Printf("Error blacklisting token: %v", err)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "auth-token",
		Value:    "",
		Expires:  time.Now().Add(-30 * time.Minute),
		HTTPOnly: false,
		Secure:   true,
		SameSite: "Noneict",
		Path:     "/",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Logged out successfully",
	})
}

func (h *CustomerHandler) ValidateSession(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No token provided",
		})
	}

	token = strings.TrimPrefix(token, "Bearer ")

	claims, err := utils.ValidateJWT(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"isValid": true,
		"user": fiber.Map{
			"id":    claims.UserID,
			"email": claims.Email,
		},
	})
}
