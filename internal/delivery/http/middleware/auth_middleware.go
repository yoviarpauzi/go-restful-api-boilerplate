package middleware

import (
	"strings"

	"go-restful-api/internal/infrastructure/token"

	"github.com/gofiber/fiber/v2"
)

func NewAuth(tokenGenerator *token.TokenGenerator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "UNAUTHORIZED",
					"message": "unauthorized",
				},
			})
		}

		fields := strings.Fields(authHeader)
		if len(fields) != 2 || strings.ToLower(fields[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "UNAUTHORIZED",
					"message": "invalid authorization header format",
				},
			})
		}

		accessToken := fields[1]
		userID, err := tokenGenerator.ValidateAccessToken(accessToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "UNAUTHORIZED",
					"message": "invalid token",
				},
			})
		}

		c.Locals("currentUser", userID)
		return c.Next()
	}
}
