package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// When a restart reason is not specified we assume user wants to restart
// an application because its data changed.
const defaultReason = "data-changed"

// RestartRequest contains the host to restart and the reason for restart.
type RestartRequest struct {
	Host   string `json:"host"`
	Reason string `json:"reason"`
}

// POST /restart
// {
//   "host": "example.com",
//   "reason": "data-updated"
// }
func restart(w http.ResponseWriter, r *http.Request) {
	req, err := getRestartRequest(r.Body)
	if err != nil {
		errMsg := fmt.Sprintf("invalid restart request: %s", err)
		log.Print(errMsg)
		writeJSONResponse(w, http.StatusBadRequest, "error", errMsg)
		return
	}

	logger.Printf("%s: restart request received. Reason: '%s'", req.Host, req.Reason)

	deployment, err := GetDeployment(req.Host)
	if err != nil {
		logger.Printf("%s: failed to get deployment: %s", req.Host, err)
		writeJSONResponse(
			w,
			http.StatusInternalServerError,
			"error",
			fmt.Sprintf("failed to get app with given host"),
		)
		return
	}
	if deployment == nil {
		logger.Printf("%s: deployment not found.", req.Host)
		writeJSONResponse(
			w,
			http.StatusNotFound,
			"error",
			fmt.Sprintf("app with given host not found"),
		)
		return
	}

	logger.Print(deployment)

	writeJSONResponse(w, http.StatusNotImplemented, "error", "TODO ⚠️")
}

func getRestartRequest(r io.Reader) (RestartRequest, error) {
	req := RestartRequest{}
	err := json.NewDecoder(r).Decode(&req)
	if err != nil {
		return req, fmt.Errorf("failed to parse restart request: %s", err)
	}

	if req.Host == "" {
		return req, fmt.Errorf("'host' can't be blank")
	}
	if req.Reason == "" {
		req.Reason = defaultReason
	}

	return req, nil
}

func healthz(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, http.StatusOK, "message", "OK 👍🏼")
}

func writeJSONResponse(w http.ResponseWriter, status int, messageKey string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		messageKey: message,
	})
}
