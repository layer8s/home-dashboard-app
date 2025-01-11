package main

import (
	"net/http"

	"github.com/Robert-litts/fantasy-football-archive/internal/data"
)

func (app *application) showLeagueHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	league := data.League{
		Id:          id,
		LeagueId:    12345,
		Year:        2013,
		TeamCount:   8,
		CurrentWeek: 15,
		NflWeek:     0,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"league": league}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
