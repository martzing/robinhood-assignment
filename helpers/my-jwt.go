package helpers

import (
	"robinhood-assignment/config"
	"robinhood-assignment/internal/core/domains"
	"robinhood-assignment/internal/core/ports"

	"github.com/golang-jwt/jwt/v5"
)

type myJWT struct {
}

func NewMyJWT() ports.MyJWT {
	return &myJWT{}
}

func (j myJWT) NewWithClaims(method jwt.SigningMethod, claims domains.Claims, opts ...jwt.TokenOption) *jwt.Token {
	token := jwt.NewWithClaims(method, claims, opts...)
	return token
}

func (j myJWT) ParseWithClaims(tokenString string, claims *domains.Claims, keyFunc jwt.Keyfunc, opts ...jwt.ParserOption) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc, opts...)
	return token, err
}

func (j myJWT) ParseToken(token *jwt.Token) (interface{}, error) {
	return []byte(config.Get().Auth.JwtSecret), nil
}
