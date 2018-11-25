package errors

import (
	"net/http"

	"github.com/labstack/echo"
)

var (
	ErrInternal           = echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	ErrObjectId           = echo.NewHTTPError(http.StatusBadRequest, "invalid objectId")
	ErrObjectNotFound     = echo.NewHTTPError(http.StatusNotFound, "requested object not found")
	ErrPictureNotValid    = echo.NewHTTPError(http.StatusBadRequest, "picture data is not valid")
	ErrPictureTooLarge    = echo.NewHTTPError(http.StatusRequestEntityTooLarge, "picture is too large")
	ErrSuccess            = echo.NewHTTPError(http.StatusOK, "success")
	ErrCreated            = echo.NewHTTPError(http.StatusCreated, "success")
	ErrInvalidParams      = echo.NewHTTPError(http.StatusBadRequest, "input param(s) are not valid")
	ErrAccessDenied       = echo.NewHTTPError(http.StatusForbidden, "can not access to this resource/method")
	ErrAccountVerified    = echo.NewHTTPError(http.StatusBadRequest, "your account has already been verified")
	ErrUserEmailNotFound  = echo.NewHTTPError(http.StatusNotFound, "your email address was not found")
	ErrTokenIsNotValid    = echo.NewHTTPError(http.StatusNotFound, "token parse error")
	ErrUserNotFound       = echo.NewHTTPError(http.StatusNotFound, "requested user not found")
	ErrAdminNotFound      = echo.NewHTTPError(http.StatusNotFound, "requested admin not found")
	ErrUsernameExists     = echo.NewHTTPError(http.StatusConflict, "requested username is already exists")
	ErrEmailExists        = echo.NewHTTPError(http.StatusConflict, "account with this email is alreay registered")
	ErrAttemptsReached    = echo.NewHTTPError(http.StatusRequestTimeout, "maximum number of auth attempts reached")
	ErrInvalidCredentials = echo.NewHTTPError(http.StatusForbidden, "credentials are invalid")
	ErrSessionNotFound    = echo.NewHTTPError(http.StatusNotFound, "session not found")
)
