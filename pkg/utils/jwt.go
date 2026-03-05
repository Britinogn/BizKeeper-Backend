package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken      = errors.New("invalid token")
	ErrTokenExpired      = errors.New("token expired")
	ErrMissingClaims     = errors.New("missing required claims")
	ErrUnexpectedSigning = errors.New("unexpected signing method")
)

type Cliams struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role,omitempty"`
	jwt.RegisteredClaims
}

func getSecretKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is not set")
	}
	return []byte(secret)
}

func GenerateToken(userID, email, role string) (string, error) {
	// read expiration from env
	expirationTime := os.Getenv("JWT_EXPIRES_IN")
	if expirationTime == "" {
		expirationTime = "24h"
	}

	duration, err := time.ParseDuration(expirationTime)
	if err != nil {
		duration = 24 * time.Hour
	}

	claims := Cliams{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "bizkeeper",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(getSecretKey())
	if err != nil {
		return "", err
	}
	return signed, nil
}

func ValidateToken(tokenString string) (*Cliams, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Cliams{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigning
		}
		return getSecretKey(), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Cliams)
	if !ok {
		return nil, ErrMissingClaims
	}

	return claims, nil
}
