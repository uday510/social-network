package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw(http.StatusText(http.StatusInternalServerError),
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error(),
	)
	_ = writeJSON(w, http.StatusInternalServerError, err.Error())
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw(http.StatusText(http.StatusBadRequest),
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error(),
	)
	_ = writeJSON(w, http.StatusBadRequest, err.Error())
}

func (app *application) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw(http.StatusText(http.StatusConflict),
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error(),
	)
	_ = writeJSON(w, http.StatusConflict, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw(http.StatusText(http.StatusNotFound),
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error(),
	)
	_ = writeJSON(w, http.StatusNotFound, err.Error())
}
