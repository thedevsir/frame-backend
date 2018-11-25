package jwt

import (
	"fmt"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/thedevsir/frame-backend/services/errors"
)

type EmailData struct {
	Action   string
	UserID   string
	Email    string
	Username string
}

func ParseJWT(tokenString string, secret []byte) (*jwt.Token, error) {

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return token, nil
	}

	return nil, err
}

func ParseEmailToken(token string, secret []byte) (*EmailData, error) {

	data, err := ParseJWT(token, secret)
	if err != nil {
		return nil, errors.ErrTokenIsNotValid
	}
	claims := data.Claims.(jwt.MapClaims)
	return &EmailData{
		Action:   claims["action"].(string),
		UserID:   claims["userId"].(string),
		Email:    claims["email"].(string),
		Username: claims["username"].(string),
	}, nil
}
