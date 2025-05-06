package main

import (
	"errors"
	"net/http"

	"github.com/RickinShah/BuzzChat/internal/data"
	"github.com/RickinShah/BuzzChat/internal/model"
	"github.com/RickinShah/BuzzChat/internal/validator"
	"github.com/jackc/pgx/v5/pgtype"
)

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
	}

	input.Username = app.readStringPath(r, "username", "")

	v := validator.New()
	if data.ValidateUsername(v, input.Username); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByUsername(input.Username)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	user.SetMarshalType(model.Frontend)

	if err := app.writeJson(w, http.StatusOK, envelope{"user": user}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username        string
		Email           string
		Name            pgtype.Text
		Password        string
		ConfirmPassword string
	}

	if err := app.readJson(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := model.NewUser(nil)
	user.Username = input.Username
	user.Email = input.Email
	user.Name = input.Name

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if data.ValidatePasswordPlainText(v, input.Password); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if data.ValidateConfirmPassword(v, input.Password, input.ConfirmPassword); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	password, err := data.SetPassword(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	user.PasswordHash = password

	if err = app.models.Users.Insert(user); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrDuplicateUsername):
			v.AddError("username", "a user with this username already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err = app.writeJson(w, http.StatusOK, envelope{"message": "Account created succesfully!"}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getProfileHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	user.SetMarshalType(model.Self)

	if err := app.writeJson(w, http.StatusOK, envelope{"user": user}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) checkEmailHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string
	}

	qs := r.URL.Query()
	input.Email = app.readString(qs, "email", "")

	v := validator.New()
	if data.ValidateEmail(v, input.Email); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	_, err := app.models.Users.GetByEmail(input.Email)
	if err == nil {
		app.conflictResponse(w, r, "a user with this email already exists.")
		return
	}

	if err := app.writeJson(w, http.StatusOK, nil, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) checkUsernameHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string
	}

	qs := r.URL.Query()
	input.Username = app.readString(qs, "username", "")

	v := validator.New()
	if data.ValidateUsername(v, input.Username); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	_, err := app.models.Users.GetByUsername(input.Username)
	if err == nil {
		app.conflictResponse(w, r, "a user with this username already exists.")
		return
	}

	if err := app.writeJson(w, http.StatusOK, nil, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
