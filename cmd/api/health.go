package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	data := map[string]string{
		"status":  "OK",
		"env":     app.config.env,
		"version": version,
	}
	if err := writeJSON(w, http.StatusOK, data); err != nil {
		_ = writeJSONError(w, http.StatusInternalServerError, "err.Error()")
	}
}
