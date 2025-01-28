package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Robert-litts/fantasy-football-archive/internal/db"
	"github.com/Robert-litts/fantasy-football-archive/internal/mailer"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"github.com/markbates/goth/gothic"
)

const version = "1.0.0"

type AuthConfig struct {
	BaseCallbackURL string
	Providers       map[string]ProviderConfig
}

// ProviderConfig holds the configuration for a specific auth provider
type ProviderConfig struct {
	Name         string
	ClientID     string
	ClientSecret string
	Scopes       []string
	ExtraConfig  map[string]string
}

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
	auth       AuthConfig
	sessionKey string
}

type application struct {
	config       config
	logger       *slog.Logger
	queries      *db.Queries
	sessionStore sessions.Store
	mailer       *mailer.Mailer
}

func main() {
	// Load environment variables
	port, env, dsn, dbMaxOpenConns, dbMaxIdleConns, dbMaxIdleTime, sessionKey, sendGridKey := loadEnvironment()

	var cfg config

	//Command line flags, default values set via .env variables
	flag.IntVar(&cfg.port, "port", port, "API Server Port")
	flag.StringVar(&cfg.env, "env", env, "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", dsn, "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", dbMaxOpenConns, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", dbMaxIdleConns, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", dbMaxIdleTime, "PostgreSQL max connection idle time")
	flag.Parse()

	cfg.auth = AuthConfig{
		BaseCallbackURL: fmt.Sprintf("http://localhost:%d", port),
		Providers: map[string]ProviderConfig{
			"auth0": {
				Name:         "auth0",
				ClientID:     os.Getenv("AUTH0_KEY"),
				ClientSecret: os.Getenv("AUTH0_SECRET"),
				ExtraConfig: map[string]string{
					"domain": os.Getenv("AUTH0_DOMAIN"),
				},
			},
			"github": {
				Name:         "github",
				ClientID:     os.Getenv("GITHUB_KEY"),
				ClientSecret: os.Getenv("GITHUB_SECRET"),
				Scopes:       []string{"user:email", "read:user"},
			},
			// Add more providers as needed
		},
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	dbConn, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer dbConn.Close()

	logger.Info("database connection pool established")

	queries := db.New(dbConn)

	cfg.sessionKey = sessionKey

	app := &application{
		config:       cfg,
		logger:       logger,
		queries:      queries,
		sessionStore: sessions.NewCookieStore([]byte(cfg.sessionKey)),
		mailer:       mailer.New(sendGridKey, "FFArchive <robert@litts.org>"),
	}

	// Initialize the auth providers
	if err := app.InitProviders(); err != nil {
		logger.Error("failed to initialize auth providers", "error", err)
		os.Exit(1)
	}
	gothic.Store = app.sessionStore

	// Call app.serve() to start the server.
	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config
	// struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	db.SetConnMaxIdleTime(cfg.db.maxIdleTime)

	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use PingContext() to establish a new connection to the database, passing in the
	// context we created above as a parameter. If the connection couldn't be
	// established successfully within the 5 second deadline, then this will return an
	// error. If we get this error, or any other, we close the connection pool and
	// return the error.
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	// Return the sql.DB connection pool.
	return db, nil
}
