package main

import (
	"net/http"

	"github.com/Robert-litts/fantasy-football-archive/internal/data"
)

func (app *application) showLeagueHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
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

	err = app.writeJSON(w, http.StatusOK, league, nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
