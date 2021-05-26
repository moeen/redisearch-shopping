package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// secretKey is the secret used to sign tokens
// @TODO: get secret key from command or env
var secretKey = []byte("super-secret")

// expirationTime is the default expiration time of the generated tokens
const DefaultExpirationTime = 24 * time.Hour

// GenerateToken receives a customer id and generates a token for that customer
func GenerateToken(customerID int, expireAt time.Time) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["customer_id"] = customerID
	claims["exp"] = expireAt

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ParseToken tries to parse a given token a returns the customer id if it was ok
func ParseToken(tokenStr string) (int, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["customer_id"].(float64)
		if !ok {
			return 0, fmt.Errorf("failed to get cliams: %w", err)
		}
		return int(userID), nil
	} else {
		return 0, fmt.Errorf("failed to get cliams: %w", err)
	}
}
