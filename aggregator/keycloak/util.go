package keycloak

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

/* Generate a JWT with a given expiration time. Useful for test cases. */
func GenerateJWT(exp int64) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": exp,
		"iat": time.Now().Unix(),
		"sub": "1234567890",
	})

	secretKey := []byte("secret")
	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return ""
	}

	return tokenString
}
