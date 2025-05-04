package main

import "net/http"

func (app *application) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)
	router.HandleFunc("GET /v1/users/{username}", app.getUserHandler)
	router.HandleFunc("POST /v1/users", app.registerUserHandler)
	return router
}
