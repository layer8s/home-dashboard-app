package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)                 // 404
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse) // 405

	authHandlers := &AuthHandlers{app: app}

	router.HandlerFunc(http.MethodGet, "/", app.loginHandler)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/leagues", app.listLeaguesHandler)
	router.HandlerFunc(http.MethodGet, "/v1/leagues/:id", app.showLeagueHandler)
	router.HandlerFunc(http.MethodGet, "/v1/leagues/:id/teams/:id", app.showTeamHandler)
	router.HandlerFunc(http.MethodGet, "/v1/auth/:provider/callback", authHandlers.HandleCallback)
	router.HandlerFunc(http.MethodGet, "/v1/auth/:provider/logout", authHandlers.HandleLogout)
	router.HandlerFunc(http.MethodGet, "/v1/auth/:provider", authHandlers.HandleAuth)

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)

	return app.recoverPanic(router)
}
