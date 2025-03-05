package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/simance-ai/smdx/rpcx/ws/internal/types"
)

func ParseAccountAuthClaims(accessToken string, secret string) (*types.AccountAuthClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &types.AccountAuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("token parse error")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*types.AccountAuthClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("token claims error")
	}
}

func SignAccountAuthClaims(userID int64, secret string, expiresIn int) (string, error) {
	code := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"email":   "test@183.com",
		"env":     "test",
		"nbf":     time.Now().Unix(),
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Duration(expiresIn) * time.Second).Unix(),
	})

	token, err := code.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return token, nil
}
