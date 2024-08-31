package utils

import (
	"fmt"
	model "main/DB"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func CompareHashPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func ParseToken(tokenString string) (claims *model.Claims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("ilovegolang"), nil
		// return controllers.JwtKey, nil
	})

	if err != nil {
		fmt.Print("err", err)
		return nil, err
	}

	claims, ok := token.Claims.(*model.Claims)

	if !ok {
		return nil, err
	}

	return claims, nil
}
