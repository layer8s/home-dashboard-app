package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/Robert-litts/fantasy-football-archive/internal/db"
	"github.com/Robert-litts/fantasy-football-archive/internal/mailer"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type AuthConfig struct {
	BaseCallbackURL string
	Providers       map[string]ProviderConfig
}

// ProviderConfig holds the configuration for a specific auth provider
type ProviderConfig struct {
	Name            string
	ClientID        string
	ClientSecret    string
	Scopes          []string
	BaseCallbackURL string
	ExtraConfig     map[string]string
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
	auth struct {
		baseCallbackURL string
		providers       map[string]ProviderConfig
	}
	sessionKey string
}

type application struct {
	config       config
	logger       *slog.Logger
	queries      *db.Queries
	sessionStore sessions.Store
	mailer       *mailer.Mailer
	wg           sync.WaitGroup
	authManager  *AuthManager
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

	cfg.auth.baseCallbackURL = fmt.Sprintf("http://localhost:%d", port)
	cfg.auth.providers = map[string]ProviderConfig{
		"auth0": {
			Name:            "auth0",
			ClientID:        os.Getenv("AUTH0_CLIENT_ID"),
			ClientSecret:    os.Getenv("AUTH0_CLIENT_SECRET"),
			Scopes:          []string{"openid", "profile", "email"},
			BaseCallbackURL: cfg.auth.baseCallbackURL,
			ExtraConfig: map[string]string{
				"issuer": fmt.Sprintf("https://%s/", os.Getenv("AUTH0_DOMAIN")),
			},
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
		mailer:       mailer.New(sendGridKey, "FFArchive <robert@litts.org>", logger),
		authManager:  NewAuthManager(),
	}

	// Initialize the auth providers
	ctx := context.Background()
	for _, providerConfig := range cfg.auth.providers {
		if err := app.authManager.RegisterProvider(ctx, providerConfig); err != nil {
			logger.Error("failed to register provider", "error", err, "provider", providerConfig.Name)
			os.Exit(1)
		}
	}

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
