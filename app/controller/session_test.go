package controller

import (
	"fmt"
	"net/http"
	"testing"

	"gopkg.in/mgo.v2/bson"

	jwt "github.com/dgrijalva/jwt-go"
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

var sessionCollection *mongodm.Model

func sessionBeforeTest() (*jwt.Token, string) {

	config.Composer("../../.env")

	db := test.DBComposer("../../resource/locals/locals.json")
	db.Shoot()

	sessionCollection = database.Connection.Model(model.SessionCollection)
	sessionCollection.RemoveAll(nil)

	secret := []byte("secret")
	userID := bson.NewObjectId().Hex()
	sid, key, _ := repository.SessionCreate("127.0.0.1", userID, ":::USER-AGENT:::")

	token := &auth.UserToken{
		key,
		sid,
		userID,
		jwt.StandardClaims{},
	}

	tc, _ := token.Create(secret)
	tokenParsed, _ := j.ParseJWT(tc, secret)

	return tokenParsed, sid
}

func sessionAfterTest() {
	sessionCollection.RemoveAll(nil)
}

func TestSessions(t *testing.T) {

	tokenParsed, _ := sessionBeforeTest()
	defer sessionAfterTest()

	c, rec := test.MakeRequest(echo.GET, "")
	c.Set("user", tokenParsed)

	if assert.NoError(t, Sessions(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestSignout(t *testing.T) {

	tokenParsed, sid := sessionBeforeTest()
	defer sessionAfterTest()

	t.Run("BindErr", func(t *testing.T) {

		c, _ := test.MakeRequest(echo.DELETE, `{,}`)
		c.Set("user", tokenParsed)
		assert.Equal(t, errors.ErrInvalidParams, Signout(c))
	})

	t.Run("ValidateErr", func(t *testing.T) {

		JSONData := `{"id":"",}`
		c, _ := test.MakeRequest(echo.DELETE, JSONData)
		c.Set("user", tokenParsed)
		assert.Regexp(t, "code=400", Signout(c))
	})

	t.Run("SessionNotFound", func(t *testing.T) {

		JSONData := fmt.Sprintf(`{"id":"%s"}`, bson.NewObjectId().Hex())
		c, _ := test.MakeRequest(echo.DELETE, JSONData)
		c.Set("user", tokenParsed)

		assert.Equal(t, errors.ErrSessionNotFound, Signout(c))
	})

	t.Run("SignoutSuccessfully", func(t *testing.T) {

		JSONData := fmt.Sprintf(`{"id":"%s"}`, sid)
		c, _ := test.MakeRequest(echo.DELETE, JSONData)
		c.Set("user", tokenParsed)

		assert.Equal(t, errors.ErrSuccess, Signout(c))
	})
}
