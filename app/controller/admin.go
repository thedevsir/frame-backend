package controller

import (
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsir/frame-backend/app/repository"
	"github.com/thedevsir/frame-backend/config"
	"github.com/thedevsir/frame-backend/services/auth"
	"github.com/thedevsir/frame-backend/services/errors"
	"github.com/thedevsir/frame-backend/services/request"
	r "github.com/thedevsir/frame-backend/services/response"
)

type (
	AdminSigninSchema struct {
		Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
		Password string `json:"password" validate:"required,min=8,max=50"`
	}
)

// AdminSignin godoc
// @Summary Admin signin
// @Tags admin
// @Accept json
// @Produce json
// @Param username body string true "Username"
// @Param password body string true "Password"
// @Success 200 {object} response.Message
// @Router /admin/signin [post]
func AdminSignin(c echo.Context) (err error) {

	params := new(AdminSigninSchema)
	if err = request.GetInputs(c, params); err != nil {
		return err
	}

	admin, err := repository.FindAdminByCredentials(params.Username, params.Password)
	if err != nil {
		return err
	}

	adminID := admin.Id.Hex()
	uuid, err := repository.SetAdminSession(adminID)
	if err != nil {
		return err
	}

	token := &auth.AdminToken{
		uuid,
		adminID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	jwtToken, err := token.Create([]byte(config.AdminSigningKey))
	if err != nil {
		return err
	}

	return r.CustomErrorWithHeader(http.StatusOK, "Success", map[string]string{echo.HeaderAuthorization: jwtToken}, c)
}

// AdminSignout godoc
// @Summary Admin signout
// @Tags admin
// @Accept json
// @Produce json
// @Security AdminApiKeyAuth
// @Success 200 {object} response.Message
// @Router /admin/signout [delete]
func AdminSignout(c echo.Context) (err error) {

	admin := request.AuthenticatedAdmin(c)

	if err = repository.TerminateAdminSession(admin.ID); err != nil {
		return err
	}

	return errors.ErrSuccess
}
