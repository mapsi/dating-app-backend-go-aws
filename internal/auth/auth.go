package auth

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
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

func GetUserIDFromToken(c *fiber.Ctx) (string, error) {
	if len(jwtKey) == 0 {
		return "", errors.New("JWT key not initialized")
	}

	user := c.Locals("user").(*jwt.Token)
	if user == nil {
		return "", errors.New("no token found in context")
	}

	claims := user.Claims.(jwt.MapClaims)
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("user_id not found in token claims")
	}

	return userID, nil
}
