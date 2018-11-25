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

var sessionCollection *mongodm.Model

func sessionBeforeTest() {

	config.Composer("../../.env")

	db := test.DBComposer("../../resource/locals/locals.json")
	db.Shoot()

	sessionCollection = database.Connection.Model(model.SessionCollection)
	sessionCollection.RemoveAll(nil)
}

func sessionAfterTest() {
	sessionCollection.RemoveAll(nil)
}

func TestSession(t *testing.T) {

	sessionBeforeTest()
	defer sessionAfterTest()
	var err error
	var SID, sess string
	var ip, userID, userAgent = "127.0.0.1", bson.NewObjectId().Hex(), ":::USER-AGENT:::"

	t.Run("CreateSession", func(t *testing.T) {
		assert.NotPanics(t, func() {
			SID, sess, err = SessionCreate(ip, userID, userAgent)
			assert.Nil(t, err)
		})
	})

	t.Run("SessionFindByID", func(t *testing.T) {

		t.Run("Success", func(t *testing.T) {
			_, err := SessionFindByID(SID)
			assert.Nil(t, err)
		})

		t.Run("SessionNotFound", func(t *testing.T) {
			_, err := SessionFindByID(bson.NewObjectId().Hex())
			assert.Error(t, err)
			assert.Equal(t, errors.ErrSessionNotFound, err)
		})
	})

	t.Run("SessionFindByCredentials", func(t *testing.T) {

		t.Run("Success", func(t *testing.T) {
			assert.Nil(t, SessionFindByCredentials(sess, SID))
		})

		t.Run("FakeSession", func(t *testing.T) {
			assert.Equal(t, errors.ErrInvalidCredentials, SessionFindByCredentials("fake", SID))
		})
	})

	t.Run("SessionUpdateLastActivity", func(t *testing.T) {

		t.Run("Success", func(t *testing.T) {
			assert.Nil(t, SessionUpdateLastActivity(SID))
		})

		t.Run("FakeSession", func(t *testing.T) {
			assert.Equal(t, errors.ErrInvalidCredentials, SessionUpdateLastActivity(bson.NewObjectId().Hex()))
		})
	})

	t.Run("GetUserSessions", func(t *testing.T) {

		t.Run("Success", func(t *testing.T) {
			_, err := GetUserSessions(userID, 1, 10)
			assert.Nil(t, err)
		})

		t.Run("FakeUser", func(t *testing.T) {
			_, err := GetUserSessions(bson.NewObjectId().Hex(), 1, 10)
			assert.Equal(t, errors.ErrSessionNotFound, err)
		})
	})

	t.Run("TerminateSession", func(t *testing.T) {

		t.Run("Success", func(t *testing.T) {
			assert.Nil(t, TerminateSession(SID))
		})

		t.Run("FakeSession", func(t *testing.T) {
			assert.Equal(t, errors.ErrSessionNotFound, TerminateSession(bson.NewObjectId().Hex()))
		})
	})

	t.Run("TerminateAllSessions", func(t *testing.T) {
		assert.Nil(t, TerminateAllSessions(userID))
	})
}
