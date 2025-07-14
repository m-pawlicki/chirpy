package handlers

import (
	"net/http"
)

func (apiCfg *APIHandler) ResetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	apiCfg.Config.FileserverHits.Store(0)
	w.Write([]byte("Hits reset to 0"))
	w.WriteHeader(http.StatusOK)
}
