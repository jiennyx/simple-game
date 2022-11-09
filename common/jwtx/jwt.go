package jwtx

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JwtInfo struct {
	ID uint64
}

type Claims struct {
	Info JwtInfo
	jwt.StandardClaims
}

var (
	authSecret    = []byte("simpleGameSecret")
	refreshSecret = []byte("simpleGameSecret")
)

func GenerateToken(info JwtInfo, duration time.Duration) (string, error) {
	expirationTime := time.Now().Add(duration)
	cls := &Claims{
		Info: info,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "simpleGame",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cls)

	return token.SignedString(token)
}

func ParseToken(token string) (*Claims, error) {
	cls := new(Claims)
	_, err := jwt.ParseWithClaims(token, cls,
		func(t *jwt.Token) (interface{}, error) {
			return authSecret, nil
		},
	)

	return cls, err
}
