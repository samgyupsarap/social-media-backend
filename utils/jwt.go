package utils

import (
	"backend/db/socmed"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type Claims struct {
	UserUUID string `json:"user_uuid"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	jwt.RegisteredClaims
}

func GenerateToken(user socmed.User) (string, error) {
	claims := &Claims{
		UserUUID: user.UserUuid,
		UserName: user.UserName.String,
		Email:    user.Email.String,
		FullName: user.FullName.String,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
