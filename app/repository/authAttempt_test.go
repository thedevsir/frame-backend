package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zebresel-com/mongodm"
	"github.com/thedevsir/frame-backend/app/model"
	"github.com/thedevsir/frame-backend/config"
	"github.com/thedevsir/frame-backend/config/database"
	"github.com/thedevsir/frame-backend/services/errors"
	"github.com/thedevsir/frame-backend/services/test"
	"github.com/thedevsir/frame-backend/services/utils"
)

var authAttemptCollection *mongodm.Model

func authAttemptBeforeTest() {

	config.Composer("../../.env")

	db := test.DBComposer("../../resource/locals/locals.json")
	db.Shoot()

	authAttemptCollection = database.Connection.Model(model.AuthAttemptCollection)
	authAttemptCollection.RemoveAll(nil)
}

func authAttemptAfterTest() {
	authAttemptCollection.RemoveAll(nil)
}

func TestAuthAttempt(t *testing.T) {

	authAttemptBeforeTest()
	defer authAttemptAfterTest()
	var ip, username = "127.0.0.1", "Irani"

	t.Run("SubmitAttempt", func(t *testing.T) {
		for i := 0; i <= 5; i++ {
			assert.Nil(t, SubmitAttempt(ip, username))
		}
	})

	t.Run("MaximumAttemptsNotReached", func(t *testing.T) {

		config.AbuseIP = 10
		config.AbuseIPUsername = 10

		assert.Nil(t, CheckAbuse(ip, username))
	})

	t.Run("MaximumAttemptsReached", func(t *testing.T) {

		config.AbuseIP = 3
		config.AbuseIPUsername = 3

		err := CheckAbuse(ip, username)
		assert.Error(t, err)
		assert.Equal(t, errors.ErrAttemptsReached, err)
	})
}
