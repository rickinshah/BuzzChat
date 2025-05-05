package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/RickinShah/BuzzChat/internal/cache"
	"github.com/RickinShah/BuzzChat/internal/db"
	"github.com/RickinShah/BuzzChat/internal/validator"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

var AnonymousUser = &User{}

type MarshalType uint8

const (
	Minimal MarshalType = iota
	Frontend
	Self
	Full
)

var (
	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrDuplicateUsername = errors.New("duplicate username")
)

type User struct {
	db.User
	marshalType MarshalType
}

func NewUser(user *db.User) *User {
	if user == nil {
		return &User{
			marshalType: Full,
		}
	}
	return &User{
		User:        *user,
		marshalType: Full,
	}
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func (u *User) MarshalJSON() ([]byte, error) {
	user := make(map[string]any, 11)
	if u.marshalType >= Minimal {
		user["userId"] = strconv.FormatInt(u.UserPid, 10)
		user["username"] = u.Username
		user["profilePic"] = u.ProfilePic
		user["name"] = u.Name
	}

	if u.marshalType >= Frontend {
		user["bio"] = u.Bio
	}

	if u.marshalType >= Self {
		user["activated"] = u.Activated
		user["createdAt"] = u.CreatedAt
		user["updatedAt"] = u.UpdatedAt
		user["email"] = u.Email
	}

	if u.marshalType >= Full {
		user["passwordHash"] = u.PasswordHash
		user["version"] = u.Version
	}

	return json.Marshal(user)
}

func (u *User) UnmarshalJSON(data []byte) error {
	var temp struct {
		UserID       string             `json:"userId"`
		Username     string             `json:"username"`
		ProfilePic   pgtype.Text        `json:"profilePic"`
		Name         pgtype.Text        `json:"name"`
		Bio          pgtype.Text        `json:"bio"`
		Activated    bool               `json:"activated"`
		CreatedAt    pgtype.Timestamptz `json:"createdAt"`
		UpdatedAt    pgtype.Timestamptz `json:"updatedAt"`
		Email        string             `json:"email"`
		PasswordHash []byte             `json:"passwordHash"`
		Version      int32              `json:"version"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	userID, err := strconv.ParseInt(temp.UserID, 10, 64)
	if err != nil {
		return err
	}

	u.UserPid = userID
	u.Username = temp.Username
	u.ProfilePic = temp.ProfilePic
	u.Name = temp.Name
	u.Bio = temp.Bio
	u.Activated = temp.Activated
	u.CreatedAt = temp.CreatedAt
	u.UpdatedAt = temp.UpdatedAt
	u.Email = temp.Email
	u.PasswordHash = temp.PasswordHash
	u.Version = temp.Version

	return nil
}

func (u *User) SetMarshalType(mt MarshalType) {
	u.marshalType = mt
}

func SetPassword(plainTextPassword string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func MatchesPassword(hashedPassword []byte, plainTextPassword string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(plainTextPassword)); err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidateUsername(v *validator.Validator, username string) {
	v.Check(username != "", "username", "must be provided")
	v.Check(len(username) <= 30, "username", "must not be more than 30 characters")
	v.Check(validator.Matches(username, validator.UsernameRX), "username", "should only contain alphanumeric characters and underscore")
}

func ValidateUsernameOrEmail(v *validator.Validator, username string) {
	v.Check(username != "", "username/email", "must be provided")
	v.Check(validator.Matches(username, validator.UsernameEmailRX), "username/email", "is not valid")
}

func ValidateBio(v *validator.Validator, bio *pgtype.Text) {
	if bio.Valid {
		v.Check(len(bio.String) <= 300, "bio", "must not be more than 500 characters")
	}
	validator.NullifyEmptyText(bio)
}

func ValidateName(v *validator.Validator, name *pgtype.Text) {
	if name.Valid {
		v.Check(len(name.String) <= 50, "name", "must not be more than 50 characters")
	}
	validator.NullifyEmptyText(name)
}

func ValidatePasswordPlainText(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must not be less than 8 characters")
	v.Check(len(password) <= 72, "password", "must not be more than 72 characters")
}

func ValidateConfirmPassword(v *validator.Validator, password string, confirmPassword string) {
	v.Check(confirmPassword != "", "confirm password", "must be provided")
	v.Check(password == confirmPassword, "password", "doesn't match")
}

func ValidateUser(v *validator.Validator, user *User) {
	ValidateUsername(v, user.Username)
	ValidateEmail(v, user.Email)
	ValidateBio(v, &user.Bio)
	ValidateName(v, &user.Name)
}

type UserModel struct {
	DB    *db.Queries
	Redis *redis.Client
}

func (m UserModel) Insert(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := db.InsertUserParams{
		Username:     user.Username,
		Email:        user.Email,
		Name:         user.Name,
		PasswordHash: user.PasswordHash,
		ProfilePic:   user.ProfilePic,
	}
	row, err := m.DB.InsertUser(ctx, args)
	if err != nil {
		var pgErr *pgconn.PgError
		switch {
		case errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "users_email_key":
			return ErrDuplicateEmail
		case errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "users_username_key":
			return ErrDuplicateUsername
		default:
			return err
		}
	}
	user.UserPid = row.UserPid
	user.CreatedAt = row.CreatedAt
	user.UpdatedAt = row.UpdatedAt
	user.Version = row.Version
	user.marshalType = Self

	if err := m.SetCachedUser(user); err != nil {
		logger.PrintInfo("failed to cache user:"+err.Error(), nil)
	}

	return nil
}

func (m UserModel) UpdatePassword(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := db.UpdatePasswordParams{
		PasswordHash: user.PasswordHash,
		UserPid:      user.UserPid,
		Version:      user.Version,
	}

	row, err := m.DB.UpdatePassword(ctx, args)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	user.UpdatedAt = row.UpdatedAt
	user.Version = row.Version

	if err := m.SetCachedUser(user); err != nil {
		logger.PrintInfo("failed to cache user:"+err.Error(), nil)
	}

	return nil
}

func (m UserModel) Update(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := db.UpdateUserParams{
		UserPid:    user.UserPid,
		Username:   user.Username,
		Email:      user.Email,
		Name:       user.Name,
		Activated:  user.Activated,
		Bio:        user.Bio,
		ProfilePic: user.ProfilePic,
		Version:    user.Version,
	}
	row, err := m.DB.UpdateUser(ctx, args)
	if err != nil {
		var pgErr *pgconn.PgError
		switch {
		case errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "users_email_key":
			return ErrDuplicateEmail
		case errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "users_username_key":
			return ErrDuplicateUsername
		default:
			return err
		}
	}
	user.UpdatedAt = row.UpdatedAt
	user.Version = row.Version

	if err := m.SetCachedUser(user); err != nil {
		logger.PrintInfo("failed to cache user:"+err.Error(), nil)
	}

	return nil
}

func (m UserModel) Delete(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.DeleteUser(ctx, user.UserPid); err != nil {
		return err
	}

	if err := m.DelCachedUser(user.Username); err != nil {
		logger.PrintInfo("failed to delete cached user:"+err.Error(), nil)
	}
	return nil
}

func (m UserModel) GetByUsername(username string) (*User, error) {
	cachedUser, isCached, _ := m.GetCachedUser(username)
	if isCached {
		return cachedUser, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user, err := m.DB.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	customUser := NewUser(&user)

	if err := m.SetCachedUser(customUser); err != nil {
		logger.PrintInfo("failed to cache user:"+err.Error(), nil)
	}

	return customUser, nil
}

func (m UserModel) GetByEmailOrUsername(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cachedUser, isCached, _ := m.GetCachedUser(email)
	if isCached {
		return cachedUser, nil
	}

	user, err := m.DB.GetUserByEmailOrUsername(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	customUser := NewUser(&user)

	if err = m.SetCachedUser(customUser); err != nil {
		logger.PrintInfo("failed to cache user:"+err.Error(), nil)
	}

	return customUser, nil
}

func (m UserModel) GetByToken(tokenScope, tokenPlaintext string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	cachedUser, isCached, _ := m.GetCachedUserByToken(tokenScope, tokenHash[:])
	if isCached {
		return cachedUser, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := db.GetUserByTokenParams{
		Hash:  tokenHash[:],
		Scope: tokenScope,
		Expiry: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	}

	user, err := m.DB.GetUserByToken(ctx, args)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	customUser := NewUser(&user)

	if err = m.SetCachedUserByToken(tokenScope, tokenHash[:], customUser); err != nil {
		logger.PrintInfo("failed to cache token:"+err.Error(), nil)
	}

	return customUser, nil
}

func (m UserModel) GetCachedUser(username string) (*User, bool, error) {
	key := cacheKeyForUser(username)

	return cache.Get[User](m.Redis, key)
}

func (m UserModel) SetCachedUser(user *User) error {
	key := cacheKeyForUser(user.Username)

	return cache.Set(m.Redis, key, user, time.Hour)
}

func (m UserModel) DelCachedUser(username string) error {
	key := cacheKeyForUser(username)
	return cache.Del(m.Redis, key)
}

func (m UserModel) SetCachedUserByToken(scope string, hash []byte, user *User) error {
	key := cacheKeyForUserByToken(scope, hash)
	return cache.Set(m.Redis, key, user, time.Hour)
}

func (m UserModel) GetCachedUserByToken(scope string, hash []byte) (*User, bool, error) {
	key := cacheKeyForUserByToken(scope, hash)
	return cache.Get[User](m.Redis, key)
}

func cacheKeyForUser(username string) string {
	return "user:" + username
}

func cacheKeyForUserByToken(scope string, hash []byte) string {
	return "token:" + scope + ":" + hex.EncodeToString(hash)
}
