package middleware

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	Client *redis.Client
}

func NewRedisStore(addr, password string, db int) *RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisStore{Client: client}
}

func (r *RedisStore) SaveJwtRefreshToken(ctx context.Context, jti, token string, userID uint, expiration time.Duration) error {
	key := "refresh_token:" + strconv.FormatUint(uint64(userID), 10) + ":" + jti
	return r.Client.Set(ctx, key, token, expiration).Err()
}

func (r *RedisStore) IsJwtRefreshTokenExists(ctx context.Context, jti string, userID uint) (bool, error) {
	key := "refresh_token:" + strconv.FormatUint(uint64(userID), 10) + ":" + jti
	result, err := r.Client.Exists(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

func (r *RedisStore) DeleteJwtRefreshToken(ctx context.Context, jti string, userID uint) error {
	key := "refresh_token:" + strconv.FormatUint(uint64(userID), 10) + ":" + jti
	return r.Client.Del(ctx, key).Err()
}

func (r *RedisStore) BlacklistAccessToken(ctx context.Context, jti string, expiration time.Duration) error {
	key := "bl:access:" + jti
	return r.Client.Set(ctx, key, "blacklisted", expiration).Err()
}
