package main

import (
	"database/sql"
	"net/http"

	"github.com/Robert-litts/fantasy-football-archive/internal/data"
	"github.com/Robert-litts/fantasy-football-archive/internal/db"
	"github.com/Robert-litts/fantasy-football-archive/internal/validator"
	"github.com/Robert-litts/fantasy-football-archive/templates"
)

func (app *application) showLeagueHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	app.logger.Info("attempting to fetch league", "id", id)

	// Use the SQLC-generated query method
	league, err := app.queries.GetLeagueById(r.Context(), int32(id))
	if err != nil {
		app.logger.Error("database error", "error", err)
		if err == sql.ErrNoRows {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"league": league}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listLeaguesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Id          int32
		LeagueId    int32
		Year        int32
		TeamCount   int32
		CurrentWeek int32
		NflWeek     int32
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	// Read the page, page_size, and sort query string values into the embedded struct.
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "nflWeek", "year", "currentWeek", "teamCount", "-id", "-year", "-teamCount"}

	input.CurrentWeek = app.readIntQuery(qs, "currentWeek", v)
	input.NflWeek = app.readIntQuery(qs, "nflWeek", v)
	input.Year = app.readIntQuery(qs, "year", v)
	input.TeamCount = app.readIntQuery(qs, "teamCount", v)
	input.Id = app.readIntQuery(qs, "id", v)
	input.LeagueId = app.readIntQuery(qs, "leagueId", v)

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Sort the leagues by the specified column and direction.
	sortVal := input.Filters.SortColumn()
	sortDir := input.Filters.SortDirection()

	// Set up the query parameters based on the filters, convert based on sort Direction.
	baseParams := db.GetLeaguesAscParams{
		ID:          input.Id,
		LeagueId:    input.LeagueId,
		Year:        input.Year,
		TeamCount:   input.TeamCount,
		CurrentWeek: input.CurrentWeek,
		NflWeek:     input.NflWeek,
		Limit:       int32(input.Filters.PageSize),
		Offset:      int32((input.Filters.Page - 1) * input.Filters.PageSize),
		Column9:     sortVal,
	}

	app.logger.Info("League Params", baseParams)
	var leagues []db.League
	var err error

	// Check the sort direction and call the appropriate method.
	if sortDir == "DESC" {
		descParams := db.GetLeaguesDescParams(baseParams)
		leagues, err = app.queries.GetLeaguesDesc(r.Context(), descParams)
	} else {
		leagues, err = app.queries.GetLeaguesAsc(r.Context(), baseParams)
	}

	if err != nil {
		app.logger.Error("database error", "error", err)
		if err == sql.ErrNoRows {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"leagues": leagues}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// leaguesPageHandler renders the leagues page for authenticated users
func (app *application) leaguesPageHandler(w http.ResponseWriter, r *http.Request) {
	// Get user session data
	// session, _ := app.sessionStore.Get(r, "auth-session")

	// Set up the query parameters for getting leagues
	baseParams := db.GetLeaguesAscParams{
		Limit:       100,
		Offset:      0,
		Column9:     "id",
		ID:          -1,
		LeagueId:    -1,
		Year:        -1,
		TeamCount:   -1,
		CurrentWeek: -1,
		NflWeek:     -1,
	}
	// Get leagues from database
	leagues, err := app.queries.GetLeaguesAsc(r.Context(), baseParams)
	if err != nil {
		app.logger.Error("database error", "error", err)
		app.serverErrorResponse(w, r, err)
		return
	}

	//app.logger.Info("Leagues fetched:", leagues)

	// // Prepare template data
	// data := map[string]interface{}{
	// 	"Email":    session.Values["email"],
	// 	"Name":     session.Values["name"],
	// 	"Provider": session.Values["provider"],
	// 	"Leagues":  leagues,
	// }

	// // Render the template
	// err = app.renderTemplate(w, "leagues.tmpl", data)
	// if err != nil {
	// 	app.serverErrorResponse(w, r, err)
	// }

	// If it's an HTMX request, return just the table
	if r.Header.Get("HX-Request") == "true" {
		err := templates.LeaguesTable(leagues).Render(r.Context(), w)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Render the full page inside the base layout
	err = templates.Base(
		templates.Leagues(leagues),
	).Render(r.Context(), w)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// leaguesRefreshHandler handles HTMX requests to refresh the leagues table
func (app *application) leaguesRefreshHandler(w http.ResponseWriter, r *http.Request) {
	baseParams := db.GetLeaguesAscParams{
		Limit:       100,
		Offset:      0,
		Column9:     "id",
		ID:          -1,
		LeagueId:    -1,
		Year:        -1,
		TeamCount:   -1,
		CurrentWeek: -1,
		NflWeek:     -1,
	}

	leagues, err := app.queries.GetLeaguesAsc(r.Context(), baseParams)
	if err != nil {
		app.logger.Error("database error", "error", err)
		if err == sql.ErrNoRows {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	data := map[string]interface{}{
		"Leagues": leagues,
	}

	err = app.renderTemplate(w, "leagues-partial.tmpl", data)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
