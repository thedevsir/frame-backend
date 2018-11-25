package request

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsir/frame-backend/services/errors"
	"github.com/thedevsir/frame-backend/services/response"
)

type (
	User struct {
		Session string
		SID     string
		ID      string
	}
	Admin struct {
		Session string
		ID      string
	}
)

func GetInputs(c echo.Context, params interface{}) (err error) {

	if err = c.Bind(params); err != nil {
		return errors.ErrInvalidParams
	}

	if err = c.Validate(params); err != nil {
		return response.ValidationError(err)
	}

	return nil
}

func AuthenticatedUser(c echo.Context) User {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return User{
		Session: claims["session"].(string),
		SID:     claims["sid"].(string),
		ID:      claims["userId"].(string),
	}
}

func AuthenticatedAdmin(c echo.Context) Admin {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return Admin{
		Session: claims["session"].(string),
		ID:      claims["userId"].(string),
	}
}

func IsRootAdmin(c echo.Context) bool {

	admin := AuthenticatedAdmin(c)
	return admin.ID == "000000000000000000000000"
}
