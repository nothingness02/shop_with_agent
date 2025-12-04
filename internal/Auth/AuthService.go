package auth

import (
	"context"
	"errors"

	user "github.com/myproject/shop/internal/User"
	"github.com/myproject/shop/pkg/middleware"
	"github.com/myproject/shop/pkg/utils"
)

type AuthService struct {
	userRepo *user.UserRepository
	redis    *middleware.RedisStore
}

func NewAuthService(uRepo *user.UserRepository, redisStore *middleware.RedisStore) *AuthService {
	return &AuthService{userRepo: uRepo, redis: redisStore}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (access, refresh string, userID uint, role uint, err error) {
	u, err := s.userRepo.GetUserByName(username)
	if err != nil {
		return "", "", 0, 0, errors.New("invalid credentials")
	}
	// existing project uses utils.HashPassword as placeholder
	if utils.HashPassword(password) != u.Password {
		return "", "", 0, 0, errors.New("invalid credentials")
	}
	access, refresh, _, err = utils.GenerateTokens(u.ID, u.Role)
	if err != nil {
		return "", "", 0, 0, err
	}
	// parse refresh token to get jti
	claims, err := utils.ParseToken(refresh)
	if err != nil {
		return "", "", 0, 0, err
	}
	jti, _ := claims["jti"].(string)
	// save refresh into redis
	if s.redis != nil {
		_ = s.redis.SaveJwtRefreshToken(ctx, jti, refresh, u.ID, utils.RefreshTTL)
	}
	return access, refresh, u.ID, u.Role, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (newAccess, newRefresh string, userID uint, role uint, err error) {
	claims, err := utils.ParseToken(refreshToken)
	if err != nil {
		return "", "", 0, 0, errors.New("invalid refresh token")
	}
	// 确保是 refresh token
	if typ, ok := claims["type"].(string); ok && typ != "refresh" {
		return "", "", 0, 0, errors.New("not a refresh token")
	}
	// get user id and jti
	subFloat := claims["sub"]
	var uid uint

	switch v := subFloat.(type) {
	case float64:
		uid = uint(v)
	case int:
		uid = uint(v)
	}

	jti, _ := claims["jti"].(string)
	exists, err := s.redis.IsJwtRefreshTokenExists(ctx, jti, uid)
	if err != nil || !exists {
		return "", "", 0, 0, errors.New("refresh token invalid or revoked")
	}

	// delete old refresh
	_ = s.redis.DeleteJwtRefreshToken(ctx, jti, uid)
	// issue new tokens
	newAccess, newRefresh, _, err = utils.GenerateTokens(uid, uint(claims["role"].(float64)))
	if err != nil {
		return "", "", 0, 0, err
	}
	// parse new refresh jti and save
	newClaims, err := utils.ParseToken(newRefresh)
	if err == nil {
		newJti, _ := newClaims["jti"].(string)
		_ = s.redis.SaveJwtRefreshToken(ctx, newJti, newRefresh, uid, utils.RefreshTTL)
	}
	return newAccess, newRefresh, uid, uint(claims["role"].(float64)), nil
}

func (s *AuthService) Logout(ctx context.Context, userID uint, accessJti string) error {
	// delete all refresh tokens for user is safer; here we assume caller knows jti if needed
	// we will blacklist the access token jti
	if s.redis != nil {
		_ = s.redis.BlacklistAccessToken(ctx, accessJti, utils.AcessTTL)
	}
	return nil
}
