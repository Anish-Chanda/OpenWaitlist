package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/anish-chanda/openwaitlist/backend/internal/db"
	postgres "github.com/anish-chanda/openwaitlist/backend/internal/db/postgresql"
	"github.com/anish-chanda/openwaitlist/backend/internal/handlers"
	"github.com/anish-chanda/openwaitlist/backend/internal/logger"
	"github.com/anish-chanda/openwaitlist/web"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	authpkg "github.com/go-pkgz/auth/v2"
	"github.com/go-pkgz/auth/v2/avatar"
	"github.com/go-pkgz/auth/v2/provider"
	"github.com/go-pkgz/auth/v2/token"
)

func main() {
	// Load and validate config
	cfg := LoadConfig()

	logConfig := logger.Config{
		Level:       cfg.LogLevel,
		Environment: cfg.Environment,
		ServiceName: "api",
	}

	log := logger.New(logConfig)

	log.Info("Initializng Database...")
	var database db.Database = postgres.NewPostgresDB(*log)
	if err := database.Connect(cfg.Dsn); err != nil {
		log.Error("Database connection failed: ", err)
		return
	}
	// ping DB
	if err := database.Ping(context.Background()); err != nil {
		log.Error("Database ping failed: ", err)
		return
	}

	// run migrations
	log.Info("Running database migrations...")
	if err := database.Migrate(); err != nil {
		log.Error("Database migration failed: ", err)
		return
	}

	// setup auth options
	authOptions := authpkg.Opts{
		SecretReader: token.SecretFunc(func(id string) (string, error) { // secret key for JWT
			return cfg.JWTSecret, nil
		}),
		TokenDuration:  time.Duration(cfg.TokenDuration) * time.Minute, // token expires in X minutes
		CookieDuration: time.Duration(cfg.CookieDuration) * time.Hour,  // cookie expires in X hours
		Issuer:         "openwaitlist",
		URL:            cfg.APIBaseURL,
		DisableXSRF:    true,
		AvatarStore:    avatar.NewLocalFS(cfg.AvatarPath),
	}

	// create authservice and local provider
	authService := authpkg.NewService(authOptions)
	authService.AddDirectProvider("local", provider.CredCheckerFunc(func(user, password string) (ok bool, err error) {
		return handlers.HandleLogin(database, user, password)
	}))

	// create router and attach paths
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	// API routes
	router.Post("/signup", handlers.SignupHandler(database, *log))

	// Auth and avatar handlers
	authHandler, avatarHandler := authService.Handlers()
	router.Mount("/auth", authHandler)
	router.Mount("/avatar", avatarHandler)

	// create file server to serve static frontend files
	fileServer := http.FileServer(http.FS(web.DistDirFS))

	// Wrapper to handle SPA routing - serves index.html for non-asset requests
	spaHandler := func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")

		// Check if it's a static asset (has file extension)
		if strings.Contains(path, ".") {
			// It's a file, serve it directly
			fileServer.ServeHTTP(w, r)
			return
		}

		// It's a route, serve index.html for SPA routing
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	}

	// Handle static files for specific routes
	router.Get("/login", spaHandler)
	router.Get("/signup", spaHandler)
	router.Get("/dashboard", spaHandler)

	// Handle static assets (CSS, JS, etc.)
	router.Get("/chunk-*", func(w http.ResponseWriter, r *http.Request) {
		fileServer.ServeHTTP(w, r)
	})

	// Root path redirects to login (handled by React Router, but serve the app)
	router.Get("/", spaHandler) // Start listening
	log.Info(fmt.Sprintf("Starting server on port %d", cfg.ApiPort))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.ApiPort), router); err != nil {
		log.Error("Server failed to start", err)
	}
}
