package models

type RequestEmail struct {
	Email string `json:"email" validate:"required,email"`
}

type ChangePassword struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
	Id          string `json:"id" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ForgotPassword struct {
	Email string `json:"email" validate:"required,email"`
}

type ForgotPasswordChange struct {
	Email       string `json:"email" validate:"required,email"`
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
}
