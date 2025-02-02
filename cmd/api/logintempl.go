package main

import (
	"net/http"

	"github.com/Robert-litts/fantasy-football-archive/templates"
)

func (app *application) loginTemplHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := app.sessionStore.Get(r, "auth-session")
	name, ok := session.Values["name"].(string)

	component := templates.Hello(name)
	if !ok {
		name = "Guest"
	}

	// Render the component
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
