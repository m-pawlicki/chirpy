package handlers

import "net/http"

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	RespondWithMsg(w, 200, "OK")
}
