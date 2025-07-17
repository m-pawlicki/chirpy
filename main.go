package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/m-pawlicki/chirpy/internal/config"
	"github.com/m-pawlicki/chirpy/internal/database"
	"github.com/m-pawlicki/chirpy/internal/handlers"
	"github.com/m-pawlicki/chirpy/internal/middleware"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load(".env")
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("SECRET")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Error opening database: %s", err)
		return
	}
	dbQueries := database.New(db)

	apiCfg := &config.APIConfig{
		DB:       dbQueries,
		Platform: platform,
		Secret:   secret,
	}

	cfgHandlers := handlers.NewAPIHandler(apiCfg)
	cfgMiddleware := middleware.NewAPIMiddleware(apiCfg)

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	mux.HandleFunc("GET /admin/metrics", cfgHandlers.HitHandler)
	mux.HandleFunc("POST /admin/reset", cfgHandlers.DeleteUsersHandler)

	mux.HandleFunc("GET /api/healthz", handlers.HealthHandler)

	mux.HandleFunc("POST /api/chirps", cfgHandlers.PostChirpHandler)
	mux.HandleFunc("GET /api/chirps", cfgHandlers.GetChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfgHandlers.GetChirpByIDHandler)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfgHandlers.DeleteChirpByIDHandler)

	mux.HandleFunc("POST /api/users", cfgHandlers.CreateUserHandler)
	mux.HandleFunc("PUT /api/users", cfgHandlers.UpdateUserHandler)

	mux.HandleFunc("POST /api/login", cfgHandlers.LoginUserHandler)
	mux.HandleFunc("POST /api/refresh", cfgHandlers.RefreshAccessTokenHandler)
	mux.HandleFunc("POST /api/revoke", cfgHandlers.RevokeTokenHandler)

	mux.HandleFunc("POST /api/polka/webhooks", cfgHandlers.UpgradeUserToRedHandler)

	appReqPath := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", cfgMiddleware.MiddlewareMetricsInc(appReqPath))

	server.ListenAndServe()
}
