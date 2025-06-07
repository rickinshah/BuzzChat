package model

import (
	"github.com/RickinShah/BuzzChat/internal/db"
)

const (
	ScopeAuthentication = "authentication"
	ScopeRegistration   = "registration"
	ScopePasswordReset  = "password_reset"
)

type Token struct {
	db.Token
	Plaintext string
}
