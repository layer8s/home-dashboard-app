package main

import (
	"net/http"

	"github.com/Robert-litts/fantasy-football-archive/templates"
)

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	loginPage := templates.Login()
	err := loginPage.Render(r.Context(), w)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
