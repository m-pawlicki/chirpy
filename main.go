package main

import (
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	strBytes := []byte("OK")
	w.Write(strBytes)
}

func main() {
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	mux.HandleFunc("/healthz", handler)
	appReqPath := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", appReqPath)
	server.ListenAndServe()
}
