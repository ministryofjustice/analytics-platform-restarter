package main

import "net/http"

func healthz(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, http.StatusOK, "message", "OK 👍🏼")
}
