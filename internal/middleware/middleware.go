package middleware

import (
	"net/http"

	"github.com/m-pawlicki/chirpy/internal/config"
)

type APIMiddleware struct {
	Config *config.APIConfig
}

func NewAPIMiddleware(cfg *config.APIConfig) *APIMiddleware {
	return &APIMiddleware{
		Config: cfg,
	}
}

func (apiCfg *APIMiddleware) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCfg.Config.FileserverHits.Add(int32(1))
		next.ServeHTTP(w, r)
	})
}
