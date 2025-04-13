package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"auth-service/internal/domain"
	"auth-service/internal/handlers"
	"auth-service/internal/test"
	"auth-service/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/stretchr/testify/require"
)

func TestServerRoutes(t *testing.T) {
	repositories, err := test.SetupTestContainers()
	require.NoError(t, err, "Test containers setup should not produce an error")

	customerHandler := handlers.NewCustomerHandler(repositories.CustomerRepository)

	app := fiber.New(fiber.Config{
		AppName: "auth-service test API",
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
	v1.Post("/customers", customerHandler.Create)

	t.Run("POST /v1/customers", func(t *testing.T) {
		params := domain.CreateCustomerRequest{
			Name:     utils.RandomString(16),
			Email:    "test@gmail.com",
			Phone:    "47999999999",
			CPFCNPJ:  "99999999999",
			Password: utils.RandomString(8),
		}
		jsonData, _ := json.Marshal(params)

		req := httptest.NewRequest("POST", "/v1/customers", bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err, "Request to /v1/customers should not produce an error")

		var createResponse domain.CreateCustomerResponse
		bodyBytes, err := io.ReadAll(resp.Body)
		require.NoError(t, err, "Reading response body should not produce an error")

		err = json.Unmarshal(bodyBytes, &createResponse)
		require.NoError(t, err, "Unmarshaling response body should not produce an error")

		require.Equal(t, http.StatusCreated, resp.StatusCode, "Status code should be 201 Created")
	})
}
