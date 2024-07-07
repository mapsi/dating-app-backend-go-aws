package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey []byte

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func InitJWTSecret(secret string) {
	jwtKey = []byte(secret)
}

func GenerateToken(userID string) (string, error) {
	if len(jwtKey) == 0 {
		return "", errors.New("JWT secret not initialized")
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
