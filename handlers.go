package main

import (
	"encoding/json"
	"net/http"
)

func healthz(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, http.StatusOK, "message", "OK 👍🏼")
}

func restart(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, http.StatusNotImplemented, "error", "TODO ⚠️")
}

func writeJSONResponse(w http.ResponseWriter, status int, messageKey string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		messageKey: message,
	})
}
