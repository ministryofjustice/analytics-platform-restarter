package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

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
