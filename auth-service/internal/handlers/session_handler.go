package handlers

import (
	"auth-service/internal/domain"
	"auth-service/internal/infra/repository"
	"auth-service/pkg/util"
	"log"
	"os"

	"time"

	"github.com/gofiber/fiber/v2"
)

type SessionHandler struct {
	repo *repository.SessionRepository
}

func NewSessionHandler(repo *repository.SessionRepository) *SessionHandler {
	return &SessionHandler{repo: repo}
}

func (h *SessionHandler) CreateSession(c *fiber.Ctx) error {
	var req domain.CreateSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	token, err := util.GenerateJWT(req.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	session := &domain.Session{
		ID:        token,
		UserID:    req.UserID,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}

	err = h.repo.Create(c.Context(), session)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to store session",
		})
	}

	secureMode := os.Getenv("RAILWAY_ENVIRONMENT_NAME") == "production"

	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    token,
		Secure:   secureMode,
		Expires:  session.ExpiresAt,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteNoneMode,
	})

	return c.JSON(fiber.Map{
		"token":   token,
		"message": "Session created successfully",
	})
}

func (h *SessionHandler) ValidateSession(c *fiber.Ctx) error {
	tokenString := c.Cookies("session_id")
	log.Println("validatesession>tokenstring:", tokenString)
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Session ID is missing",
		})
	}

	err := util.ValidateJWT(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	session, err := h.repo.Validate(c.Context(), tokenString)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to validate session",
		})
	}

	if session == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired session",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Session is valid",
		"session": session,
	})
}

func (h *SessionHandler) DeleteSession(c *fiber.Ctx) error {
	sessionID := c.Cookies("session_id")
	if sessionID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Session ID is missing",
		})
	}

	err := h.repo.Delete(c.Context(), sessionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete session",
		})
	}

	c.ClearCookie("session_id")

	return c.JSON(fiber.Map{
		"message": "Session deleted successfully",
	})
}
