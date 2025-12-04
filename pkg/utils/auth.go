package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var (
	AcessTTL   = time.Minute * 15
	RefreshTTL = time.Hour * 24 * 7
	JWTSecret  = []byte("replace-with-env-secret")
)

func GenerateTokens(userID uint, role uint) (string, string, string, error) {
	now := time.Now()
	accessJTI := uuid.NewString()
	accessClaims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  now.Add(AcessTTL).Unix(),
		"iat":  now.Unix(),
		"jti":  accessJTI,
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	acessStr, err := at.SignedString(JWTSecret)
	if err != nil {
		return "", "", "", err
	}
	refreshJTI := uuid.NewString()
	refreshClaims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  now.Add(RefreshTTL).Unix(),
		"iat":  now.Unix(),
		"jti":  refreshJTI,
		"type": "refresh",
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshStr, err := rt.SignedString(JWTSecret)
	if err != nil {
		return "", "", "", err
	}
	return acessStr, refreshStr, accessJTI, nil
}

func ParseToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return JWTSecret, nil
	})
	if err != nil || token == nil || !token.Valid {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrInvalidKeyType
	}
	return claims, nil
}

func HashPassword(password string) string {
	return "hashed_" + password
}
