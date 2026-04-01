package main

import (
	"context"
	_ "go-restful-api/docs"                                     // Side-effect import for Swagger documentation
	"go-restful-api/internal/config"                            // Import for Bootstrap and BootstrapConfig
	infraConfig "go-restful-api/internal/infrastructure/config" // Import for NewViper
	"go-restful-api/internal/domain/entity"
	"go-restful-api/internal/infrastructure/database"
	"go-restful-api/internal/infrastructure/fiber"
	"go-restful-api/internal/infrastructure/logger"
	"go-restful-api/internal/infrastructure/validation"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// @title Go Restful API
// @version 1.0
// @description This is a sample server for Go Restful API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load environment variables from .env file if it exists
	_ = godotenv.Load()

	viperConfig := infraConfig.NewViper()
	zapLog := logger.NewLogger(viperConfig)
	gormDB := database.NewDatabase(viperConfig, zapLog)
	validator := validation.NewValidator()
	fiberApp := fiber.NewFiber(viperConfig, zapLog)
	
	// GORM AutoMigrate
	if err := gormDB.AutoMigrate(&entity.User{}); err != nil {
		zapLog.Fatal("failed to migrate database", zap.Error(err))
	}

	bootstrapConfig := &config.BootstrapConfig{
		DB:       gormDB,
		App:      fiberApp,
		Log:      zapLog,
		Validate: validator,
		Config:   viperConfig,
	}

	config.Bootstrap(bootstrapConfig)

	errCh := make(chan error, 1)
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		appPort := viperConfig.GetString("APP_PORT")

		zapLog.Info("Server is starting", zap.String("port", appPort), zap.String("app_name", viperConfig.GetString("APP_NAME")))
		if err := fiberApp.Listen(":" + appPort); err != nil {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		zapLog.Fatal("Server start failed", zap.Error(err))
	case sig := <-quitCh:
		zapLog.Info("Shutdown signal received", zap.String("signal", sig.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := fiberApp.ShutdownWithContext(ctx); err != nil {
			zapLog.Fatal("Server shutdown failed", zap.Error(err))
		}

		zapLog.Info("Server stopped successfully")
	}
}
