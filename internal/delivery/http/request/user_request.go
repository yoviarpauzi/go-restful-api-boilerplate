package request

type UpdateUserRequest struct {
	Name  string `json:"name" validate:"omitempty,min=3"`
	Email string `json:"email" validate:"omitempty,email"`
}
