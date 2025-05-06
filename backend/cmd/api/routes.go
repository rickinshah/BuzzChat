package main

import "net/http"

func (app *application) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)
	router.HandleFunc("GET /v1/users/{username}", app.requireAuthenticatedUser(app.getUserHandler))
	router.HandleFunc("POST /v1/auth/register", app.registerUserHandler)
	router.HandleFunc("GET /v1/users/me", app.requireAuthenticatedUser(app.getProfileHandler))
	router.HandleFunc("POST /v1/auth/login", app.authenticationHandler)
	router.HandleFunc("GET /v1/users/check-email", app.checkEmailHandler)
	router.HandleFunc("GET /v1/users/check-username", app.checkUsernameHandler)
	return app.enableCORS(app.authenticate(router))
}
