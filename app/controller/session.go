package controller

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/thedevsir/frame-backend/app/repository"
	"github.com/thedevsir/frame-backend/services/errors"
	"github.com/thedevsir/frame-backend/services/paginate"
	"github.com/thedevsir/frame-backend/services/request"
	r "github.com/thedevsir/frame-backend/services/response"
)

type (
	SignoutSchema struct {
		ID string `json:"id" validate:"required"`
	}
)

// Sessions godoc
// @Summary Get user sessions
// @Tags session
// @Accept json
// @Produce json
// @Param page query number false "Page"
// @Param limit query number false "Limit"
// @Security ApiKeyAuth
// @Success 200 {object} response.Message
// @Router /users/auth/sessions [get]
func Sessions(c echo.Context) (err error) {

	user := request.AuthenticatedUser(c)

	page, limit := paginate.HandleQueries(c)
	sessions, err := repository.GetUserSessions(user.ID, page, limit)
	if err != nil {
		return err
	}

	return r.CustomErrorJson(http.StatusOK, sessions, c)
}

// Signout godoc
// @Summary Delete session
// @Tags session
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id body string false "Session ID"
// @Success 200 {object} response.Message
// @Router /users/auth/signout [delete]
func Signout(c echo.Context) (err error) {

	user := request.AuthenticatedUser(c)

	params := new(SignoutSchema)
	if err = request.GetInputs(c, params); err != nil {
		return err
	}

	session, err := repository.SessionFindByID(params.ID)
	if err != nil {
		return err
	}

	if session.UserID != user.ID {
		return errors.ErrAccessDenied
	}

	if err = repository.TerminateSession(user.SID); err != nil {
		return err
	}

	return errors.ErrSuccess
}
