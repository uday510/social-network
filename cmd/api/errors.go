package main

import (
	"net/http"
)

func (app *application) internalServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw(http.StatusText(http.StatusInternalServerError),
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error(),
	)
	_ = writeJSON(w, http.StatusInternalServerError, err.Error())
}

func (app *application) badRequestErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw(http.StatusText(http.StatusBadRequest),
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error(),
	)
	_ = writeJSON(w, http.StatusBadRequest, err.Error())
}

func (app *application) conflictErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw(http.StatusText(http.StatusConflict),
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error(),
	)
	_ = writeJSON(w, http.StatusConflict, err.Error())
}

func (app *application) notFoundErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw(http.StatusText(http.StatusNotFound),
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error(),
	)
	_ = writeJSON(w, http.StatusNotFound, err.Error())
}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw(http.StatusText(http.StatusUnauthorized),
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error(),
	)
	_ = writeJSON(w, http.StatusUnauthorized, err.Error())
}

func (app *application) unauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw(http.StatusText(http.StatusUnauthorized),
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error(),
	)
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	_ = writeJSON(w, http.StatusUnauthorized, err.Error())
}

func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	app.logger.Warnw(http.StatusText(http.StatusUnauthorized),
		"method", r.Method,
		"path", r.URL.Path,
	)

	_ = writeJSONError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
}
