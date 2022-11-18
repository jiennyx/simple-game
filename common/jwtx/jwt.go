package jwtx

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTType string

const (
	JWTTypeAuth    JWTType = "auth"
	JWTTypeRefresh JWTType = "refresh"
)

type JwtInfo struct {
	ID       uint64
	Username string
}

type Claims struct {
	Info JwtInfo
	jwt.StandardClaims
}

const (
	issuer = "simpleGame"
)

var (
	authSecret    = []byte("simpleGameSecret")
	refreshSecret = []byte("simpleGameSecret")
)

func GenerateToken(tp JWTType, info JwtInfo, duration time.Duration) (string, error) {
	expirationTime := time.Now().Add(duration)
	cls := &Claims{
		Info: info,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cls)

	return token.SignedString(getJWTSecret(tp))
}

func ParseToken(tp JWTType, token string) (*Claims, error) {
	var i = len(token)
	for token[i-1] == ' ' {
		i--
	}
	token = token[:i]
	cls := new(Claims)
	_, err := jwt.ParseWithClaims(token, cls,
		func(t *jwt.Token) (interface{}, error) {
			return getJWTSecret(tp), nil
		},
	)

	return cls, err
}

func IsValidIssuer(str string) bool {
	return str == issuer
}

func getJWTSecret(tp JWTType) []byte {
	switch tp {
	case JWTTypeAuth:
		return authSecret
	}

	return refreshSecret
}
