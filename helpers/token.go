package helpers

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTToken struct {
	secretKey string
}

type JWTClaim struct {
	jwt.StandardClaims
	UserID int64 `json:"user_id"`
	Exp    int64 `json:"exp"`
}

func NewJWTToken(secretKey string) *JWTToken {
	return &JWTToken{secretKey: secretKey}
}

func (j *JWTToken) CreateToken(userId uint) (string, error) {
	claims := JWTClaim{
		UserID: int64(userId),
		Exp:    time.Now().Add(time.Minute * 60 * 24 * 7).Unix(), // Set expiration time (30 minutes)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(j.secretKey))

	if err != nil {
		// Consider logging the error (e.g., with a logging library)
		return "", fmt.Errorf("error creating token: %w", err)
	}
	return string(tokenString), nil
}

func (j *JWTToken) VerifyToken(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaim{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid authentication token")
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		// Consider logging the error
		return 0, fmt.Errorf("invalid authentication token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaim)

	if !ok {
		// Consider logging the error
		return 0, fmt.Errorf("invalid authentication token")
	}

	if claims.Exp < time.Now().Unix() {
		return 0, fmt.Errorf("token has expired")
	}

	return claims.UserID, nil
}
