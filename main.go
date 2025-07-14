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

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load(".env")
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Error opening database: %s", err)
		return
	}

	dbQueries := database.New(db)
	apiCfg := &config.APIConfig{DB: dbQueries}
	cfgHandlers := handlers.NewAPIHandler(apiCfg)

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	mux.HandleFunc("GET /api/healthz", handlers.HealthHandler)
	mux.HandleFunc("GET /admin/metrics", cfgHandlers.HitHandler)
	mux.HandleFunc("POST /admin/reset", cfgHandlers.ResetHandler)
	mux.HandleFunc("POST /api/validate_chirp", cfgHandlers.ValidateHandler)

	appReqPath := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", cfgHandlers.MiddlewareMetricsInc(appReqPath))

	server.ListenAndServe()
}
