package objectId

import (
	"bytes"
	"io/ioutil"

	"github.com/labstack/echo"
	"github.com/thedevsir/frame-backend/services/errors"
	"gopkg.in/mgo.v2/bson"
)

func Check(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		var bodyBytes []byte
		if c.Request().Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request().Body)
		}

		c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		params := echo.Map{}
		if err := c.Bind(&params); err == nil {
			id, ok := params["id"]
			if ok && bson.IsObjectIdHex(id.(string)) != true {
				return errors.ErrObjectId
			}
		}

		paramID := c.Param("id")
		if paramID != "" && bson.IsObjectIdHex(paramID) != true {
			return errors.ErrObjectId
		}

		c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		return next(c)
	}
}
