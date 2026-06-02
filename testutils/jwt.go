package testutils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const TestJWTSecret = "test-secret"

// GenerateTestToken creates a signed JWT for use in test requests.
func GenerateTestToken(noTelp string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": noTelp,
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	s, _ := token.SignedString([]byte(TestJWTSecret))
	return s
}
