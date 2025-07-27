package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/uday510/go-crud-app/internal/store/cache"

	"github.com/uday510/go-crud-app/internal/auth"

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

//	@title			GopherSocial API
//	@description	API for SocialNetwork,
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
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
		redisCfg: redisConfig{
			addr:    env.GetString("REDIS_ADDR", "localhost:6379"),
			pw:      env.GetString("REDIS_PW", ""),
			db:      env.GetInt("REDIS_DB", 0),
			enabled: env.GetBool("REDIS_ENABLED", false),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			expiry:    time.Hour * 24 * 3,
			fromEmail: env.GetString("FROM_EMAIL", ""),
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
		},
		auth: authConfig{
			basic: basicConfig{
				user: env.GetString("AUTH_BASIC_USER", "admin"),
				pass: env.GetString("AUTH_BASIC_PASS", "admin"),
			},
			token: tokenConfig{
				secret: env.GetString("AUTH_TOKEN_SECRET", "fallback"),
				exp:    time.Hour * 24 * 3,
				iss:    "socialnetwork",
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

	var redisDB *redis.Client
	if cfg.redisCfg.enabled {
		redisDB = cache.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.pw, cfg.redisCfg.db)
		logger.Info("redis cache connection established")
	}

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
	cacheStorage := cache.NewRedisStorage(redisDB)

	newMailer := mailer.NewSendgrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)

	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.iss, cfg.auth.token.iss)

	app := &application{
		config:        cfg,
		store:         newStorage,
		cacheStorage:  cacheStorage,
		logger:        logger,
		mailer:        newMailer,
		authenticator: jwtAuthenticator,
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
