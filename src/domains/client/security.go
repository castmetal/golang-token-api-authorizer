package client

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateTokenJWT(tokenDuration time.Duration, clientId string, salt []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorized": true,
		"clientId":   clientId,
		"exp":        time.Now().Add(time.Minute * 2).Unix(),
	})

	tokenString, err := token.SignedString(salt)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateTokenJWT(tokenStr string, clientId string, salt []byte) (bool, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return salt, nil
	})
	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["authorized"] != true || claims["clientId"] != clientId {
			return false, fmt.Errorf("Invalid Token")
		}
	}

	return true, nil
}
