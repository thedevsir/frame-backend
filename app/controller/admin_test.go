package controller

import (
	"net/http"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/zebresel-com/mongodm"
	"github.com/thedevsir/frame-backend/app/model"
	"github.com/thedevsir/frame-backend/app/repository"
	"github.com/thedevsir/frame-backend/config"
	"github.com/thedevsir/frame-backend/config/database"
	"github.com/thedevsir/frame-backend/services/auth"
	"github.com/thedevsir/frame-backend/services/errors"
	j "github.com/thedevsir/frame-backend/services/jwt"
	"github.com/thedevsir/frame-backend/services/test"
)

var adminCollection *mongodm.Model

func adminBeforeTest() {

	config.Composer("../../.env")

	db := test.DBComposer("../../resource/locals/locals.json")
	db.Shoot()

	adminCollection = database.Connection.Model(model.AdminCollection)
	adminCollection.RemoveAll(nil)
}

func adminAfterTest() {
	adminCollection.RemoveAll(nil)
}

func TestAdminUser(t *testing.T) {

	adminBeforeTest()
	defer adminAfterTest()

	secret := []byte("secret")
	username, password := "admin", "12345678"
	adminID, _ := repository.CreateAdmin(username, password)

	token := &auth.AdminToken{
		Session: bson.NewObjectId().Hex(),
		ID:      adminID,
	}
	tc, _ := token.Create(secret)

	t.Run("AdminSignin", func(t *testing.T) {

		t.Run("BindErr", func(t *testing.T) {

			c, _ := test.MakeRequest(echo.POST, "")
			assert.Equal(t, errors.ErrInvalidParams, AdminSignin(c))
		})

		t.Run("ValidateErr", func(t *testing.T) {

			JSONData := `{"username":"ir","password":"12345","email":"fakeMail"}`
			c, _ := test.MakeRequest(echo.POST, JSONData)
			assert.Regexp(t, "code=400", AdminSignin(c))
		})

		t.Run("UserNotFound", func(t *testing.T) {

			JSONData := `{"username":"admin2","password":"` + password + `"}`
			c, _ := test.MakeRequest(echo.POST, JSONData)
			assert.Equal(t, errors.ErrAdminNotFound, AdminSignin(c))
		})

		t.Run("InvalidCredentials", func(t *testing.T) {

			JSONData := `{"username":"` + username + `","password":"87654321"}`
			c, _ := test.MakeRequest(echo.POST, JSONData)
			assert.Equal(t, errors.ErrInvalidCredentials, AdminSignin(c))
		})

		t.Run("Success", func(t *testing.T) {

			JSONData := `{"username":"` + username + `","password":"` + password + `"}`
			c, rec := test.MakeRequest(echo.POST, JSONData)
			if assert.NoError(t, AdminSignin(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
			}
		})
	})

	t.Run("AdminLogout", func(t *testing.T) {

		c, _ := test.MakeRequest(echo.DELETE, "")
		tokenParsed, _ := j.ParseJWT(tc, secret)
		c.Set("user", tokenParsed)

		assert.Equal(t, errors.ErrSuccess, AdminSignout(c))
	})
}
