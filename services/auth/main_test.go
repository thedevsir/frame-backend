package auth

import (
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestJWTUserToken(t *testing.T) {

	t.Run("UserToken", func(t *testing.T) {
		composer := &UserToken{
			"session",
			"SID",
			"ID",
			jwt.StandardClaims{},
		}
		assert.NotPanics(t, func() { composer.Create([]byte("secret")) })
	})

	t.Run("AdminToken", func(t *testing.T) {
		composer := &AdminToken{
			"session",
			"ID",
			jwt.StandardClaims{},
		}
		assert.NotPanics(t, func() { composer.Create([]byte("secret")) })
	})
}
