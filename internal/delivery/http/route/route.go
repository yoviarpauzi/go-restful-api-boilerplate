package route

import (
	"go-restful-api/internal/delivery/http/handler"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/huxulm/fiber-swagger-v2"
)

type RouteConfig struct {
	App            *fiber.App
	AuthHandler    *handler.AuthHandler
	UserHandler    *handler.UserHandler
	AuthMiddleware fiber.Handler
}

func (c *RouteConfig) Setup() {
	authConfig := AuthRouteConfig{
		App:            c.App,
		AuthHandler:    c.AuthHandler,
		AuthMiddleware: c.AuthMiddleware,
	}

	authConfig.Setup()

	userConfig := UserRouteConfig{
		App:            c.App,
		UserHandler:    c.UserHandler,
		AuthMiddleware: c.AuthMiddleware,
	}
	userConfig.Setup()

	api := c.App.Group("/api/v1")
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "healthy",
		})
	})

	c.App.Get("/swagger/*", fiberSwagger.WrapHandler)
}
