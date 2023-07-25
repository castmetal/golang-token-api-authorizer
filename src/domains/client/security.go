package client

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
)

func GenerateTokenJWT(tokenPeriod string, tokenDuration int32, clientId string, salt []byte) (string, error) {
	var exp int64

	now := time.Now().Local()

	switch tokenPeriod {
	case "days":
		exp = now.AddDate(0, 0, int(tokenDuration)).Unix()
	case "years":
		exp = now.AddDate(int(tokenDuration), 0, 0).Unix()
	case "minutes":
		exp = now.Add(time.Minute * time.Duration(tokenDuration)).Unix()
	case "seconds":
		exp = now.Add(time.Second * time.Duration(tokenDuration)).Unix()
	default:
		exp = now.Add(time.Hour * time.Duration(tokenDuration)).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorized": true,
		"clientId":   clientId,
		"exp":        exp,
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

func GetNewUlid() ulid.ULID {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())

	newUlid, err := ulid.New(ms, entropy)
	if err != nil {
		newUlid = ulid.Make()
	}

	return newUlid
}
