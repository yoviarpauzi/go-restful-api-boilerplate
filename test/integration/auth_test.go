package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-restful-api/internal/config"
	"go-restful-api/internal/delivery/http/request"
	"go-restful-api/internal/delivery/http/response"
	"go-restful-api/internal/domain/entity"
	"go-restful-api/internal/infrastructure/database"
	"go-restful-api/internal/infrastructure/fiber"
	"go-restful-api/internal/infrastructure/logger"
	"go-restful-api/internal/infrastructure/validation"

	"github.com/joho/godotenv"
	infraConfig "go-restful-api/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// AuthSuite defines the suite for auth integration tests
type AuthSuite struct {
	suite.Suite
	gormDB *gorm.DB
	app    *config.BootstrapConfig
}

func (s *AuthSuite) SetupSuite() {
	// Load .env
	_ = godotenv.Load("../../.env")

	// Initialize Viper Config
	v := infraConfig.NewViper()
	
	// Force some test-specific config if needed, or rely on .env
	v.Set("PASETO_ACCESS_TOKEN_SECRET", "12345678901234567890123456789012")
	v.Set("PASETO_REFRESH_TOKEN_SECRET", "12345678901234567890123456789012")

	log := logger.NewLogger(v)
	validate := validation.NewValidator()

	// Initialize Real Database
	s.gormDB = database.NewDatabase(v, log)

	// Migrate User entity
	_ = s.gormDB.Migrator().DropTable(&entity.User{})
	_ = s.gormDB.AutoMigrate(&entity.User{})

	// Initialize Fiber App
	fiberApp := fiber.NewFiber(v, log)

	s.app = &config.BootstrapConfig{
		DB:       s.gormDB,
		App:      fiberApp,
		Log:      log,
		Validate: validate,
		Config:   v,
	}

	config.Bootstrap(s.app)
	s.T().Log("SetupSuite (Real DB) completed")
}

func (s *AuthSuite) SetupTest() {
	// Clean up database before each test (Hard Delete)
	s.gormDB.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&entity.User{})
}

func (s *AuthSuite) TestRegister_Success() {
	reqBody := request.RegisterRequest{
		Name:     "Integration User",
		Email:    "int@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := s.app.App.Test(req)

	assert.Equal(s.T(), http.StatusCreated, resp.StatusCode)

	var result response.SuccessResponse
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(s.T(), result.Success)

	// Verify database state
	var user entity.User
	err := s.gormDB.Where("email = ?", "int@example.com").First(&user).Error
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "Integration User", user.Name)
}

func (s *AuthSuite) TestLogin_Success() {
	// First register a user
	regReq := request.RegisterRequest{
		Name:     "Login User",
		Email:    "login@example.com",
		Password: "password123",
	}
	regBody, _ := json.Marshal(regReq)
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(regBody))
	req1.Header.Set("Content-Type", "application/json")
	s.app.App.Test(req1)

	// Now try to Login
	loginReq := request.LoginRequest{
		Email:    "login@example.com",
		Password: "password123",
	}
	loginBody, _ := json.Marshal(loginReq)
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	req2.Header.Set("Content-Type", "application/json")
	resp, _ := s.app.App.Test(req2)

	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

	var result response.SuccessResponse
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(s.T(), result.Success)

	// Verify cookie (refresh token)
	cookies := resp.Cookies()
	var refreshTokenFound bool
	for _, cookie := range cookies {
		if cookie.Name == "refresh_token" {
			refreshTokenFound = true
			assert.NotEmpty(s.T(), cookie.Value)
			assert.True(s.T(), cookie.HttpOnly)
		}
	}
	assert.True(s.T(), refreshTokenFound)

	// Verify response data (access token)
	data := result.Data.(map[string]interface{})
	assert.NotEmpty(s.T(), data["access_token"])
}

func TestAuthIntegration(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}
