package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	cases := []struct {
		customerID int
		expire     time.Time
	}{
		{customerID: 1, expire: time.Now().Add(time.Hour)},
		{customerID: 100, expire: time.Now().Add(2 * time.Hour)},
		{customerID: 1000, expire: time.Now().Add(3 * time.Hour)},
	}

	for _, tc := range cases {
		token, err := GenerateToken(tc.customerID, tc.expire)
		assert.NoError(t, err)

		pt, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		assert.NoError(t, err)

		claims, ok := pt.Claims.(jwt.MapClaims)
		assert.True(t, ok)
		assert.True(t, pt.Valid)

		ex, err := time.Parse(time.RFC3339, claims["exp"].(string))
		assert.NoError(t, err)
		assert.Equal(t, time.Duration(0), tc.expire.Sub(ex))

		assert.Equal(t, float64(tc.customerID), claims["customer_id"].(float64))
	}
}

func TestParseToken(t *testing.T) {
	t.Run("test valid tokens", func(t *testing.T) {
		cases := []int{10, 100, 100}

		for _, tc := range cases {
			token := jwt.New(jwt.SigningMethodHS256)

			claims := token.Claims.(jwt.MapClaims)
			claims["customer_id"] = tc
			tokenStr, _ := token.SignedString(secretKey)

			cid, err := ParseToken(tokenStr)
			assert.NoError(t, err)
			assert.Equal(t, tc, cid)
		}
	})

	t.Run("test invalid token strings", func(t *testing.T) {
		cid, err := ParseToken("invalid-jwt-token")
		assert.Error(t, err)
		assert.Equal(t, 0, cid)
	})

	t.Run("test token with no claims", func(t *testing.T) {
		token := jwt.New(jwt.SigningMethodHS256)
		tokenStr, _ := token.SignedString(secretKey)
		cid, err := ParseToken(tokenStr)
		assert.Error(t, err)
		assert.Equal(t, 0, cid)
	})

	t.Run("test token with invalid claims", func(t *testing.T) {
		token := jwt.New(jwt.SigningMethodHS256)
		tokenStr, _ := token.SignedString(secretKey)
		cid, err := ParseToken(tokenStr)
		assert.Error(t, err)
		assert.Equal(t, 0, cid)
	})
}
