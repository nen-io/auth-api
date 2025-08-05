package services

import (
	"context"
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

type GoogleClaims struct {
	Email    string `json:"email"`
	Verified bool   `json:"email_verified"`
}

func InitGoogleOAuth() {

	oauthCfg = &oauth2.Config{
		ClientID:     os.Getenv("G_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("G_OAUTH_SECRET_KEY"),
		RedirectURL:  os.Getenv("G_OAUTH_CALLBACK"),
		Endpoint:     google.Endpoint, // helper pkg  [oai_citation:10‡pkg.go.dev](https://pkg.go.dev/golang.org/x/oauth2/google?utm_source=chatgpt.com)
		Scopes:       []string{"openid", "email", "profile"},
	}

	var err error
	oidcProv, err = oidc.NewProvider(context.TODO(), "https://accounts.google.com")
	if err != nil {
		log.Fatalf("oidc provider: %v", err)
	}
}

func HandleGoogleAuth(c fiber.Ctx) error {

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

	return c.Redirect().To(url)
}

func HandleGoogleCallback(c fiber.Ctx) (GoogleClaims, error) {
	if c.Params("provider") != "google" {
		return GoogleClaims{}, c.SendStatus(fiber.StatusBadRequest)
	}
	sess := session.FromContext(c)
	if c.Query("state") != sess.Get("state") {
		return GoogleClaims{}, c.SendStatus(fiber.StatusBadRequest)
	}
	verifier := sess.Get("verifier").(string)

	token, err := oauthCfg.Exchange(c.Context(), c.Query("code"),
		oauth2.SetAuthURLParam("code_verifier", verifier))
	if err != nil {
		return GoogleClaims{}, err
	}

	// Set Cookies for access and refresh tokens
	accessTokenCookie := &fiber.Cookie{
		Name:     "accessToken",
		Value:    token.AccessToken,
		Expires:  time.Now().Add(time.Minute * 50),
		HTTPOnly: true,
	}

	c.Cookie(accessTokenCookie)

	// ----- verify ID-token -----
	rawID, ok := token.Extra("id_token").(string)
	if !ok {
		return GoogleClaims{}, fiber.NewError(fiber.StatusInternalServerError, "missing id_token")
	}

	slog.Info("Verifying Google ID token", "token", token)

	verifierOIDC := oidcProv.Verifier(&oidc.Config{ClientID: oauthCfg.ClientID})
	idTok, err := verifierOIDC.Verify(c.Context(), rawID)
	if err != nil {
		return GoogleClaims{}, err
	}
	var claims GoogleClaims
	if err := idTok.Claims(&claims); err != nil {
		return GoogleClaims{}, err
	}

	// save session & redirect

	j, err := json.Marshal(&token)
	if err != nil {
		slog.Error("Failed to marshal token", "error", err)
		return GoogleClaims{}, fiber.NewError(fiber.StatusInternalServerError, "failed to marshal token")
	}
	sess.Set("provider", "google")
	sess.Set("token", string(j))
	sess.Set("id_token", rawID)

	c.Locals("email", claims.Email)

	return claims, nil

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
		slog.Error("No valid Google token found in session")
		return nil, fiber.NewError(fiber.StatusUnauthorized, "no valid Google token found in session")
	}

	accessToken := c.Request().Header.Cookie("accessToken")

	if string(accessToken) == "" {
		slog.Error("Missing access or refresh token in request headers")
		return nil, fiber.NewError(fiber.StatusUnauthorized, "missing access token")
	}

	if !ok || tokenStr == "" {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "no valid Google token found in session")
	}

	var token oauth2.Token
	if err := json.Unmarshal([]byte(tokenStr), &token); err != nil {
		slog.Error("Failed to unmarshal token", "error", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to unmarshal token")
	}

	if token.AccessToken != string(accessToken) {
		slog.Error("access token and refresh token do not match session token")
		return nil, fiber.NewError(fiber.StatusUnauthorized, "access or refresh token does not match the expected value")
	}

	clientID := os.Getenv("G_OAUTH_CLIENT_ID")

	rawID := sess.Get("id_token").(string)
	if rawID == "" {
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
		return nil, fiber.NewError(fiber.StatusUnauthorized, "no valid Google token found in session")
	}

	c.Locals("email", claims.Email)
	c.Locals("username", claims.Email)

	return &TokenInfo{
		Email:          claims.Email,
		IDTokenTTL:     time.Until(idTok.Expiry),
		AccessTokenTTL: time.Until(token.Expiry),
	}, nil
}
