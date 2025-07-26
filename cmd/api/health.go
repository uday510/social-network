package main

import (
	"net/http"
	"time"
)

// healthcheckHandler godoc
//
//	@Summary		Healthcheck
//	@Description	Healthcheck endpoint
//	@Tags			ops
//	@Produce		json
//	@Success		200	{object}	string	"ok"
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":    "OK",
		"env":       app.config.env,
		"version":   version,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.badRequestError(w, r, err)
	}
}
