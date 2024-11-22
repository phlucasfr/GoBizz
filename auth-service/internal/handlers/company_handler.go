package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"auth-service/internal/domain"
	"auth-service/internal/infra/repository"
)

type CompanyHandler struct {
	repo *repository.CompanyRepository
}

func NewCompanyHandler(repo *repository.CompanyRepository) *CompanyHandler {
	return &CompanyHandler{repo: repo}
}

func (h *CompanyHandler) Create(c *fiber.Ctx) error {
	var req domain.CreateCompanyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	company, err := h.repo.Create(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(company)
}

func (h *CompanyHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid company ID",
		})
	}

	company, err := h.repo.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Company not found",
		})
	}

	return c.JSON(company)
}

func (h *CompanyHandler) VerifyCompanyBySms(c *fiber.Ctx) error {
	var req domain.VerifyCompanyBySmsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	company, err := h.repo.VerifyCompanyBySms(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(company)
}
