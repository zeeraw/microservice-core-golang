package auth

import (
	"crypto/rsa"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	errorMessageMalformed     = "token malformed"
	errorMessageExpired       = "token expired or not yet valid"
	errorMessageInvalid       = "invalid token"
	errorMessageClaimsInvalid = "invalid token claims"
)

// JWTClaims represents the claims within the JWT.
type JWTClaims struct {
	Consumer Consumer `json:"consumer"`
	jwt.StandardClaims
}

// ParseJWT parses a JWT string and checks its signature validity
func ParseJWT(pk *rsa.PublicKey, raw string) (*jwt.Token, error) {
	// Parse the JWT token
	token, err := jwt.ParseWithClaims(raw, &JWTClaims{}, checkSignatureFunc(pk))

	// Bail out if the token could not be parsed
	if err != nil {
		if _, ok := err.(*jwt.ValidationError); ok {
			// Handle any token specific errors
			var errorMessage string
			if err.(*jwt.ValidationError).Errors&jwt.ValidationErrorMalformed != 0 {
				errorMessage = errorMessageMalformed
			} else if err.(*jwt.ValidationError).Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				errorMessage = errorMessageExpired
			} else {
				errorMessage = errorMessageInvalid
			}
			return nil, fmt.Errorf(errorMessage)
		}
		return nil, fmt.Errorf(errorMessageInvalid)
	}

	// Check the claims and token are valid
	if _, ok := token.Claims.(*JWTClaims); !ok || !token.Valid {
		return nil, fmt.Errorf(errorMessageClaimsInvalid)
	}

	return token, nil
}

func checkSignatureFunc(pk *rsa.PublicKey) func(t *jwt.Token) (interface{}, error) {
	return func(t *jwt.Token) (interface{}, error) {
		// Ensure the signing method was not changed
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return pk, nil
	}
}