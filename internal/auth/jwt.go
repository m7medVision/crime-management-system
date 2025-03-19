package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/m7medVision/crime-management-system/internal/model"
)

type JWTClaims struct {
	UserID         uint                 `json:"user_id"`
	Username       string               `json:"username"`
	Role           model.Role           `json:"role"`
	ClearanceLevel model.ClearanceLevel `json:"clearance_level"`
	jwt.StandardClaims
}

func GenerateToken(user *model.User, secret string, expiryHours int) (string, error) {
	claims := JWTClaims{
		UserID:         user.ID,
		Username:       user.Username,
		Role:           user.Role,
		ClearanceLevel: user.ClearanceLevel,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(expiryHours) * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "crime-management-system",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
