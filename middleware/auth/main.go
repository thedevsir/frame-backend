package auth

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsir/frame-backend/app/repository"
)

func Middleware(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		SID, session := claims["sid"].(string), claims["session"].(string)

		if err := repository.SessionFindByCredentials(session, SID); err != nil {
			return err
		}

		if err := repository.SessionUpdateLastActivity(SID); err != nil {
			return err
		}

		return next(c)
	}
}

func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		adminID, session := claims["userId"].(string), claims["session"].(string)

		if err := repository.CheckAdminSession(adminID, session); err != nil {
			return err
		}

		if err := repository.AdminSessionUpdateLastActivity(adminID); err != nil {
			return err
		}

		return next(c)
	}
}
