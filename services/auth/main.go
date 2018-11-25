package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/thedevsir/frame-backend/services/errors"
)

type (
	UserToken struct {
		Session string `json:"session"`
		SID     string `json:"sid"`
		ID      string `json:"userId"`
		jwt.StandardClaims
	}
	AdminToken struct {
		Session string `json:"session"`
		ID      string `json:"userId"`
		jwt.StandardClaims
	}
)

func (j *UserToken) Create(signingKey []byte) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, j)
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", errors.ErrInternal
	}
	return "Bearer " + tokenString, nil
}

func (j *AdminToken) Create(signingKey []byte) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, j)
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", errors.ErrInternal
	}
	return "Bearer " + tokenString, nil
}
