package controller

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"testing"

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
	"github.com/thedevsir/frame-backend/services/mail"
	"github.com/thedevsir/frame-backend/services/storage"
	"github.com/thedevsir/frame-backend/services/test"
)

var userCollection *mongodm.Model

func userBeforeTest() (*model.User, *jwt.Token) {

	config.Composer("../../.env")

	db := test.DBComposer("../../resource/locals/locals.json")
	db.Shoot()

	storage.Composer()

	userCollection = database.Connection.Model(model.UserCollection)
	userCollection.RemoveAll(nil)

	// Mock
	_SendVerficationMail = func(username, email, token string) error {
		return nil
	}
	_SendResetMail = func(username, email, token string) error {
		return nil
	}

	username, password, email := "Amir", "12345678", "live@forumX.com"

	// CreateUser
	user, _ := repository.CreateUser(username, password, email)

	secret := []byte("secret")
	token := &auth.UserToken{
		Session: "Session",
		SID:     "SID",
		ID:      user.Id.Hex(),
	}
	tc, _ := token.Create(secret)
	tokenParsed, _ := j.ParseJWT(tc, secret)

	return user, tokenParsed
}

func userAfterTest() {
	userCollection.RemoveAll(nil)
}

func TestSignup(t *testing.T) {

	user, _ := userBeforeTest()
	defer userAfterTest()

	t.Run("BindErr", func(t *testing.T) {

		c, _ := test.MakeRequest(echo.POST, "")
		assert.Equal(t, errors.ErrInvalidParams, Signup(c))
	})

	t.Run("ValidateErr", func(t *testing.T) {

		JSONData := `{"username":"ir","password":"12345","email":"fakeMail"}`
		c, _ := test.MakeRequest(echo.POST, JSONData)
		assert.Regexp(t, "code=400", Signup(c))
	})

	t.Run("DuplicateUsername", func(t *testing.T) {

		JSONData := `{"username":"` + user.Username + `","password":"12345678","email":"live@forumX.com"}`
		c, _ := test.MakeRequest(echo.POST, JSONData)
		assert.Equal(t, errors.ErrUsernameExists, Signup(c))
	})

	t.Run("DuplicateEmail", func(t *testing.T) {

		JSONData := `{"username":"Amir2","password":"12345678","email":"` + user.Email + `"}`
		c, _ := test.MakeRequest(echo.POST, JSONData)
		assert.Equal(t, errors.ErrEmailExists, Signup(c))
	})

	t.Run("Success", func(t *testing.T) {

		username := "Amir2"
		email := "live2@forumX.com"

		JSONData := `{"username":"` + username + `","password":"12345678","email":"` + email + `"}`
		c, _ := test.MakeRequest(echo.POST, JSONData)
		assert.Equal(t, errors.ErrCreated, Signup(c))
	})
}

func TestResend(t *testing.T) {

	user, _ := userBeforeTest()
	defer userAfterTest()

	t.Run("BindErr", func(t *testing.T) {

		c, _ := test.MakeRequest(echo.POST, "")
		assert.Equal(t, errors.ErrInvalidParams, Resend(c))
	})

	t.Run("ValidateErr", func(t *testing.T) {

		JSONData := `{"key":"value"}`
		c, _ := test.MakeRequest(echo.POST, JSONData)
		assert.Regexp(t, "code=400", Resend(c))
	})

	t.Run("EmailNotFound", func(t *testing.T) {

		JSONData := `{"email":"fake@gmail.com"}`
		c, _ := test.MakeRequest(echo.POST, JSONData)
		assert.Equal(t, errors.ErrUserEmailNotFound, Resend(c))
	})

	t.Run("Success", func(t *testing.T) {

		JSONData := `{"email":"live@forumX.com"}`
		c, _ := test.MakeRequest(echo.POST, JSONData)
		assert.Equal(t, errors.ErrSuccess, Resend(c))
	})

	t.Run("VerfiedBefore", func(t *testing.T) {

		repository.UserActivation(user.Id.Hex())
		JSONData := `{"email":"live@forumX.com"}`
		c, _ := test.MakeRequest(echo.POST, JSONData)
		assert.Equal(t, errors.ErrAccountVerified, Resend(c))
	})
}

func TestVerification(t *testing.T) {

	user, _ := userBeforeTest()
	defer userAfterTest()

	secret := []byte(config.SigningKey)
	token0, _ := mail.MakeEmailToken("verify", user.Id.Hex(), user.Username, user.Email, secret)
	token1, _ := mail.MakeEmailToken("reset", user.Id.Hex(), user.Username, user.Email, secret)

	t.Run("BindErr", func(t *testing.T) {

		c, _ := test.MakeRequest(echo.POST, "")
		assert.Equal(t, errors.ErrInvalidParams, Verification(c))
	})

	t.Run("ValidateErr", func(t *testing.T) {

		JSONData := `{"key":"value"}`
		c, _ := test.MakeRequest(echo.POST, JSONData)
		assert.Regexp(t, "code=400", Verification(c))
	})

	t.Run("WrongAction", func(t *testing.T) {

		JSONData := fmt.Sprintf(`{"token":"%s"}`, token1)
		c, _ := test.MakeRequest(echo.POST, JSONData)
		assert.Equal(t, errors.ErrAccessDenied, Verification(c))
	})

	t.Run("Success", func(t *testing.T) {

		JSONData := fmt.Sprintf(`{"token":"%s"}`, token0)
		c, _ := test.MakeRequest(echo.POST, JSONData)
		assert.Equal(t, errors.ErrSuccess, Verification(c))
	})
}

func TestSignin(t *testing.T) {

	userBeforeTest()
	defer userAfterTest()

	t.Run("BindErr", func(t *testing.T) {

		c, _ := test.MakeRequest(echo.POST, "")
		assert.Equal(t, errors.ErrInvalidParams, Signin(c))
	})

	t.Run("ValidateErr", func(t *testing.T) {

		JSONData := `{"username":"ir"}`
		c, _ := test.MakeRequest(echo.POST, JSONData)
		assert.Regexp(t, "code=400", Signin(c))
	})

	t.Run("UserNotFound", func(t *testing.T) {

		JSONData := `{"username":"ir","password":"12345"}`
		c, _ := test.MakeRequest(echo.POST, JSONData)
		assert.Equal(t, errors.ErrUserNotFound, Signin(c))
	})

	t.Run("LoginSuccessfully", func(t *testing.T) {

		JSONData := `{"username":"amir","password":"12345678"}`
		c, rec := test.MakeRequest(echo.POST, JSONData)

		if assert.NoError(t, Signin(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
}

func TestForgot(t *testing.T) {

	user, _ := userBeforeTest()
	defer userAfterTest()

	t.Run("BindErr", func(t *testing.T) {

		c, _ := test.MakeRequest(echo.POST, "")
		assert.Equal(t, errors.ErrInvalidParams, Forgot(c))
	})

	t.Run("ValidateErr", func(t *testing.T) {

		JSONData := `{"key":"value"}`
		c, _ := test.MakeRequest(echo.POST, JSONData)
		assert.Regexp(t, "code=400", Forgot(c))
	})

	t.Run("EmailNotFound", func(t *testing.T) {

		c, _ := test.MakeRequest(echo.POST, `{"email":"fake@fake.io"}`)
		assert.Equal(t, errors.ErrUserEmailNotFound, Forgot(c))
	})

	t.Run("Success", func(t *testing.T) {

		c, _ := test.MakeRequest(echo.POST, `{"email":"`+user.Email+`"}`)
		assert.Equal(t, errors.ErrSuccess, Forgot(c))
	})
}

func TestReset(t *testing.T) {

	user, _ := userBeforeTest()
	defer userAfterTest()

	secret := []byte(config.SigningKey)
	token0, _ := mail.MakeEmailToken("reset", user.Id.Hex(), user.Username, user.Email, secret)
	token1, _ := mail.MakeEmailToken("verify", user.Id.Hex(), user.Username, user.Email, secret)

	t.Run("BindErr", func(t *testing.T) {

		c, _ := test.MakeRequest(echo.PUT, "")
		assert.Equal(t, errors.ErrInvalidParams, Reset(c))
	})

	t.Run("ValidateErr", func(t *testing.T) {

		JSONData := `{"password":"12345"}`
		c, _ := test.MakeRequest(echo.PUT, JSONData)
		assert.Regexp(t, "code=400", Reset(c))
	})

	t.Run("WrongAction", func(t *testing.T) {

		JSONData := fmt.Sprintf(`{"password":"12345678","token":"%s"}`, token1)
		c, _ := test.MakeRequest(echo.PUT, JSONData)
		assert.Equal(t, errors.ErrAccessDenied, Reset(c))
	})

	t.Run("Success", func(t *testing.T) {

		JSONData := fmt.Sprintf(`{"password":"12345678","token":"%s"}`, token0)
		c, _ := test.MakeRequest(echo.PUT, JSONData)
		assert.Equal(t, errors.ErrSuccess, Reset(c))
	})
}

func TestChangeUsername(t *testing.T) {

	_, tokenParsed := userBeforeTest()
	defer userAfterTest()

	JSONData := `{"username":"Irani"}`
	c, _ := test.MakeRequest(echo.PUT, JSONData)
	c.Set("user", tokenParsed)

	assert.Equal(t, errors.ErrSuccess, ChangeUsername(c))
}

func TestChangePassword(t *testing.T) {

	_, tokenParsed := userBeforeTest()
	defer userAfterTest()

	JSONData := `{"password":"87654321"}`
	c, _ := test.MakeRequest(echo.PUT, JSONData)
	c.Set("user", tokenParsed)

	assert.Equal(t, errors.ErrSuccess, ChangePassword(c))
}

func TestGetAccount(t *testing.T) {

	_, tokenParsed := userBeforeTest()
	defer userAfterTest()

	c, rec := test.MakeRequest(echo.GET, "")
	c.Set("user", tokenParsed)

	if assert.NoError(t, GetAccount(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestChangeEmail(t *testing.T) {

	_, tokenParsed := userBeforeTest()
	defer userAfterTest()

	JSONData := `{"email":"rock@amir.ir"}`
	c, _ := test.MakeRequest(echo.PUT, JSONData)
	c.Set("user", tokenParsed)

	assert.Equal(t, errors.ErrSuccess, ChangeEmail(c))
}

func TestGetUser(t *testing.T) {

	user, tokenParsed := userBeforeTest()
	defer userAfterTest()

	c, rec := test.MakeRequest(echo.GET, "")
	c.Set("user", tokenParsed)
	c.SetParamNames("id")
	c.SetParamValues(user.Id.Hex())

	if assert.NoError(t, GetUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestAvatar(t *testing.T) {

	user, tokenParsed := userBeforeTest()
	defer userAfterTest()

	t.Run("PutAvatar", func(t *testing.T) {

		r := base64.NewDecoder(base64.StdEncoding, strings.NewReader(test.TestPNGBase64))

		c, _ := test.MakeFormdataRequest(echo.PUT, r)
		c.Set("user", tokenParsed)

		assert.Equal(t, errors.ErrSuccess, PutAvatar(c))
	})

	t.Run("DeleteAvatar", func(t *testing.T) {

		c, _ := test.MakeRequest(echo.DELETE, "")
		c.Set("user", tokenParsed)

		assert.Equal(t, errors.ErrSuccess, DeleteAvatar(c))
	})

	// Clean the storage
	storage.Delete(user.Id.Hex(), "avatar")
}
