package handlers

import (
	"fmt"
	"net/http"
)

func (apiCfg *APIHandler) HitHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	hits := int(apiCfg.Config.FileserverHits.Load())
	resp := "<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>"
	fmt.Fprintf(w, resp, hits)
}
