package main

import (
	"net/http"

	"github.com/layer8s/home-dashboard-app/templates"
)

// func (app *application) showLeagueHandler(w http.ResponseWriter, r *http.Request) {
// 	id, err := app.readIDParam(r)
// 	if err != nil {
// 		app.notFoundResponse(w, r)
// 		return
// 	}

// 	app.logger.Info("attempting to fetch league", "id", id)

// 	// Use the SQLC-generated query method
// 	league, err := app.queries.GetLeagueById(r.Context(), int32(id))
// 	if err != nil {
// 		app.logger.Error("database error", "error", err)
// 		if err == sql.ErrNoRows {
// 			app.notFoundResponse(w, r)
// 			return
// 		}
// 		app.serverErrorResponse(w, r, err)
// 		return
// 	}

// 	err = app.writeJSON(w, http.StatusOK, envelope{"league": league}, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}
// }

// func (app *application) listLeaguesHandler(w http.ResponseWriter, r *http.Request) {
// 	var input struct {
// 		Id          int32
// 		LeagueId    int32
// 		Year        int32
// 		TeamCount   int32
// 		CurrentWeek int32
// 		NflWeek     int32
// 		data.Filters
// 	}

// 	v := validator.New()
// 	qs := r.URL.Query()

// 	// Read the page, page_size, and sort query string values into the embedded struct.
// 	input.Filters.Page = app.readInt(qs, "page", 1, v)
// 	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
// 	input.Filters.Sort = app.readString(qs, "sort", "id")
// 	input.Filters.SortSafelist = []string{"id", "nflWeek", "year", "currentWeek", "teamCount", "-id", "-year", "-teamCount"}

// 	input.CurrentWeek = app.readIntQuery(qs, "currentWeek", v)
// 	input.NflWeek = app.readIntQuery(qs, "nflWeek", v)
// 	input.Year = app.readIntQuery(qs, "year", v)
// 	input.TeamCount = app.readIntQuery(qs, "teamCount", v)
// 	input.Id = app.readIntQuery(qs, "id", v)
// 	input.LeagueId = app.readIntQuery(qs, "leagueId", v)

// 	if data.ValidateFilters(v, input.Filters); !v.Valid() {
// 		app.failedValidationResponse(w, r, v.Errors)
// 		return
// 	}
// 	// Sort the leagues by the specified column and direction.
// 	sortVal := input.Filters.SortColumn()
// 	sortDir := input.Filters.SortDirection()

// 	// Set up the query parameters based on the filters, convert based on sort Direction.
// 	baseParams := db.GetLeaguesAscParams{
// 		ID:          input.Id,
// 		LeagueId:    input.LeagueId,
// 		Year:        input.Year,
// 		TeamCount:   input.TeamCount,
// 		CurrentWeek: input.CurrentWeek,
// 		NflWeek:     input.NflWeek,
// 		Limit:       int32(input.Filters.PageSize),
// 		Offset:      int32((input.Filters.Page - 1) * input.Filters.PageSize),
// 		Column9:     sortVal,
// 	}

// 	app.logger.Info("League Params", baseParams)
// 	var leagues []db.League
// 	var err error

// 	// Check the sort direction and call the appropriate method.
// 	if sortDir == "DESC" {
// 		descParams := db.GetLeaguesDescParams(baseParams)
// 		leagues, err = app.queries.GetLeaguesDesc(r.Context(), descParams)
// 	} else {
// 		leagues, err = app.queries.GetLeaguesAsc(r.Context(), baseParams)
// 	}

// 	if err != nil {
// 		app.logger.Error("database error", "error", err)
// 		if err == sql.ErrNoRows {
// 			app.notFoundResponse(w, r)
// 			return
// 		}
// 		app.serverErrorResponse(w, r, err)
// 		return
// 	}

// 	err = app.writeJSON(w, http.StatusOK, envelope{"leagues": leagues}, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}
// }

// leaguesPageHandler renders the leagues page for authenticated users
func (app *application) magicMirrorHandler(w http.ResponseWriter, r *http.Request) {
	// Get user session data
	// session, _ := app.sessionStore.Get(r, "auth-session")

	// Set up the query parameters for getting leagues

	// Render the full page inside the base layout
	templates.Base(templates.MagicMirror()).Render(r.Context(), w)
}

// // leaguesRefreshHandler handles HTMX requests to refresh the leagues table
// func (app *application) leaguesRefreshHandler(w http.ResponseWriter, r *http.Request) {
// 	baseParams := db.GetLeaguesAscParams{
// 		Limit:       100,
// 		Offset:      0,
// 		Column9:     "id",
// 		ID:          -1,
// 		LeagueId:    -1,
// 		Year:        -1,
// 		TeamCount:   -1,
// 		CurrentWeek: -1,
// 		NflWeek:     -1,
// 	}

// 	leagues, err := app.queries.GetLeaguesAsc(r.Context(), baseParams)
// 	if err != nil {
// 		app.logger.Error("database error", "error", err)
// 		if err == sql.ErrNoRows {
// 			app.notFoundResponse(w, r)
// 			return
// 		}
// 		app.serverErrorResponse(w, r, err)
// 		return
// 	}

// 	data := map[string]interface{}{
// 		"Leagues": leagues,
// 	}

// 	err = app.renderTemplate(w, "leagues-partial.tmpl", data)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}
// }

// func (app *application) leaguesIndexHandler(w http.ResponseWriter, r *http.Request) {
// 	component := templates.LeaguesPage()
// 	component.Render(r.Context(), w)
// }
