package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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

//	@title			SocialNetwork API
//	@version		1.0
//	@description	This is a social network server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath					/api/v1
//	@securityDefinitions.basic	BasicAuth

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description

func main() {
	log.Println("loading configuration...")

	cfg := config{
		addr:   env.GetString("ADDR", ":8080"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://user:password@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNECTIONS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNECTIONS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
	}

	log.Printf("configuration loaded: addr=%s, db_addr=%s", cfg.addr, cfg.db.addr)

	log.Println("initializing database connection...")

	database, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Fatalf("failed to create database pool: %v", err)
	}
	log.Println("database connection pool established")

	defer func() {
		log.Println("closing database connection...")
		if err := database.Close(); err != nil {
			log.Printf("error closing database connection: %v", err)
		} else {
			log.Println("database connection closed")
		}
	}()

	log.Println("initializing storage layer...")
	newStorage := store.NewStorage(database)

	app := &application{
		config: cfg,
		store:  newStorage,
	}

	// Graceful shutdown handling
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		sig := <-stop
		log.Printf("received signal: %s. initiating shutdown...", sig)
		os.Exit(0)
	}()

	mux := app.mount()

	log.Printf("starting HTTP server on %s", cfg.addr)
	log.Fatal(app.run(mux))
}
