package controllers

import (
	"errors"
	"log/slog"
	"time"

	"example.com/src/models"
	"example.com/src/services"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"gorm.io/gorm"
)

func Login(c fiber.Ctx) error {

	body := c.Request().Body()

	loginRequest := new(models.LoginRequest)

	if err := c.Bind().JSON(loginRequest); err != nil {
		slog.Error("login decode body", "error", err, "body", body)
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("Invalid login request body"))
	}

	// only fetch what we need for login
	user := models.User{
		Email: loginRequest.Email,
	}
	userLoginFields := []string{"password", "verified", "id", "user_name", "email"}

	if err := user.ValidateEmail(loginRequest.Email); err != nil {
		user.Email = ""
		user.UserName = loginRequest.Email
	}

	if err := user.GetFields(userLoginFields); errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Error("Failed to find user", "User:", user)
		return c.Status(fiber.StatusForbidden).JSON(models.MakeError("Failed to login"))
	}

	// check password
	if err := user.ComparePassword(loginRequest.Password); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(models.MakeError("Failed to login"))
	}

	if user.Verified == false {
		resp := map[string]string{
			"status":  "VERIFY EMAIL",
			"message": "User not verified",
			"email":   user.Email,
		}

		return c.JSON(resp)
	}

	// generate jwt tokens
	accessToken, err := services.CreateJWT(user.UserName, user.Email, user.ID, "accessToken", time.Minute*15)
	if err != nil {
		slog.Error("Failed to Create AccessToken", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.MakeError("Failed to login"))
	}

	refreshToken, err := services.CreateJWT(user.UserName, user.Email, user.ID, "refreshToken", time.Hour*24*15)
	if err != nil {
		slog.Error("Failed to Create RefreshToken", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.MakeError("Failed to login"))

	}

	accessTokenCookie := &fiber.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		Expires:  time.Now().Add(time.Minute * 15),
		HTTPOnly: true,
	}

	refreshTokenCookie := &fiber.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Hour * 24 * 15),
		HTTPOnly: true,
	}

	sess := session.FromContext(c)

	sess.Set("user", user.Email)
	sess.Set("accessToken", accessToken)
	sess.Set("provider", "google")

	c.Cookie(refreshTokenCookie)
	c.Cookie(accessTokenCookie)

	// create a session
	services.SessionManager.State[user.ID] = true

	loginResp := map[string]any{
		"success":  true,
		"username": user.UserName,
		"email":    user.Email,
		"message":  "User logged in",
	}

	return c.JSON(loginResp)

}
