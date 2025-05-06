package model

import (
	"encoding/json"
	"strconv"

	"github.com/RickinShah/BuzzChat/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type MarshalType uint8

var AnonymousUser = &User{}

const (
	Minimal MarshalType = iota
	Frontend
	Self
	Full
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
