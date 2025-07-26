package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/uday510/go-crud-app/internal/mailer"

	"go.uber.org/zap"

	"github.com/go-playground/validator/v10"
	_ "github.com/swaggo/http-swagger/v2"
	"github.com/uday510/go-crud-app/internal/db"
	"github.com/uday510/go-crud-app/internal/env"
	"github.com/uday510/go-crud-app/internal/store"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

const version = "0.0.1"

// Swagger annotations omitted for brevity...

func main() {
	log.Println("loading configuration...")

	cfg := config{
		addr:        env.GetString("ADDR", ":8080"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:4000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://user:password@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNECTIONS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNECTIONS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			expiry:    time.Hour * 24 * 3,
			fromEmail: env.GetString("FROM_EMAIL", ""),
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
		},
	}

	// Use color logger in development
	var zapLogger *zap.Logger
	if cfg.env == "development" {
		zapLogger = NewDevLogger()
	} else {
		zapLogger = zap.Must(zap.NewProduction())
	}
	logger := zapLogger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	logger.Infow("configuration loaded", "addr", cfg.addr, "db_addr", cfg.db.addr)

	logger.Info("initializing database connection...")
	database, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatalf("failed to create database pool: %v", err)
	}
	logger.Info("database connection pool established")

	defer func() {
		logger.Info("closing database connection...")
		if err := database.Close(); err != nil {
			logger.Errorf("error closing database connection: %v", err)
		} else {
			logger.Info("database connection closed")
		}
	}()

	logger.Info("initializing storage layer...")
	newStorage := store.NewStorage(database)

	mailtrap := mailer.NewSendgrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)

	app := &application{
		config: cfg,
		store:  newStorage,
		logger: logger,
		mailer: mailtrap,
	}

	mux := app.mount()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-stop
		logger.Infof("received signal: %s. initiating shutdown...", sig)
		if err := database.Close(); err != nil {
			logger.Errorf("error closing database: %v", err)
		}
		os.Exit(0)
	}()

	logger.Infow("starting HTTP server", "addr", cfg.addr, "env", cfg.env)
	logger.Fatal(app.run(mux))
}
