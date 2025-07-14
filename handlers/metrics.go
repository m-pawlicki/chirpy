package handlers

import (
	"fmt"
	"net/http"

	"github.com/m-pawlicki/chirpy/internal/config"
)

type APIHandler struct {
	Config *config.APIConfig
}

func NewAPIHandler(cfg *config.APIConfig) *APIHandler {
	return &APIHandler{
		Config: cfg,
	}
}

func (apiCfg *APIHandler) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCfg.Config.FileserverHits.Add(int32(1))
		next.ServeHTTP(w, r)
	})
}

func (apiCfg *APIHandler) HitHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	hits := int(apiCfg.Config.FileserverHits.Load())
	resp := "<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>"
	fmt.Fprintf(w, resp, hits)
}
