package models

type RequestEmail struct {
	Email string `json:"email"`
}

type ChangePassword struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
	Id          string `json:"id"`
}
