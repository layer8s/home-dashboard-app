package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)                 // 404
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse) // 405

	router.HandlerFunc(http.MethodGet, "/", app.loginHandler)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/leagues", app.listLeaguesHandler)
	router.HandlerFunc(http.MethodGet, "/v1/leagues/:id", app.showLeagueHandler)
	router.HandlerFunc(http.MethodGet, "/v1/leagues/:id/teams/:id", app.showTeamHandler)
	router.HandlerFunc(http.MethodGet, "/v1/auth/:provider/callback", app.HandleCallback)
	router.HandlerFunc(http.MethodGet, "/v1/auth/:provider/logout", app.HandleLogout)
	router.HandlerFunc(http.MethodGet, "/v1/auth/:provider", app.HandleAuth)

	// route with authentication middleware
	router.HandlerFunc(http.MethodGet, "/v1/dashboard",
		app.requireAuthenticated(app.dashboardHandler))

	// route for HTMX to refresh user info
	router.HandlerFunc(http.MethodGet, "/v1/dashboard/refresh",
		app.requireAuthenticated(app.dashboardHandler))

	router.HandlerFunc(http.MethodGet, "/v1/dashboard/leagues",
		app.requireAuthenticated(app.leaguesPageHandler))

	router.HandlerFunc(http.MethodGet, "/v1/dashboard/leagues/refresh",
		app.requireAuthenticated(app.leaguesPageHandler))

	router.HandlerFunc(http.MethodGet, "/v1/dashboard/index",
		app.requireAuthenticated(app.leaguesIndexHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)

	router.HandlerFunc(http.MethodGet, "/login", app.loginTemplHandler)

	router.HandlerFunc(http.MethodGet, "/mm", app.magicMirrorHandler)

	return app.recoverPanic(router)
}
