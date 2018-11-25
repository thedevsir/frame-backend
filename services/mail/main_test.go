package mail

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thedevsir/frame-backend/config"
	"github.com/thedevsir/frame-backend/config/mail"
)

func TestMakeEmailToken(t *testing.T) {

	assert.NotPanics(t, func() {
		MakeEmailToken("verify", "userID", "username", "@", []byte("secret"))
	})
}

func TestSendMail(t *testing.T) {

	config.Composer("../../.env")
	mail.Composer()

	t.Run("SendVerficationMail", func(t *testing.T) {
		assert.NoError(t, SendVerficationMail("username", "freshmanlimited@gmail.com", "token"))
	})

	t.Run("SendResetMail", func(t *testing.T) {
		assert.NoError(t, SendResetMail("username", "freshmanlimited@gmail.com", "token"))
	})
}
