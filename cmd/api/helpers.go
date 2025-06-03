package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/layer8s/home-dashboard-app/internal/validator"
)

type envelope map[string]any

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid ID parameter")
	}
	return id, nil
}

func (app *application) readProviderParam(r *http.Request) (string, error) {
	params := httprouter.ParamsFromContext(r.Context())
	provider := params.ByName("provider")

	if provider == "" {
		return "", errors.New("no provider specified")
	}
	return provider, nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body must not be empty")
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)
		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	js = append(js, '\n')
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	return s
}
func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)
	if csv == "" {
		return defaultValue
	}
	return strings.Split(csv, ",")
}

func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}
	return i
}

func (app *application) readIntQuery(qs url.Values, key string, v *validator.Validator) int32 {
	s := qs.Get(key)
	if s == "" {
		return -1 // Sentinel value to indicate "no value"
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return -1 // Return sentinel value on error
	}
	return int32(i) // Return the valid int32 value
}

func loadEnvironment() (int, string, string, int, int, time.Duration, string, string, string, string, int) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("PORT environment variable not set")
	}
	env := os.Getenv("ENV")
	if env == "" {
		log.Fatal("ENV environment variable not set")
	}

	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("DB_URL environment variable not set")
	}
	dbMaxOpenConns, err := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNS"))
	if err != nil {
		log.Fatal("DB_MAX_OPEN_CONNS environment variable not set")
	}
	dbMaxIdleConns, err := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS"))
	if err != nil {
		log.Fatal("DB_MAX_IDLE_CONNS environment variable not set")
	}
	dbMaxIdleTime, err := time.ParseDuration(os.Getenv("DB_MAX_IDLE_TIME"))
	if err != nil {
		log.Fatal("DB_MAX_IDLE_TIME environment variable not set")
	}

	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		log.Fatal("SESSION_KEY must be set")
	}

	sendGridKey := os.Getenv("SENDGRID_API_KEY")
	if sendGridKey == "" {
		log.Fatal("SENDGRID_API_KEY environment variable is required")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatal("Failed to convert redisDB to int")
	}
	// if redisAddr == "" || redisPassword == "" {
	// 	log.Fatal("REDIS environment variables are required")
	// }

	//return the env variables
	return port, env, dsn, dbMaxOpenConns, dbMaxIdleConns, dbMaxIdleTime, sessionKey, sendGridKey, redisAddr, redisPassword, redisDB
}

// The background() helper accepts an arbitrary function as a parameter.
func (app *application) background(fn func()) {
	// Increment the WaitGroup counter.
	app.wg.Add(1)
	go func() {
		// Recover any panic.

		defer app.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				app.logger.Error(fmt.Sprintf("%v", err))
			}
		}()

		// Execute the arbitrary function that we passed as the parameter.
		fn()
	}()
}

func (app *application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid or missing authentication token"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *application) renderTemplate(w http.ResponseWriter, name string, data interface{}) error {
	templates := make(map[string]*template.Template)

	// Define the paths for different template types
	templatePaths := []string{
		filepath.Join("ui", "html", "base.tmpl"),
		filepath.Join("ui", "html", name),
	}

	// If we're rendering the dashboard or leagues, we also need their respective partial templates
	switch name {
	case "dashboard.tmpl", "dashboard-partial.tmpl":
		templatePaths = append(templatePaths, filepath.Join("ui", "html", "dashboard-partial.tmpl"))
	case "leagues.tmpl", "leagues-partial.tmpl":
		templatePaths = append(templatePaths, filepath.Join("ui", "html", "leagues-partial.tmpl"))
	}

	// Log the templates we're attempting to parse
	app.logger.Info("parsing templates",
		"paths", templatePaths)

	// Parse all required templates
	tmpl, err := template.ParseFiles(templatePaths...)
	if err != nil {
		// Log the error with detailed information
		app.logger.Error("template parsing failed",
			"error", err,
			"paths", templatePaths)
		return fmt.Errorf("error parsing template files: %w", err)
	}

	// Store the template in our cache
	templates[name] = tmpl

	// Execute the appropriate template based on the type
	switch name {
	case "dashboard-partial.tmpl":
		err = tmpl.ExecuteTemplate(w, "user-info", data)
	case "leagues-partial.tmpl":
		err = tmpl.ExecuteTemplate(w, "leagues-table", data)
	default:
		// For full pages, we execute the base template
		err = tmpl.ExecuteTemplate(w, "base", data)
	}

	if err != nil {
		app.logger.Error("template execution failed",
			"error", err,
			"template", name)
		return fmt.Errorf("error executing template: %w", err)
	}

	return nil
}
