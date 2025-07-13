package main

import "net/http"

func (app *application) healthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("OK"))
	if err != nil {
		return
	}
}
