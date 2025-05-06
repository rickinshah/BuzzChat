package model

import (
	"github.com/RickinShah/BuzzChat/internal/db"
)

const (
	ScopeAuthentication = "authentication"
	ScopeRegistration   = "registration"
)

type Token struct {
	db.Token
	Plaintext string
}
