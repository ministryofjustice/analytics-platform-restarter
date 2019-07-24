package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// NewServer returns the HTTP server configured to respond to POST /restart
// and GET /healthz
func NewServer(port int) *http.Server {
	m := mux.NewRouter()
	m.HandleFunc("/healthz", healthz).Methods("GET")
	m.HandleFunc("/restart", restart).Methods("POST")

	s := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: m,

		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	return &s
}

func writeJSONResponse(w http.ResponseWriter, status int, messageKey string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		messageKey: message,
	})
}
