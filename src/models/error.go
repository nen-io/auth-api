package models

type ApiError struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func MakeError(message string) ApiError {
	return ApiError{
		Error:   true,
		Message: message,
	}
}
