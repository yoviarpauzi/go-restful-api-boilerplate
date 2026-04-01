package config

import (
	"go-restful-api/internal/delivery/http/handler"
	"go-restful-api/internal/delivery/http/middleware"
	"go-restful-api/internal/delivery/http/route"
	"go-restful-api/internal/infrastructure/database"
	"go-restful-api/internal/infrastructure/token"
	"go-restful-api/internal/repository"
	"go-restful-api/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *zap.Logger
	Validate *validator.Validate
	Config   *viper.Viper
	Token    *token.TokenGenerator
}

func Bootstrap(config *BootstrapConfig) {
	// 1. Setup Infrastructure components
	tokenGen, err := token.NewTokenGenerator(config.Config)
	if err != nil {
		config.Log.Fatal("failed to initialize token generator", zap.Error(err))
	}

	// 2. Setup Repositories
	userRepository := repository.NewUserRepository(config.DB, config.Log)
	txManager := database.NewTransactionManager(config.DB)

	// 3. Setup Use Cases
	userUseCase := usecase.NewUserUseCase(userRepository, config.Log)
	authUseCase := usecase.NewAuthUseCase(userRepository, tokenGen, config.Log, txManager)

	// 4. Setup Handlers
	userHandler := handler.NewUserHandler(userUseCase, config.Log, config.Validate)
	authHandler := handler.NewAuthHandler(authUseCase, config.Log, config.Validate)

	// 5. Setup Middleware
	authMiddleware := middleware.NewAuth(tokenGen)

	// 6. Setup Routes
	routeConfig := route.RouteConfig{
		App:            config.App,
		UserHandler:    userHandler,
		AuthHandler:    authHandler,
		AuthMiddleware: authMiddleware,
	}

	routeConfig.Setup()
}
