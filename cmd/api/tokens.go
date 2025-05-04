package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/RickinShah/BuzzChat/internal/data"
	"github.com/RickinShah/BuzzChat/internal/validator"
)

func (app *application) authenticationHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := app.readJson(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateUsernameOrEmail(v, input.Username); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if data.ValidatePasswordPlainText(v, input.Password); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmailOrUsername(input.Username)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	match, err := data.MatchesPassword(user.PasswordHash, input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	token, err := app.models.Tokens.New(user.UserPid, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	user.SetMarshalType(data.Self)

	cookie := http.Cookie{
		Name:        "auth_token",
		Value:       token.Plaintext,
		Expires:     token.Expiry.Time.UTC(),
		HttpOnly:    true,
		SameSite:    http.SameSiteNoneMode,
		Secure:      true,
		Path:        "/",
		Partitioned: true,
	}

	http.SetCookie(w, &cookie)

	if err := app.writeJson(w, http.StatusOK, envelope{"user": user}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
