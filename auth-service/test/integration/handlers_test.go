package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"auth-service/internal/domain"
	"auth-service/internal/handlers"
	"auth-service/pkg/util"
	"auth-service/test"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestServerRoutes(t *testing.T) {
	repositories, err := test.SetupTestContainers(t)
	assert.NoError(t, err)

	companyHandler := handlers.NewCompanyHandler(repositories.CompanyRepository, repositories.SessionRepository)

	app := fiber.New(fiber.Config{
		AppName: "Your Company Management API",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
	}))

	v1 := app.Group("/v1")
	v1.Post("/companies", companyHandler.Create)
	v1.Get("/companies/:id", companyHandler.GetByID)

	var companyID uuid.UUID
	t.Run("POST /v1/companies", func(t *testing.T) {
		params := domain.CreateCompanyRequest{
			Name:     util.RandomString(16),
			Email:    "test@gmail.com",
			Phone:    "47999999999",
			CPFCNPJ:  "99999999999",
			Password: util.RandomString(8),
		}
		jsonData, _ := json.Marshal(params)
		req := httptest.NewRequest("POST", "/v1/companies", bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)

		var createResponse domain.CreateCompanyResponse
		bodyBytes, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		err = json.Unmarshal(bodyBytes, &createResponse)
		assert.NoError(t, err)

		companyID = createResponse.ID
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	//TODO: Fix test
	t.Run("GET /v1/companies/:id", func(t *testing.T) {
		req := httptest.NewRequest("GET", fmt.Sprintf("/v1/companies/%s", companyID.String()), nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET /v1/companies/:id", func(t *testing.T) {
		companyID = uuid.New()

		req := httptest.NewRequest("GET", fmt.Sprintf("/v1/companies/%s", companyID.String()), nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("GET /v1/companies/:id", func(t *testing.T) {

		req := httptest.NewRequest("GET", fmt.Sprintf("/v1/companies/%s", "1"), nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
