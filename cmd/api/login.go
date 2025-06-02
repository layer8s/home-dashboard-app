package main

import (
	"net/http"

	"github.com/layer8s/home-dashboard-app/templates"
)

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	loginPage := templates.Login()
	err := loginPage.Render(r.Context(), w)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
