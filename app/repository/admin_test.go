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

func TestAdmin(t *testing.T) {

	adminBeforeTest()
	defer adminAfterTest()
	var err error
	var adminID string
	var username, password = "admin", "12345"

	t.Run("CreateAdmin", func(t *testing.T) {

		t.Run("Success", func(t *testing.T) {
			adminID, err = CreateAdmin(username, password)
			assert.Nil(t, err)
		})

		t.Run("UsernameExists", func(t *testing.T) {
			_, err = CreateAdmin(username, password)
			assert.Equal(t, errors.ErrUsernameExists, err)
		})
	})

	t.Run("FindAdminByCredentials", func(t *testing.T) {

		t.Run("Success", func(t *testing.T) {

			result, err := FindAdminByCredentials(username, password)
			if assert.Nil(t, err) {
				assert.Equal(t, username, result.Username)
			}
		})

		t.Run("UsernameNotFound", func(t *testing.T) {

			_, err = FindAdminByCredentials("admin2", password)
			assert.Equal(t, errors.ErrAdminNotFound, err)
		})

		t.Run("PasswordHasNotMatch", func(t *testing.T) {

			_, err := FindAdminByCredentials(username, "wrongpassword")
			assert.Equal(t, errors.ErrInvalidCredentials, err)
		})
	})

	t.Run("AdminSession", func(t *testing.T) {

		var key string

		t.Run("SetAdminSession", func(t *testing.T) {

			t.Run("AdminNotFound", func(t *testing.T) {
				key, err = SetAdminSession(bson.NewObjectId().Hex())
				assert.Equal(t, errors.ErrAdminNotFound, err)
			})

			t.Run("Success", func(t *testing.T) {
				key, err = SetAdminSession(adminID)
				assert.Nil(t, err)
			})
		})

		t.Run("CheckAdminSession", func(t *testing.T) {

			t.Run("AdminNotFound", func(t *testing.T) {
				assert.Equal(t, errors.ErrAdminNotFound, CheckAdminSession(bson.NewObjectId().Hex(), key))
			})

			t.Run("FakeSession", func(t *testing.T) {
				assert.Equal(t, errors.ErrInvalidCredentials, CheckAdminSession(adminID, "wrongKey"))
			})

			t.Run("Success", func(t *testing.T) {
				assert.Nil(t, CheckAdminSession(adminID, key))
			})
		})

		t.Run("TerminateAdminSession", func(t *testing.T) {

			t.Run("AdminNotFound", func(t *testing.T) {
				assert.Equal(t, errors.ErrAdminNotFound, TerminateAdminSession(bson.NewObjectId().Hex()))
			})

			t.Run("Success", func(t *testing.T) {
				assert.Nil(t, TerminateAdminSession(adminID))
			})
		})
	})

	t.Run("GetAdminByID", func(t *testing.T) {

		t.Run("Success", func(t *testing.T) {

			_, err := GetAdminByID(adminID)
			assert.Nil(t, err)
		})

		t.Run("UserNotFound", func(t *testing.T) {

			_, err := GetAdminByID(bson.NewObjectId().Hex())
			assert.Equal(t, errors.ErrAdminNotFound, err)
		})
	})

	t.Run("AdminChangeUsername", func(t *testing.T) {

		username = "root"

		t.Run("AdminNotFound", func(t *testing.T) {

			err := AdminChangeUsername(bson.NewObjectId().Hex(), username)
			assert.Equal(t, errors.ErrAdminNotFound, err)
		})

		t.Run("Success", func(t *testing.T) {

			err := AdminChangeUsername(adminID, username)
			assert.Nil(t, err)
		})

		t.Run("UsernameExists", func(t *testing.T) {

			err := AdminChangeUsername(adminID, username)
			assert.Equal(t, errors.ErrUsernameExists, err)
		})
	})

	t.Run("AdminChangePassword", func(t *testing.T) {

		password = "12345678"

		t.Run("Success", func(t *testing.T) {

			err := AdminChangePassword(adminID, password)
			assert.Nil(t, err)
		})

		t.Run("UserNotFound", func(t *testing.T) {

			err := AdminChangePassword(bson.NewObjectId().Hex(), "12345678")
			assert.Equal(t, errors.ErrAdminNotFound, err)
		})
	})

	t.Run("AdminSessionUpdateLastActivity", func(t *testing.T) {

		t.Run("Success", func(t *testing.T) {
			assert.Nil(t, AdminSessionUpdateLastActivity(adminID))
		})

		t.Run("FakeSession", func(t *testing.T) {
			assert.Equal(t, errors.ErrAdminNotFound, AdminSessionUpdateLastActivity(bson.NewObjectId().Hex()))
		})
	})

	t.Run("ChangeAdminStatus", func(t *testing.T) {

		t.Run("Success", func(t *testing.T) {
			assert.Nil(t, ChangeAdminStatus(adminID, false))
		})

		t.Run("AdminNotFound", func(t *testing.T) {
			assert.Equal(t, errors.ErrAdminNotFound, ChangeAdminStatus(bson.NewObjectId().Hex(), false))
		})
	})

	t.Run("GetAdmins", func(t *testing.T) {

		t.Run("Success", func(t *testing.T) {
			_, err := GetAdmins(1, 10)
			assert.Nil(t, err)
		})

		adminCollection.RemoveAll(nil)

		t.Run("AdminNotFound", func(t *testing.T) {
			_, err := GetAdmins(1, 10)
			assert.Equal(t, errors.ErrAdminNotFound, err)
		})
	})
}
