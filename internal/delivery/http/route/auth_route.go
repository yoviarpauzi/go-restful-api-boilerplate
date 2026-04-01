package route

import (
	"go-restful-api/internal/delivery/http/handler"

	"github.com/gofiber/fiber/v2"
)

type AuthRouteConfig struct {
	App            *fiber.App
	AuthHandler    *handler.AuthHandler
	AuthMiddleware fiber.Handler
}

func (c *AuthRouteConfig) Setup() {
	auth := c.App.Group("/api/v1/auth")

	auth.Post("/register", c.AuthHandler.Register)
	auth.Post("/login", c.AuthHandler.Login)
	auth.Post("/reset-password", c.AuthHandler.ResetPassword)

	auth.Post("/change-password", c.AuthMiddleware, c.AuthHandler.ChangePassword)
}
