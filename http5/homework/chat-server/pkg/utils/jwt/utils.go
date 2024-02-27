package jwt

import "github.com/golang-jwt/jwt/v5"

func CreateJWT(payload jwt.MapClaims, method jwt.SigningMethod, secret string) (string, error) {
	return jwt.NewWithClaims(method, payload).SignedString(secret)
}
