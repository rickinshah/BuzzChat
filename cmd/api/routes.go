package main

import "net/http"

func (app *application) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)
	router.HandleFunc("GET /v1/user/{username}", app.requireAuthenticatedUser(app.getUserHandler))
	router.HandleFunc("POST /v1/user", app.registerUserHandler)
	router.HandleFunc("GET /v1/user/profile", app.requireAuthenticatedUser(app.getProfileHandler))
	router.HandleFunc("POST /v1/user/authentication", app.authenticationHandler)
	return app.authenticate(router)
}
