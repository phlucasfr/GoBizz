package server

import (
	"auth-service/internal/handlers"
	"auth-service/middleware"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App, customerHandler *handlers.CustomerHandler, linksHandler *handlers.LinksHandler, eventsHandler *handlers.EventsHandler, rdb *redis.Client) {
	v1 := app.Group("/v1")

	// Auth routes
	v1.Post("/auth", customerHandler.Create)
	v1.Post("/auth/login", customerHandler.Login)
	v1.Post("/auth/logout", customerHandler.Logout)
	v1.Post("/auth/recovery", customerHandler.RecoverPassword)
	v1.Post("/auth/refresh-token", customerHandler.RefreshToken)

	v1.Put("/auth/reset-password", customerHandler.ResetPassword)
	v1.Put("/auth/email-verification", customerHandler.VerifyCustomerByEmail)

	v1.Get("/auth/validate-session", customerHandler.ValidateSession)

	// Links routes - protected by auth middleware
	links := v1.Group("/links", middleware.AuthMiddleware(rdb))
	links.Post("/", linksHandler.CreateLinkHTTP)
	links.Put("/:id", linksHandler.UpdateLinkHTTP)
	links.Put("/:id/clicks", linksHandler.UpdateLinkClicksHTTP)
	links.Get("/:shortUrl", linksHandler.GetLinkHTTP)
	links.Get("/customer/:customerId", linksHandler.GetCustomerLinksHTTP)
	links.Delete("/:id", linksHandler.DeleteLinkHTTP)

	// Events routes - protected by auth middleware
	events := v1.Group("/events", middleware.AuthMiddleware(rdb))
	events.Get("/occurrences", eventsHandler.ListOccurrencesHTTP)
	events.Get("/", eventsHandler.ListEventsHTTP)
	events.Post("/", eventsHandler.CreateEventHTTP)
	events.Put("/:id", eventsHandler.UpdateEventHTTP)
	events.Get("/:id", eventsHandler.GetEventHTTP)
	events.Delete("/:id", eventsHandler.DeleteOrCutEventHTTP)
}
