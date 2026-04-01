package route

import (
	"go-restful-api/internal/delivery/http/handler"

	"github.com/gofiber/fiber/v2"
)

type UserRouteConfig struct {
	App            *fiber.App
	UserHandler    *handler.UserHandler
	AuthMiddleware fiber.Handler
}

func (c *UserRouteConfig) Setup() {
	users := c.App.Group("/api/v1/users", c.AuthMiddleware)

	users.Get("/:id", c.UserHandler.GetByID)
	users.Put("/:id", c.UserHandler.UpdateByID)
	users.Delete("/:id", c.UserHandler.DeleteByID)
	users.Get("/", c.UserHandler.GetAllUsers)
}
