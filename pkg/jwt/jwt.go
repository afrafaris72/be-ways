package jwtToken

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

var SecretKey = os.Getenv("SECRET_KEY")

func GenerateToken(claims *jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	webtoken, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return webtoken, nil
}

func VerifyToken(tokenstring string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("unexpected sign in method : %v", token.Header["alg"])
		}

		return []byte(SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func DecodeToken(tokenstring string) (jwt.MapClaims, error) {
	token, err := VerifyToken(tokenstring)
	if err != nil {
		return nil, err
	}

	claims, tokenOK := token.Claims.(jwt.MapClaims)
	if tokenOK && token.Valid {
		return claims, nil

	}
	return nil, fmt.Errorf("Invalid Token")
}
