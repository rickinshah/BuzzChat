package data

import (
	"errors"
	"os"

	"github.com/RickinShah/BuzzChat/internal/db"
	"github.com/RickinShah/BuzzChat/internal/jsonlog"
	"github.com/redis/go-redis/v9"
)

var logger = jsonlog.New(os.Stdout, jsonlog.LevelInfo)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Users  *UserModel
	Tokens *TokenModel
	OTPs   *OTPModel
}

func NewModels(db *db.Queries, redis *redis.Client) Models {
	return Models{
		Users:  &UserModel{db, redis},
		Tokens: &TokenModel{db, redis},
		OTPs:   &OTPModel{db, redis},
	}
}
