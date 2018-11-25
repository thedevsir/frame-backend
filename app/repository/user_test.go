package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zebresel-com/mongodm"
	"github.com/thedevsir/frame-backend/app/model"
	"github.com/thedevsir/frame-backend/config/database"
	"github.com/thedevsir/frame-backend/services/errors"
	"github.com/thedevsir/frame-backend/services/test"
	"github.com/thedevsir/frame-backend/services/utils"
	"gopkg.in/mgo.v2/bson"
)

var userCollection *mongodm.Model

func userBeforeTest() {

	config.Composer("../../.env")

	db := test.DBComposer("../../resource/locals/locals.json")
	db.Shoot()

	userCollection = database.Connection.Model(model.UserCollection)
	userCollection.RemoveAll(nil)
}

func userAfterTest() {
	userCollection.RemoveAll(nil)
}

func TestUser(t *testing.T) {

	userBeforeTest()
	defer userAfterTest()

	var err error
	var insertedData *model.User
	var username, password, email = "Irani", "12345", "freshmanlimited@gmail.com"

	t.Run("CreateUser", func(t *testing.T) {
		insertedData, err = CreateUser(username, password, email)
		assert.Nil(t, err)
	})

	t.Run("Activation", func(t *testing.T) {

		t.Run("Success", func(t *testing.T) {
			assert.Nil(t, UserActivation(insertedData.Id.Hex()))
		})

		t.Run("UserNotFound", func(t *testing.T) {
			assert.Equal(t, errors.ErrUserNotFound, UserActivation(bson.NewObjectId().Hex()))
		})
	})

	t.Run("CheckEmail", func(t *testing.T) {

		t.Run("EmailNotFound", func(t *testing.T) {
			_, err = CheckEmail("fake@service.domain")
			assert.Nil(t, err)
		})

		t.Run("AlreadyExists", func(t *testing.T) {
			_, err = CheckEmail(email)
			assert.Equal(t, err, errors.ErrEmailExists)
		})
	})

	t.Run("CheckUsername", func(t *testing.T) {

		t.Run("UserNotFound", func(t *testing.T) {
			_, err = CheckUsername("fakeuser")
			assert.Nil(t, err)
		})

		t.Run("AlreadyExists", func(t *testing.T) {
			_, err = CheckUsername(username)
			assert.Equal(t, err, errors.ErrUsernameExists)
		})
	})

	t.Run("ChangePassword", func(t *testing.T) {

		password = "12345678"

		t.Run("UserNotFound", func(t *testing.T) {
			assert.Equal(t, errors.ErrUserNotFound, ChangePassword(bson.NewObjectId().Hex(), password, false))
		})

		t.Run("Success", func(t *testing.T) {
			assert.Nil(t, ChangePassword(insertedData.Id.Hex(), password, false))
		})
	})

	t.Run("ChangeUsername", func(t *testing.T) {

		username = "amirhossein"

		t.Run("UserNotFound", func(t *testing.T) {
			assert.Equal(t, errors.ErrUserNotFound, ChangeUsername(bson.NewObjectId().Hex(), username, false))
		})

		t.Run("Success", func(t *testing.T) {
			assert.Nil(t, ChangeUsername(insertedData.Id.Hex(), username, false))
		})
	})

	t.Run("ChangeEmail", func(t *testing.T) {

		email = "amirhossein@rock.age"

		t.Run("UserNotFound", func(t *testing.T) {
			assert.Equal(t, errors.ErrUserNotFound, ChangeEmail(bson.NewObjectId().Hex(), email, false))
		})

		t.Run("Success", func(t *testing.T) {
			assert.Nil(t, ChangeEmail(insertedData.Id.Hex(), email, false))
		})
	})

	t.Run("GetAccountInfo", func(t *testing.T) {

		t.Run("UserNotFound", func(t *testing.T) {
			_, err = GetAccountInfo(bson.NewObjectId().Hex())
			assert.Equal(t, errors.ErrUserNotFound, err)
		})

		t.Run("Success", func(t *testing.T) {
			_, err = GetAccountInfo(insertedData.Id.Hex())
			assert.Nil(t, err)
		})
	})

	t.Run("FindUserByCredentials", func(t *testing.T) {

		t.Run("Success", func(t *testing.T) {

			result, err := FindUserByCredentials(username, password)
			if assert.Nil(t, err) {
				assert.Equal(t, username, result.Username)
				assert.Equal(t, email, result.Email)
			}
		})

		t.Run("UsernameNotFound", func(t *testing.T) {

			_, err = FindUserByCredentials("Irani@gmail.com", "12345")
			assert.Equal(t, errors.ErrUserNotFound, err)
		})

		t.Run("PasswordHasNotMatch", func(t *testing.T) {

			_, err := FindUserByCredentials(username, "123")
			assert.Equal(t, errors.ErrInvalidCredentials, err)
		})
	})

	t.Run("GetUserByID", func(t *testing.T) {

		t.Run("UserNotFound", func(t *testing.T) {
			_, err = GetUserByID(bson.NewObjectId().Hex())
			assert.Equal(t, errors.ErrUserNotFound, err)
		})

		t.Run("Success", func(t *testing.T) {
			_, err = GetUserByID(insertedData.Id.Hex())
			assert.Nil(t, err)
		})
	})

	t.Run("ChangeUserStatus", func(t *testing.T) {

		t.Run("UserNotFound", func(t *testing.T) {
			assert.Equal(t, errors.ErrUserNotFound, ChangeUserStatus(bson.NewObjectId().Hex(), false))
		})

		t.Run("Success", func(t *testing.T) {
			assert.Nil(t, ChangeUserStatus(insertedData.Id.Hex(), false))
		})
	})

	t.Run("GetUserByIDFromAdmin", func(t *testing.T) {

		t.Run("UserNotFound", func(t *testing.T) {
			_, err = GetUserByIDFromAdmin(bson.NewObjectId().Hex())
			assert.Equal(t, errors.ErrUserNotFound, err)
		})

		t.Run("Success", func(t *testing.T) {
			_, err = GetUserByIDFromAdmin(insertedData.Id.Hex())
			assert.Nil(t, err)
		})
	})

	t.Run("GetUsers", func(t *testing.T) {

		t.Run("Success", func(t *testing.T) {
			_, err := GetUsers(1, 10)
			assert.Nil(t, err)
		})

		userCollection.RemoveAll(nil)

		t.Run("UserNotFound", func(t *testing.T) {
			_, err := GetUsers(1, 10)
			assert.Equal(t, errors.ErrUserNotFound, err)
		})
	})
}
