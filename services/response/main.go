package response

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type Message struct {
	Message   interface{} `json:"Message"`
	ErrorCode int         `json:"Code"`
}

func CustomError(errorCode int, errorMessage string) *echo.HTTPError {

	return echo.NewHTTPError(errorCode, errorMessage)
}

func CustomErrorJson(errorCode int, errorMessage interface{}, c echo.Context) error {

	m := Message{
		Message:   errorMessage,
		ErrorCode: errorCode,
	}

	return c.JSON(errorCode, m)
}

func CustomErrorWithHeader(errorCode int, errorMessage string, headers map[string]string, c echo.Context) error {

	m := Message{
		Message:   errorMessage,
		ErrorCode: errorCode,
	}

	for k, v := range headers {
		c.Response().Header().Set(k, v)
	}

	return c.JSON(errorCode, m)
}

func ValidationError(err error) *echo.HTTPError {

	return echo.NewHTTPError(http.StatusBadRequest, err.Error())
}

func ErrorHandler(err error, c echo.Context) {

	m := Message{
		Message:   err.Error(),
		ErrorCode: http.StatusInternalServerError,
	}

	if he, ok := err.(*echo.HTTPError); ok {
		m.ErrorCode = he.Code
		m.Message = fmt.Sprint(he.Message)
	}

	// Send response
	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD {
			err = c.NoContent(m.ErrorCode)
		} else {
			err = c.JSON(m.ErrorCode, m)
		}
	}
}
