package client

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateTokenJWT(tokenDuration time.Duration, clientId string, salt []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
		"exp":        time.Now().Add(tokenDuration),
		"authorized": true,
		"clientId":   clientId,
	})

	tokenString, err := token.SignedString(salt)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateTokenJWT(tokenStr string, salt []byte) (bool, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return salt, nil
	})
	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["authorize"] != false {
			return false, fmt.Errorf("Invalid Token")
		}
	}

	return true, nil
}
