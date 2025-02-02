package main

import (
	"net/http"

	"github.com/Robert-litts/fantasy-football-archive/templates"
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
	var name, email, provider string
	var namePtr, emailPtr, providerPtr *string

	// Handle name
	if nameVal, ok := session.Values["name"].(string); ok {
		name = nameVal
		namePtr = &name
	}

	// Handle email
	if emailVal, ok := session.Values["email"].(string); ok {
		email = emailVal
		emailPtr = &email
	}

	// Handle provider
	if providerVal, ok := session.Values["provider"].(string); ok {
		provider = providerVal
		providerPtr = &provider
	}

	// // Check if any required values are missing
	// if namePtr == nil || emailPtr == nil || providerPtr == nil {
	//     app.sessionRequiredResponse(w, r)
	//     return
	// }

	// If it's an HTMX request, return just the user info partial
	if r.Header.Get("HX-Request") == "true" {
		err := templates.UserInfo(namePtr, emailPtr, providerPtr).Render(r.Context(), w)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Render the full dashboard inside the base layout
	err = templates.Base(
		templates.Dashboard(namePtr, emailPtr, providerPtr),
	).Render(r.Context(), w)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
