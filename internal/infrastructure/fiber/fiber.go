package fiber

import (
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewFiber(config *viper.Viper, log *zap.Logger) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: config.GetString("APP_NAME"),
	})

	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: log,
	}))

	return app
}
