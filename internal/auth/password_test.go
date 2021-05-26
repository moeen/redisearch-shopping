package auth

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestCheckPasswordHash(t *testing.T) {
	t.Run("test with correct password", func(t *testing.T) {
		pass := "test"
		h, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
		assert.NoError(t, err)
		ok := CheckPasswordHash(pass, string(h))
		assert.True(t, ok)
	})

	t.Run("test with invalid hash", func(t *testing.T) {
		ok := CheckPasswordHash("test", "invalid-hash")
		assert.False(t, ok)
	})
}

func TestHashPassword(t *testing.T) {
	pass := "test"
	h, err := HashPassword(pass)
	assert.NoError(t, err)

	err = bcrypt.CompareHashAndPassword([]byte(h), []byte(pass))
	assert.NoError(t, err)
}
