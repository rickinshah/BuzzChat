package data

import (
	"context"
	"errors"
	"time"

	"github.com/RickinShah/BuzzChat/internal/cache"
	"github.com/RickinShah/BuzzChat/internal/db"
	"github.com/RickinShah/BuzzChat/internal/model"
	"github.com/RickinShah/BuzzChat/internal/security"
	"github.com/RickinShah/BuzzChat/internal/validator"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pquerna/otp/totp"
	"github.com/redis/go-redis/v9"
)

func generateOTP(accountName string, ttl time.Duration, userID int64) (*model.OTP, error) {
	otp := &model.OTP{
		Otp: db.Otp{
			UserPid: userID,
			CreatedAt: pgtype.Timestamptz{
				Time:  time.Now(),
				Valid: true,
			},
			Expiry: pgtype.Timestamptz{
				Time:  time.Now(),
				Valid: true,
			},
		},
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "BuzzChat",
		AccountName: accountName,
	})

	if err != nil {
		return nil, err
	}

	otp.SecretKey = key.Secret()
	otp.Code, err = totp.GenerateCodeCustom(otp.SecretKey, otp.CreatedAt.Time, totp.ValidateOpts{
		Period: 60 * 15,
		Skew:   1,
		Digits: 6,
	})

	if err != nil {
		return nil, err
	}

	return otp, nil
}

func MatchesOTP(otp *model.OTP) (bool, error) {
	return totp.ValidateCustom(otp.Code, otp.SecretKey, time.Now(), totp.ValidateOpts{
		Period: 60 * 15,
		Skew:   1,
		Digits: 6,
	})
}

func ValidateOTPCode(v *validator.Validator, code string) {
	v.Check(len(code) == 6, "otp", "length must be 6")
}

type OTPModel struct {
	DB    *db.Queries
	Redis *redis.Client
}

func (m OTPModel) New(accountName string, ttl time.Duration, userID int64) (*model.OTP, error) {
	otp, err := generateOTP(accountName, ttl, userID)
	if err != nil {
		return nil, err
	}

	secretKey, err := security.EncryptionWithHex(otp.SecretKey)

	if err != nil {
		return nil, err
	}

	otp.SecretKey = secretKey

	return otp, m.Insert(otp)
}

func (m OTPModel) Insert(otp *model.OTP) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := db.InsertOTPParams{
		UserPid:   otp.UserPid,
		CreatedAt: otp.CreatedAt,
		SecretKey: otp.SecretKey,
		Expiry:    otp.Expiry,
	}

	if err := m.DB.InsertOTP(ctx, args); err != nil {
		return err
	}

	if err := cache.SetCachedOTP(m.Redis, otp); err != nil {
		logger.PrintInfo("failed to cache otp:"+err.Error(), nil)
	}

	return nil
}

func (m OTPModel) Delete(userID int64) error {
	if err := cache.DelCachedOTP(m.Redis, userID); err != nil {
		logger.PrintInfo("failed to delete cached otp:"+err.Error(), nil)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.DeleteOTP(ctx, userID)
}

func (m OTPModel) Get(userID int64) (*model.OTP, error) {
	cachedOTP, isCached, _ := cache.GetCachedOTP(m.Redis, userID)
	if isCached {
		secretKey, err := security.DecryptionWithHex(cachedOTP.SecretKey)
		if err != nil {
			return nil, err
		}

		cachedOTP.SecretKey = secretKey
		return cachedOTP, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	otp, err := m.DB.GetOTP(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	secretKey, err := security.DecryptionWithHex(otp.SecretKey)
	if err != nil {
		return nil, err
	}

	otp.SecretKey = secretKey

	customOTP := model.OTP{
		Otp: otp,
	}

	return &customOTP, nil
}
