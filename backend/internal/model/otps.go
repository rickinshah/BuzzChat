package model

import "github.com/RickinShah/BuzzChat/internal/db"

type OTP struct {
	db.Otp
	Code string `json:"-"`
}
