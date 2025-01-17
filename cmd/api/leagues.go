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
		Id          int
		LeagueId    int
		Year        int
		TeamCount   int
		CurrentWeek int
		NflWeek     int
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	// Read the page, page_size, and sort query string values into the embedded struct.
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "nflWeek", "year", "currentWeek", "teamCount", "-id", "-year", "-teamCount"}

	input.CurrentWeek = app.readInt(qs, "currentWeek", 0, v)
	input.NflWeek = app.readInt(qs, "nflWeek", 0, v)
	input.Year = app.readInt(qs, "year", 0, v)
	input.TeamCount = app.readInt(qs, "teamCount", 0, v)
	input.Id = app.readInt(qs, "id", 0, v)

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	leagueParams := db.GetLeaguesParams{
		Year: int32(input.Year),
		TeamCount: sql.NullInt32{
			Int32: int32(input.TeamCount),
			Valid: input.TeamCount != 0,
		},
		CurrentWeek: sql.NullInt32{
			Int32: int32(input.CurrentWeek),
			Valid: input.CurrentWeek != 0,
		},
		NflWeek: sql.NullInt32{
			Int32: int32(input.NflWeek),
			Valid: input.NflWeek != 0,
		},
		Column5: input.Filters.Sort, // Sort Column
	}

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
