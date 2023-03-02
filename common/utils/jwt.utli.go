package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
)

type JWTClaim struct {
	UserId uint `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userId uint, jwtkey string) string {
	//var c = config.GetConfig()
	var jwtKey = []byte(jwtkey)
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &JWTClaim{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(jwtKey)

	if err != nil {
		log.Fatal("An error occurred while signing token")
	}

	return t
}

func ValidateToken(signedToken string, jwtkey string) (interface{}, error) {
	var jwtKey = []byte(jwtkey)
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		},
	)

	if claims, ok := token.Claims.(*JWTClaim); ok && token.Valid {
		return claims.UserId, nil

	} else {
		if errors.Is(err, jwt.ErrTokenExpired) {
			e := errors.New("token expired")
			return nil, e
		}
		er := errors.New("invalid token")
		return nil, er
	}
}
