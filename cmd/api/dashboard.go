package main

import (
	"net/http"
)

func (app *application) dashboardHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := app.sessionStore.Get(r, "auth-session")

	data := map[string]interface{}{
		"Email":    session.Values["email"],
		"Name":     session.Values["name"],
		"Provider": session.Values["provider"],
	}

	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// Render just the user info partial
		err := app.renderTemplate(w, "dashboard-partial.tmpl", data)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Render the full page for regular requests
	err := app.renderTemplate(w, "dashboard.tmpl", data)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
