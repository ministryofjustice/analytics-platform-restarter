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

	s := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: m,

		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	return &s
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK 👍🏼"))
}
