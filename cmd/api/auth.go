package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type OIDCProvider struct {
	Provider     *oidc.Provider
	Config       oauth2.Config
	Verifier     *oidc.IDTokenVerifier
	Name         string
	ClientID     string
	ClientSecret string
	Issuer       string
}

type AuthManager struct {
	providers map[string]*OIDCProvider
	mu        sync.RWMutex
}

func NewAuthManager() *AuthManager {
	return &AuthManager{
		providers: make(map[string]*OIDCProvider),
		mu:        sync.RWMutex{},
	}
}

func (am *AuthManager) RegisterProvider(ctx context.Context, config ProviderConfig) error {

	if am == nil {
		return fmt.Errorf("AuthManager is nil")
	}
	if am.providers == nil {
		am.providers = make(map[string]*OIDCProvider)
	}
	am.mu.Lock()
	defer am.mu.Unlock()

	provider, err := oidc.NewProvider(ctx, config.ExtraConfig["issuer"])
	if err != nil {
		return fmt.Errorf("failed to initialize provider: %w", err)
	}

	oidcConfig := &oidc.Config{
		ClientID: config.ClientID,
	}

	// Configure the OAuth2 config
	oauth2Config := oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  fmt.Sprintf("%s/v1/auth/%s/callback", config.BaseCallbackURL, config.Name),
		Endpoint:     provider.Endpoint(),
		Scopes:       append([]string{oidc.ScopeOpenID, "profile", "email"}, config.Scopes...),
	}

	am.providers[config.Name] = &OIDCProvider{
		Provider:     provider,
		Config:       oauth2Config,
		Verifier:     provider.Verifier(oidcConfig),
		Name:         config.Name,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Issuer:       config.ExtraConfig["issuer"],
	}

	return nil
}

func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (a *application) HandleAuth(w http.ResponseWriter, r *http.Request) {
	provider := strings.TrimPrefix(r.URL.Path, "/v1/auth/")
	provider = strings.TrimSuffix(provider, "/")

	a.authManager.mu.RLock()
	p, exists := a.authManager.providers[provider]
	a.authManager.mu.RUnlock()

	if !exists {
		a.notFoundResponse(w, r)
		return
	}

	state, err := generateState()
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Store state in session
	session, _ := a.sessionStore.Get(r, "auth-session")
	session.Values["state"] = state
	session.Save(r, w)

	authURL := p.Config.AuthCodeURL(state)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (a *application) HandleCallback(w http.ResponseWriter, r *http.Request) {
	provider := strings.TrimPrefix(r.URL.Path, "/v1/auth/")
	provider = strings.TrimSuffix(strings.TrimSuffix(provider, "/callback"), "/")

	a.logger.Info("auth callback received",
		"provider", provider,
		"state", r.URL.Query().Get("state"),
		"code", r.URL.Query().Get("code") != "")

	a.authManager.mu.RLock()
	p, exists := a.authManager.providers[provider]
	a.authManager.mu.RUnlock()

	if !exists {
		a.logger.Error("provider not found", "provider", provider)
		a.notFoundResponse(w, r)
		return
	}

	// Verify state
	session, _ := a.sessionStore.Get(r, "auth-session")
	if r.URL.Query().Get("state") != session.Values["state"] {
		a.invalidAuthenticationTokenResponse(w, r)
		return
	}

	// Exchange code for token
	oauth2Token, err := p.Config.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		a.logger.Error("token exchange failed",
			"error", err,
			"provider", provider,
			"error_type", fmt.Sprintf("%T", err))

		var oauthError *oauth2.RetrieveError
		if errors.As(err, &oauthError) {
			a.logger.Error("oauth2 retrieve error details",
				"status_code", oauthError.Response.StatusCode,
				"body", string(oauthError.Body))
		}
		a.serverErrorResponse(w, r, err)
		return
	}

	// Extract the ID Token
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		a.serverErrorResponse(w, r, errors.New("no id_token in token response"))
		return
	}

	// Verify the ID Token
	idToken, err := p.Verifier.Verify(r.Context(), rawIDToken)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Get user info
	var claims struct {
		Email    string `json:"email"`
		Sub      string `json:"sub"`
		Name     string `json:"name"`
		Nickname string `json:"nickname"`
	}
	if err := idToken.Claims(&claims); err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Store user info in session
	session.Values["user_id"] = claims.Sub
	session.Values["email"] = claims.Email
	session.Values["name"] = claims.Name
	session.Values["provider"] = provider
	session.Values["authenticated"] = true
	session.Save(r, w)

	// Return user info
	// err = a.writeJSON(w, http.StatusOK, envelope{"user": claims}, nil)
	// if err != nil {
	// 	a.serverErrorResponse(w, r, err)
	// }

	http.Redirect(w, r, "/v1/dashboard", http.StatusSeeOther)
}

func (a *application) HandleLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := a.sessionStore.Get(r, "auth-session")
	session.Options.MaxAge = -1
	session.Save(r, w)

	err := a.writeJSON(w, http.StatusOK, envelope{"message": "Successfully logged out"}, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
