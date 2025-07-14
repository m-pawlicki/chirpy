package handlers

import (
	"encoding/json"
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

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	fmt.Printf("Error: %s", msg)
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Write(data)
}
