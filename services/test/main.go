package test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo"
	"github.com/thedevsir/frame-backend/config"
	"github.com/thedevsir/frame-backend/config/database"
	"github.com/thedevsir/frame-backend/services/validation"
	validator "gopkg.in/go-playground/validator.v9"
)

const (
	TestGIFBase64 = `R0lGODlhAQABAIAAAP///wAAACwAAAAAAQABAAACAkQBADs=`
	TestPNGBase64 = `iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAAAAAA6fptVAAAACklEQVQYV2P4DwABAQEAWk1v8QAAAABJRU5ErkJggg==`
)

func DBComposer(localPath string) database.Composer {

	db := &database.Composer{
		Locals:   localPath,
		Addrs:    strings.Split(config.DBAddress, ","),
		Database: config.DBName + "_Test",
		Username: config.DBUsername,
		Password: config.DBPassword,
		Source:   config.DBSource + "_Test",
	}

	return *db
}

func MakeRequest(method, userJSON string) (echo.Context, *httptest.ResponseRecorder) {

	e := echo.New()
	e.Validator = &validation.DataValidator{ValidatorData: validator.New()}
	req := httptest.NewRequest(method, "/", nil)
	if userJSON != "" {
		req = httptest.NewRequest(method, "/", strings.NewReader(userJSON))
	}
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

func MakeFormdataRequest(method string, r io.Reader) (echo.Context, *httptest.ResponseRecorder) {

	e := echo.New()
	e.Validator = &validation.DataValidator{ValidatorData: validator.New()}
	buf := bytes.NewBuffer(nil)
	mw := multipart.NewWriter(buf)
	w, _ := mw.CreateFormFile("picture", "sample-picture")
	io.Copy(w, r)
	mw.Close()
	req := httptest.NewRequest(method, "/", buf)
	req.Header.Set(echo.HeaderContentType, mw.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}
