package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/auth0"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

func (app *application) authCallbackHandler(w http.ResponseWriter, r *http.Request) {
	provider, err := app.readProviderParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	app.logger.Info("processing auth callback", "provider", provider)

	// Set the provider in the context
	r = app.setProviderContext(r, provider)

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		app.logger.Error("auth completion failed", "error", err, "provider", provider)
		app.serverErrorResponse(w, r, err)
		return
	}

	// Create a session
	session, _ := app.sessionStore.Get(r, "auth-session")
	session.Values["user_id"] = user.UserID
	session.Values["provider"] = provider
	session.Values["authenticated"] = true
	err = session.Save(r, w)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) authProviderHandler(w http.ResponseWriter, r *http.Request) {
	provider, err := app.readProviderParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	if !app.IsProviderEnabled(provider) {
		app.notFoundResponse(w, r)
		return
	}

	// Check if there's an active session
	session, _ := app.sessionStore.Get(r, "auth-session")
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		userID := session.Values["user_id"]
		app.logger.Info("user already authenticated", "provider", provider, "user_id", userID)

		// You might want to verify the session or refresh user data here

		err = app.writeJSON(w, http.StatusOK, envelope{"message": "Already authenticated", "user_id": userID}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	r = app.setProviderContext(r, provider)
	// Begin the auth process
	gothic.BeginAuthHandler(w, r)
}

func (app *application) authLogoutHandler(w http.ResponseWriter, r *http.Request) {
	provider, err := app.readProviderParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	app.logger.Info("processing logout", "provider", provider)
	// Clear the session
	session, _ := app.sessionStore.Get(r, "auth-session")
	session.Values = map[interface{}]interface{}{}
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Set the provider in the context
	r = app.setProviderContext(r, provider)

	gothic.Logout(w, r)

	// Send JSON response instead of redirect
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Successfully logged out"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// Uses context to pass provider information through the request context for gothic middleware
func (app *application) setProviderContext(r *http.Request, provider string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), gothic.ProviderParamKey, provider))
}

// InitProviders initializes all configured auth providers
func (app *application) InitProviders() error {
	providers := []goth.Provider{}

	for providerName, providerConfig := range app.config.auth.Providers {
		provider, err := app.createProvider(providerName, providerConfig)
		if err != nil {
			return fmt.Errorf("failed to create provider %s: %w", providerName, err)
		}
		if provider != nil {
			providers = append(providers, provider)
		}
	}

	goth.UseProviders(providers...)
	return nil
}

// createProvider creates a specific provider based on configuration
func (app *application) createProvider(name string, config ProviderConfig) (goth.Provider, error) {
	callbackURL := fmt.Sprintf("%s/v1/auth/%s/callback", app.config.auth.BaseCallbackURL, name)

	switch name {
	case "auth0":
		domain, exists := config.ExtraConfig["domain"]
		if !exists {
			return nil, fmt.Errorf("auth0 domain not configured")
		}
		return auth0.New(
			config.ClientID,
			config.ClientSecret,
			callbackURL,
			domain,
		), nil

	case "github":
		return github.New(
			config.ClientID,
			config.ClientSecret,
			callbackURL,
			config.Scopes...,
		), nil

	case "google":
		return google.New(
			config.ClientID,
			config.ClientSecret,
			callbackURL,
			config.Scopes...,
		), nil

	default:
		return nil, fmt.Errorf("unsupported provider: %s", name)
	}
}

// IsProviderEnabled checks if a specific provider is enabled
func (app *application) IsProviderEnabled(name string) bool {
	_, exists := app.config.auth.Providers[name]
	return exists
}

// GetProviderConfig gets the configuration for a specific provider
func (app *application) GetProviderConfig(name string) (ProviderConfig, bool) {
	config, exists := app.config.auth.Providers[name]
	return config, exists
}
