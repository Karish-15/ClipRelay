package utils

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func ExtractUserIDFromJWT(tokenString string) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		// Make sure the signing method is HMAC
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}

	// Try user_id as string or numeric
	if id, ok := claims["sub"].(string); ok {
		return id, nil
	}
	if f, ok := claims["sub"].(float64); ok {
		return fmt.Sprintf("%.0f", f), nil
	}
	return "", errors.New("user_id not found in token")
}
