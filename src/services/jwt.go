package services

import (
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(username string, id string, tokenType string, duration time.Duration) (string, error) {

	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		return "", errors.New("JWT_SECRET not set")
	}

	slog.Info(id)

	claims := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"email": username,
		"id":    id,
		"sub":   tokenType,
		"exp":   time.Now().Add(time.Minute * 15).Unix(), // Expiration time
		"iat":   time.Now().Unix(),                       // Issued at
	})

	token, err := claims.SignedString([]byte(jwtSecret))

	if err != nil {
		return "", err
	}

	return token, nil
}
