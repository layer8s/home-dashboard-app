package main

import (
	"database/sql"
	"net/http"

	"github.com/Robert-litts/fantasy-football-archive/internal/data"
	"github.com/Robert-litts/fantasy-football-archive/internal/db"
	"github.com/Robert-litts/fantasy-football-archive/internal/validator"
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

	leagueParams := db.GetLeaguesParams{
		ID:          input.Id,
		LeagueId:    input.LeagueId,
		Year:        input.Year,
		TeamCount:   input.TeamCount,
		CurrentWeek: input.CurrentWeek,
		NflWeek:     input.NflWeek,
		Limit:       int32(input.Filters.PageSize),
		Offset:      int32((input.Filters.Page - 1) * input.Filters.PageSize),
	}

	app.logger.Info("League Params", leagueParams)

	leagues, err := app.queries.GetLeagues(r.Context(), leagueParams)
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
