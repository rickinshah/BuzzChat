package cache

import (
	"encoding/hex"
	"time"

	"github.com/RickinShah/BuzzChat/internal/model"
	"github.com/redis/go-redis/v9"
)

func GetCachedUser(rdb *redis.Client, username string) (*model.User, bool, error) {
	key := cacheKeyForUser(username)

	return Get[model.User](rdb, key)
}

func SetCachedUser(rdb RedisExecutor, user *model.User) error {
	key := cacheKeyForUser(user.Username)

	return Set(rdb, key, user, time.Hour)
}

func DelCachedUser(rdb RedisExecutor, username string) error {
	key := cacheKeyForUser(username)
	return Del(rdb, key)
}

func SetCachedUserByToken(rdb RedisExecutor, scope string, hash []byte, username string) error {
	key := cacheKeyForUserByToken(scope, hash)
	userKey := cacheKeyForUser(username)
	return Set(rdb, key, userKey, 30*time.Minute)
}

func GetCachedUserByToken(rdb *redis.Client, scope string, hash []byte) (*model.User, bool, error) {
	key := cacheKeyForUserByToken(scope, hash)

	userKey, isCached, err := Get[string](rdb, key)
	if !isCached {
		return nil, false, err
	}

	return Get[model.User](rdb, *userKey)
}

func cacheKeyForUser(username string) string {
	return "user:" + username
}

func cacheKeyForUserByToken(scope string, hash []byte) string {
	return "token:" + scope + ":" + hex.EncodeToString(hash)
}
