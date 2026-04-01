package handler

import (
	"go-restful-api/internal/delivery/http/request"
	"go-restful-api/internal/delivery/http/response"
	"go-restful-api/internal/domain/entity"
	"go-restful-api/internal/domain/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AuthHandler struct {
	AuthUseCase usecase.AuthUseCase
	Log         *zap.Logger
	Validate    *validator.Validate
}

func NewAuthHandler(authUseCase usecase.AuthUseCase, log *zap.Logger, validate *validator.Validate) *AuthHandler {
	return &AuthHandler{
		AuthUseCase: authUseCase,
		Log:         log,
		Validate:    validate,
	}
}

// Register godoc
// @Summary Register new user
// @Description Create a new account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body request.RegisterRequest true "Register User"
// @Success 201 {object} response.SuccessResponse "User created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	req := new(request.RegisterRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Success: false,
			Error: response.ErrorDetail{
				Code:    "BAD_REQUEST",
				Message: "cannot parse JSON",
			},
		})
	}

	if err := h.Validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Success: false,
			Error: response.ErrorDetail{
				Code:    "VALIDATION_ERROR",
				Message: "validation failed",
			},
		})
	}

	user := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.AuthUseCase.Register(c.UserContext(), user); err != nil {
		h.Log.Error("register failed", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Success: false,
			Error: response.ErrorDetail{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "failed to register user",
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse{
		Success: true,
		Data: response.UserResponse{
			ID:        user.ID.String(),
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
		Message: "register successfully",
	})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and get tokens plus user ID
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "Login Credentials"
// @Success 200 {object} response.SuccessResponse{data=response.TokenResponse} "Login successful"
// @Failure 401 {object} response.ErrorResponse "Invalid credentials"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	req := new(request.LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Success: false,
			Error: response.ErrorDetail{
				Code:    "BAD_REQUEST",
				Message: "cannot parse JSON",
			},
		})
	}

	accessToken, refreshToken, userID, err := h.AuthUseCase.Login(c.UserContext(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
			Success: false,
			Error: response.ErrorDetail{
				Code:    "UNAUTHORIZED",
				Message: "invalid credentials",
			},
		})
	}

	// Set Refresh Token in Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HTTPOnly: true,
		Secure:   false, // Set to true in production
		SameSite: "Lax",
	})

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Success: true,
		Data: response.TokenResponse{
			AccessToken: accessToken,
			UserID:      userID,
		},
		Message: "login successfully",
	})
}

// ChangePassword godoc
// @Summary Change user password
// @Description Update password for the logged-in user
// @Tags Auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body request.ChangePasswordRequest true "Change Password Request"
// @Success 200 {object} response.SuccessResponse "Password changed successfully"
// @Failure 401 {object} response.ErrorResponse "Old password does not match"
// @Router /auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	req := new(request.ChangePasswordRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Success: false,
			Error: response.ErrorDetail{
				Code:    "BAD_REQUEST",
				Message: "cannot parse JSON",
			},
		})
	}

	if err := h.Validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Success: false,
			Error: response.ErrorDetail{
				Code:    "VALIDATION_ERROR",
				Message: "validation failed",
			},
		})
	}

	userID := c.Locals("currentUser").(string)
	if err := h.AuthUseCase.ChangePassword(c.UserContext(), userID, req.OldPassword, req.NewPassword); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
			Success: false,
			Error: response.ErrorDetail{
				Code:    "UNAUTHORIZED",
				Message: "old password does not match",
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Success: true,
		Message: "password changed successfully",
	})
}

// ResetPassword godoc
// @Summary Reset user password (Mock)
// @Description Reset password without email verification (only for dev)
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body request.ResetPasswordRequest true "Reset Password Request"
// @Success 200 {object} response.SuccessResponse "Password reset successfully"
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	req := new(request.ResetPasswordRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Success: false,
			Error: response.ErrorDetail{
				Code:    "BAD_REQUEST",
				Message: "cannot parse JSON",
			},
		})
	}

	if err := h.Validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Success: false,
			Error: response.ErrorDetail{
				Code:    "VALIDATION_ERROR",
				Message: "validation failed",
			},
		})
	}

	if err := h.AuthUseCase.ResetPassword(c.UserContext(), req.Email, req.NewPassword); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Success: false,
			Error: response.ErrorDetail{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "failed to reset password",
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Success: true,
		Message: "password reset successfully",
	})
}
