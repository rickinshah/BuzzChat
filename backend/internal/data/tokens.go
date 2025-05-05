package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"

	"github.com/RickinShah/BuzzChat/internal/db"
	"github.com/RickinShah/BuzzChat/internal/validator"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
)

const (
	ScopeAuthentication = "authentication"
)

type Token struct {
	db.Token
	Plaintext string
}

func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token := Token{
		Token: db.Token{
			UserID: userID,
			Expiry: pgtype.Timestamptz{Time: time.Now().Add(ttl), Valid: true},
			Scope:  scope,
		},
	}
	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return &token, nil
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 characters long")
}

type TokenModel struct {
	DB    *db.Queries
	Redis *redis.Client
}

func (m TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := db.InsertTokenParams(token.Token)
	err = m.DB.InsertToken(ctx, args)
	return token, err
}

func (m TokenModel) DeleteAllForUser(scope string, userID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := db.DeleteAllTokensParams{
		Scope:  scope,
		UserID: userID,
	}
	return m.DB.DeleteAllTokens(ctx, args)
}
