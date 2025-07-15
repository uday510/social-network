package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("%s: %s %s - %s", http.StatusText(http.StatusInternalServerError), r.Method, r.URL.Path, err.Error())
	_ = writeJSON(w, http.StatusInternalServerError, err.Error())
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("%s: %s %s - %s", http.StatusText(http.StatusBadRequest), r.Method, r.URL.Path, err.Error())
	_ = writeJSON(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("%s: %s %s - %s", http.StatusText(http.StatusNotFound), r.Method, r.URL.Path, err.Error())
	_ = writeJSON(w, http.StatusNotFound, err.Error())
}
