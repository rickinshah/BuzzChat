package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/RickinShah/BuzzChat/internal/data"
	"github.com/RickinShah/BuzzChat/internal/mailer"
	"github.com/RickinShah/BuzzChat/internal/model"
	"github.com/RickinShah/BuzzChat/internal/validator"
)

func (app *application) generateOtpHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
	}

	if err := app.readJson(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.models.Users.GetByEmailOrUsername(input.Username)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	otp, err := app.models.OTPs.New(user.Username, 15*time.Minute, user.UserPid)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// app.background(func() {
	// 	args := map[string]any{
	// 		"Code": otp.Code,
	// 		"Mail": user.Email,
	// 		"Name": user.Name.String,
	// 	}
	// 	err = app.mailer.Send(user.Email, "otp_verification.tmpl", args)
	// 	if err != nil {
	// 		app.logger.PrintError(err, nil)
	// 	}
	// })
	job := mailer.EmailJob{
		Recipient: user.Email,
		Template:  "otp_verification.tmpl",
		Data: map[string]any{
			"Code": otp.Code,
			"Mail": user.Email,
			"Name": user.Username,
		},
	}
	if err := mailer.EnqueueEmail(app.models.OTPs.Redis, job); err != nil {
		app.serverErrorResponse(w, r, err)
	}

	if err = app.writeJson(w, http.StatusOK, envelope{"username": user.Username}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) validateOtpHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Otp      string `json:"otp"`
	}

	if err := app.readJson(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	user, err := app.models.Users.GetByEmailOrUsername(input.Username)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if data.ValidateOTPCode(v, input.Otp); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	otp, err := app.models.OTPs.Get(user.UserPid)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	otp.Code = input.Otp

	valid, err := data.MatchesOTP(otp)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !valid {
		v.AddError("otp", "is invalid")
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err = app.models.OTPs.Delete(user.UserPid); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	token, err := app.models.Tokens.New(user.UserPid, 1*time.Hour, model.ScopePasswordReset)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err = app.writeJson(w, http.StatusOK, envelope{"reset_token": token.Plaintext}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
