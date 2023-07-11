package ports

import (
	"robinhood-assignment/internal/core/domains"

	"github.com/golang-jwt/jwt/v5"
)

type MyBcrypt interface {
	GenerateFromPassword(password string, cost int) (*string, error)
	CompareHashAndPassword(hashedPassword, password string) error
}

type MyJWT interface {
	NewWithClaims(method jwt.SigningMethod, claims domains.Claims, opts ...jwt.TokenOption) *jwt.Token
	ParseWithClaims(tokenString string, claims *domains.Claims, keyFunc jwt.Keyfunc, opts ...jwt.ParserOption) (*jwt.Token, error)
	ParseToken(token *jwt.Token) (interface{}, error)
}
