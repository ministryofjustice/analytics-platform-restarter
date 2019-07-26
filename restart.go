package main

import (
	"encoding/json"
	"errors"
	"io"
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

// NewFromReader parses a RestartRequest from a reader
func NewFromReader(r io.Reader) (*RestartRequest, error) {
	req := &RestartRequest{}
	err := json.NewDecoder(r).Decode(&req)
	if err != nil {
		return req, Error{Code: DecodeError, Err: err}
	}

	if req.Host == "" {
		return req, Error{Code: BlankHostError, Err: errors.New("invalid request: 'host' can't be blank")}
	}
	if req.Reason == "" {
		req.Reason = defaultReason
	}

	return req, nil
}

// POST /restart
// {
//   "host": "example.com",
//   "reason": "data-updated"
// }
func restart(w http.ResponseWriter, r *http.Request) {
	req, err := NewFromReader(r.Body)
	if err != nil {
		switch err.(Error).Code {
		case DecodeError:
			logger.Printf("failed to decode request body: %s", err)
			writeJSONResponse(w, http.StatusBadRequest, "error", "Failed to decode request body: Invalid JSON.")
			return
		case BlankHostError:
			errMsg := "invalid request: 'host' can't be empty"
			logger.Printf(errMsg)
			writeJSONResponse(w, http.StatusBadRequest, "error", errMsg)
			return
		default:
			logger.Printf("failed to restart request for unknown reason: %s", err)
			writeJSONResponse(w, http.StatusInternalServerError, "error", "Unknown error: Failed to process your request.")
			return
		}
	}

	logger.Printf("%s: restart request received. Reason: '%s'", req.Host, req.Reason)

	err = req.Process()
	if err != nil {
		switch err.(Error).Code {
		case K8sError:
			logger.Printf("%s: request to k8s failed: %s", req.Host, err)
			writeJSONResponse(w, http.StatusInternalServerError, "error", "Error while trying to restart application.")
			return
		case NoDeploymentError:
			logger.Printf("%s: no deployments found: %s", req.Host, err)
			writeJSONResponse(w, http.StatusNotFound, "error", "Application not found.")
			return
		case TooManyDeploymentsError:
			logger.Printf("%s: too many deployments found: %s", req.Host, err)
			writeJSONResponse(w, http.StatusConflict, "error", "Too many Deployments found for application, can't process restart request.")
			return
		}
	}

	logger.Printf("%s: pods restarted", req.Host)
	writeJSONResponse(w, http.StatusOK, "message", "Restarted.")
}

// Process restart the requested app with the given host and reason
func (r *RestartRequest) Process() error {
	deploy, err := GetDeployment(r.Host, namespace)
	if err != nil {
		return err
	}

	logger.Printf("%s: Deployment '%s' found.", r.Host, deploy.Name)

	return RestartPods(deploy, r.Reason)
}
