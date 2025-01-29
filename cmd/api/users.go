package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/Robert-litts/fantasy-football-archive/internal/data"
	"github.com/Robert-litts/fantasy-football-archive/internal/db"
	"github.com/Robert-litts/fantasy-football-archive/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	// Create an anonymous struct to hold the expected data from the request body.
	var input struct {
		Name     *string `json:"name"`
		Email    *string `json:"email"`
		Password *string `json:"password"`
	}

	var (
		ErrDuplicateEmail = errors.New("duplicate email")
	)

	// Parse the request body into the anonymous struct.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy the data from the request body into a new User struct. Notice also that we
	// set the Activated field to false, which isn't strictly necessary because the
	// Activated field will have the zero-value of false by default. But setting this
	// explicitly helps to make our intentions clear to anyone reading the code.
	user := &data.User{
		Name:      *input.Name,
		Email:     *input.Email,
		Activated: false,
	}

	// Use the Password.Set() method to generate and store the hashed and plaintext
	// passwords.
	err = user.Password.Set(*input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	data.ValidateUser(v, user)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	userParams := &db.InsertUserParams{
		Name:         *input.Name,
		Email:        *input.Email,
		Activated:    false,
		PasswordHash: []byte(*user.Password.Hash()),
	}

	// Insert the user data into the database.
	userRow, err := app.queries.InsertUser(r.Context(), *userParams)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
			app.serverErrorResponse(w, r, ErrDuplicateEmail)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	user.ID = userRow.ID
	user.CreatedAt = userRow.CreatedAt

	tokenService := data.NewTokenService(app.queries)

	token, err := tokenService.New(r.Context(), user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.background(func() {

		data := map[string]any{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}
		app.logger.Info("User data info: ", data)

		err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			app.logger.Error((err.Error()))
		}

	})

	// Write a JSON response containing the user data along with a 201 Created status
	// code.
	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
