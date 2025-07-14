package main

import (
	"net/http"

	"github.com/m-pawlicki/chirpy/handlers"
	"github.com/m-pawlicki/chirpy/internal/config"
)

func main() {
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	apiCfg := &config.APIConfig{}
	apiHandlers := handlers.NewAPIHandler(apiCfg)

	mux.HandleFunc("GET /api/healthz", handlers.HealthHandler)
	mux.HandleFunc("GET /admin/metrics", apiHandlers.HitHandler)
	mux.HandleFunc("POST /admin/reset", apiHandlers.ResetHandler)

	appReqPath := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", apiHandlers.MiddlewareMetricsInc(appReqPath))

	server.ListenAndServe()
}
