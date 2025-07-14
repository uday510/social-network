package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("%s: %s %s - %s", http.StatusText(http.StatusInternalServerError), r.Method, r.URL.Path, err.Error())
	_ = writeJSONError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("%s: %s %s - %s", http.StatusText(http.StatusBadRequest), r.Method, r.URL.Path, err.Error())
	_ = writeJSONError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("%s: %s %s - %s", http.StatusText(http.StatusNotFound), r.Method, r.URL.Path, err.Error())
	_ = writeJSONError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
}
