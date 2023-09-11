package util

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func Hash(s string) string {
	hash := sha256.Sum256([]byte(s))

	return hex.EncodeToString(hash[:])
}

func GenJWT(expiration time.Time, payload map[string]any) string {
	secret := []byte(Get("JWT_SECRET"))

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = expiration.Unix()
	for key, val := range payload {
		claims[key] = val
	}

	tokenString, err := token.SignedString(secret)

	if err != nil {
		fmt.Printf("Error while signing JWT : %q", err)
		return ""
	}

	return tokenString
}

var ErrInvalidToken = errors.New("jwt token is not valid")

func VerifyJWT(raw_token string) (map[string]any, error) {
	token, err := jwt.Parse(raw_token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return map[string]any{
				"error": "Unauthorized",
			}, errors.New("Unauthorized")
		}
		return []byte(Get("JWT_SECRET")), nil
	})

	if err != nil {
		return map[string]any{}, err
	}

	var toReturn map[string]any = make(map[string]any)

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		for key, val := range claims {
			toReturn[key] = val
		}
		return toReturn, nil
	}
	return map[string]any{}, ErrInvalidToken
}
