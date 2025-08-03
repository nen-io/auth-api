package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"example.com/src/utils"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	oauthCfg *oauth2.Config
	oidcProv *oidc.Provider
)

func InitGoogleOAuth() {

	oauthCfg = &oauth2.Config{
		ClientID:     os.Getenv("G_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("G_OAUTH_SECRET_KEY"),
		RedirectURL:  os.Getenv("G_OAUTH_CALLBACK"),
		Endpoint:     google.Endpoint, // helper pkg  [oai_citation:10‡pkg.go.dev](https://pkg.go.dev/golang.org/x/oauth2/google?utm_source=chatgpt.com)
		Scopes:       []string{"openid", "email", "profile"},
	}

	slog.Info("Initializing Google OAuth....", "clientId", oauthCfg.ClientID, "clientSecret", oauthCfg.ClientSecret)

	var err error
	oidcProv, err = oidc.NewProvider(context.TODO(), "https://accounts.google.com")
	if err != nil {
		log.Fatalf("oidc provider: %v", err)
	}
}

func HandleGoogleAuth(c fiber.Ctx) error {

	slog.Info("google env details", "clientId", os.Getenv("G_OAUTH_CLIENT_ID"), "clientSecret", os.Getenv("G_OAUTH_SECRET_KEY"))

	if c.Params("provider") != "google" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	state := utils.RandString(24)

	verifier, challenge, err := utils.PCKE_Pair()

	if err != nil {
		slog.Error("Failed to generate PCKE pair", "error", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	sess, err := GetSession(c)
	if err != nil {
		slog.Error("Failed to get session", "error", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	sess.Set("state", state)
	sess.Set("verifier", verifier)

	url := oauthCfg.AuthCodeURL(state,
		oauth2.SetAuthURLParam("code_challenge", challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)

	slog.Info("Redirecting to Google OAuth", "url", url)

	return c.Redirect().To(url)
}

func HandleGoogleCallback(c fiber.Ctx) error {
	if c.Params("provider") != "google" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	sess := session.FromContext(c)
	if c.Query("state") != sess.Get("state") {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	verifier := sess.Get("verifier").(string)

	token, err := oauthCfg.Exchange(c.Context(), c.Query("code"),
		oauth2.SetAuthURLParam("code_verifier", verifier))
	if err != nil {
		return err
	}

	// Set Cookies for access and refresh tokens
	accessTokenCookie := &fiber.Cookie{
		Name:     "accessToken",
		Value:    token.AccessToken,
		Expires:  time.Now().Add(time.Minute * 15),
		HTTPOnly: true,
	}
	refreshTokenCookie := &fiber.Cookie{
		Name:     "refreshToken",
		Value:    token.RefreshToken,
		Expires:  time.Now().Add(time.Hour * 24 * 15),
		HTTPOnly: true,
	}

	c.Cookie(accessTokenCookie)
	c.Cookie(refreshTokenCookie)

	// ----- verify ID-token -----
	rawID, ok := token.Extra("id_token").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "missing id_token")
	}
	verifierOIDC := oidcProv.Verifier(&oidc.Config{ClientID: oauthCfg.ClientID})
	idTok, err := verifierOIDC.Verify(c.Context(), rawID)
	if err != nil {
		return err
	}
	var claims struct {
		Email    string `json:"email"`
		Verified bool   `json:"email_verified"`
	}
	if err := idTok.Claims(&claims); err != nil {
		return err
	}

	// save session & redirect

	j, err := json.Marshal(&token)
	if err != nil {
		slog.Error("Failed to marshal token", "error", err)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to marshal token")
	}
	sess.Set("provider", "google")
	sess.Set("token", string(j))

	c.Locals("email", claims.Email)

	return c.Redirect().To(os.Getenv("UI_DOMAIN") + "/home")
}

type TokenInfo struct {
	Email          string
	IDTokenTTL     time.Duration
	AccessTokenTTL time.Duration
}

func ValidateGoogleTokens(c fiber.Ctx) (*TokenInfo, error) {
	// 1) Verify ID token
	ctx := c.Context()
	sess := session.FromContext(c)
	tokenStr, ok := sess.Get("token").(string)
	if !ok || tokenStr == "" {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "no valid Google token found in session")
	}

	var token oauth2.Token
	if err := json.Unmarshal([]byte(tokenStr), &token); err != nil {
		slog.Error("Failed to unmarshal token", "error", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to unmarshal token")
	}

	clientID := os.Getenv("G_OAUTH_CLIENT_ID")

	rawID, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "missing id_token")
	}
	ver := oidcProv.Verifier(&oidc.Config{ClientID: clientID})
	idTok, err := ver.Verify(ctx, rawID)
	if err != nil {
		return nil, fmt.Errorf("id token verify failed: %w", err)
	}
	var claims struct {
		Email string `json:"email"`
	}
	if err := idTok.Claims(&claims); err != nil {
		return nil, err
	}

	// 2) Check access token lifetime
	if !token.Valid() {
		return nil, errors.New("access token already expired/invalid")
	}

	slog.Info("Google token validated", "email", claims.Email, "expires_in", token.Expiry)
	c.Locals("email", claims.Email)

	return &TokenInfo{
		Email:          claims.Email,
		IDTokenTTL:     time.Until(idTok.Expiry),
		AccessTokenTTL: time.Until(token.Expiry),
	}, nil
}
