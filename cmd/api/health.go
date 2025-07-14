package main

import "net/http"

func (app *application) healthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("OK"))

}
