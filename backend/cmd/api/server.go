package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         ":" + strconv.Itoa(app.config.port),
		Handler:      app.routes(),
		ErrorLog:     log.New(app.logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.logger.PrintInfo("starting server", map[string]any{
		"addr": srv.Addr,
		"port": "4000",
	})

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
