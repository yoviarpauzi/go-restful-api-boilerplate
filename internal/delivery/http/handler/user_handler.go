package handler

import (
	"go-restful-api/internal/delivery/http/request"
	"go-restful-api/internal/delivery/http/response"
	"go-restful-api/internal/domain/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type UserHandler struct {
	UserUseCase usecase.UserUseCase
	Log         *zap.Logger
	Validate    *validator.Validate
}

func NewUserHandler(userUseCase usecase.UserUseCase, log *zap.Logger, validate *validator.Validate) *UserHandler {
	return &UserHandler{
		UserUseCase: userUseCase,
		Log:         log,
		Validate:    validate,
	}
}

// GetByID godoc
// @Summary Get user profile by ID
// @Description Fetch user details by their unique ID
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Success 200 {object} response.SuccessResponse "User details"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "User not found"
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	userID := c.Params("id")
	user, err := h.UserUseCase.GetProfile(c.UserContext(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse{
			Success: false,
			Error: response.ErrorDetail{
				Code:    "NOT_FOUND",
				Message: "user not found",
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Success: true,
		Data: response.UserResponse{
			ID:        user.ID.String(),
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
		Message: "fetch profile successfully",
	})
}

// UpdateByID godoc
// @Summary Update user profile by ID
// @Description Update name or email for a user by their ID
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Param request body request.UpdateUserRequest true "Update User Profile"
// @Success 200 {object} response.SuccessResponse "User updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 404 {object} response.ErrorResponse "User not found"
// @Router /users/{id} [put]
func (h *UserHandler) UpdateByID(c *fiber.Ctx) error {
	req := new(request.UpdateUserRequest)
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

	userID := c.Params("id")
	user, err := h.UserUseCase.GetProfile(c.UserContext(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse{
			Success: false,
			Error: response.ErrorDetail{
				Code:    "NOT_FOUND",
				Message: "user not found",
			},
		})
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := h.UserUseCase.UpdateUser(c.UserContext(), user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Success: false,
			Error: response.ErrorDetail{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "failed to update user",
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Success: true,
		Message: "user profile updated successfully",
	})
}

// GetAllUsers godoc
// @Summary List all users
// @Description Admin or internal use case to list all users with pagination
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Page number (default 1)"
// @Param size query int false "Page size (default 10)"
// @Success 200 {object} response.PaginatedResponse "List of users with metadata"
// @Router /users [get]
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	size := c.QueryInt("size", 10)

	users, total, err := h.UserUseCase.GetAllUsers(c.UserContext(), page, size)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch users"})
	}

	var userList []response.UserResponse
	for _, u := range users {
		userList = append(userList, response.UserResponse{
			ID:        u.ID.String(),
			Name:      u.Name,
			Email:     u.Email,
			CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: u.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	totalPages := int((total + int64(size) - 1) / int64(size))

	return c.Status(fiber.StatusOK).JSON(response.PaginatedResponse{
		Success: true,
		Data:    userList,
		Meta: response.PaginationMeta{
			Total:      total,
			Page:       page,
			PerPage:    size,
			TotalPages: totalPages,
		},
	})
}

// DeleteByID godoc
// @Summary Delete user by ID
// @Description Permanently delete a user by their unique ID
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Success 200 {object} response.SuccessResponse "User deleted successfully"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "User not found"
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteByID(c *fiber.Ctx) error {
	userID := c.Params("id")
	if err := h.UserUseCase.DeleteUser(c.UserContext(), userID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse{
			Success: false,
			Error: response.ErrorDetail{
				Code:    "NOT_FOUND",
				Message: "user not found or delete failed",
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Success: true,
		Message: "user deleted successfully",
	})
}
