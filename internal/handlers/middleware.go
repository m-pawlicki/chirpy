package handlers

import (
	"net/http"
)

func (apiCfg *APIHandler) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCfg.Config.FileserverHits.Add(int32(1))
		next.ServeHTTP(w, r)
	})
}
