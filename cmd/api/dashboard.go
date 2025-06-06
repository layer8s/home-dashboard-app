package main

import (
	"net/http"
	"strconv"

	"github.com/layer8s/home-dashboard-app/internal/db"
	"github.com/layer8s/home-dashboard-app/templates"
)

// func (app *application) dashboardHandler(w http.ResponseWriter, r *http.Request) {
// 	session, _ := app.sessionStore.Get(r, "auth-session")

// 	data := map[string]interface{}{
// 		"Email":    session.Values["email"],
// 		"Name":     session.Values["name"],
// 		"Provider": session.Values["provider"],
// 	}

// 	// Check if this is an HTMX request
// 	if r.Header.Get("HX-Request") == "true" {
// 		// Render just the user info partial
// 		err := app.renderTemplate(w, "dashboard-partial.tmpl", data)
// 		if err != nil {
// 			app.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}

// 	// Render the full page for regular requests
// 	err := app.renderTemplate(w, "dashboard.tmpl", data)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}
// }

func (app *application) dashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Get session and handle potential error
	session, err := app.sessionStore.Get(r, "auth-session")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Extract session values and convert to pointers
	var name, email, provider, sub string
	var hasName, hasEmail bool

	// Handle name
	if nameVal, ok := session.Values["name"].(string); ok {
		name = nameVal
		hasName = true
	}

	// Handle email
	if emailVal, ok := session.Values["email"].(string); ok {
		email = emailVal
		hasEmail = true
	}

	provider = session.Values["provider"].(string)
	sub = session.Values["user_id"].(string)

	// Get user data from database
	user, err := app.queries.GetUserByProvider(r.Context(), db.GetUserByProviderParams{
		Provider:   provider,
		ProviderID: sub,
	})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var namePtr, emailPtr *string
	if hasName {
		namePtr = &name
	}
	if hasEmail {
		emailPtr = &email
	}
	// These always exist due to auth middleware
	providerPtr := &provider
	subPtr := &sub

	// If it's an HTMX request, return just the user info partial
	if r.Header.Get("HX-Request") == "true" {
		userIDStr := strconv.FormatInt(user.ID, 10)
		err := templates.UserInfo(&user.Name, emailPtr, providerPtr, &userIDStr).Render(r.Context(), w)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Render the full dashboard inside the base layout
	err = templates.Base(
		templates.Dashboard(namePtr, emailPtr, providerPtr, subPtr),
	).Render(r.Context(), w)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
