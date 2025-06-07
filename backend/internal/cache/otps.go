package cache

import (
	"strconv"
	"time"

	"github.com/RickinShah/BuzzChat/internal/model"
	"github.com/redis/go-redis/v9"
)

func SetCachedOTP(rdb RedisExecutor, otp *model.OTP) error {
	key := cacheKeyForOTP(otp.UserPid)

	return Set(rdb, key, otp, time.Hour)
}

func GetCachedOTP(rdb *redis.Client, userID int64) (*model.OTP, bool, error) {
	key := cacheKeyForOTP(userID)

	return Get[model.OTP](rdb, key)
}

func DelCachedOTP(rdb *redis.Client, userID int64) error {
	key := cacheKeyForOTP(userID)

	return Del(rdb, key)
}

func cacheKeyForOTP(userID int64) string {
	return "otp:" + strconv.FormatInt(userID, 10)
}
