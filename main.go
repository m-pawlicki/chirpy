package main

import (
	"net/http"

	"github.com/m-pawlicki/chirpy/internal/config"
	"github.com/m-pawlicki/chirpy/internal/handlers"

	_ "github.com/lib/pq"
)

func main() {
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	apiCfg := &config.APIConfig{}
	cfgHandlers := handlers.NewAPIHandler(apiCfg)

	mux.HandleFunc("GET /api/healthz", handlers.HealthHandler)
	mux.HandleFunc("GET /admin/metrics", cfgHandlers.HitHandler)
	mux.HandleFunc("POST /admin/reset", cfgHandlers.ResetHandler)
	mux.HandleFunc("POST /api/validate_chirp", cfgHandlers.ValidateHandler)

	appReqPath := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", cfgHandlers.MiddlewareMetricsInc(appReqPath))

	server.ListenAndServe()
}
