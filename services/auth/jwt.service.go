package auth_service

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwt(data jwt.MapClaims) (string, error) {
	key := []byte(os.Getenv("JWT_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetDataFromToken(token string) (jwt.MapClaims, error) {
	key := []byte(os.Getenv("JWT_SECRET"))
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}
	return claims, nil
}
