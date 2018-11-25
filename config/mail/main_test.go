package mail

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thedevsir/frame-backend/config"
	"github.com/thedevsir/frame-backend/services/utils"
	gomail "gopkg.in/gomail.v2"
)

func TestMailConnection(t *testing.T) {

	t.Run("WithOutConfigInfo", func(t *testing.T) {
		assert.Panics(t, func() { Composer() })
	})

	t.Run("SendTestEmail", func(t *testing.T) {
		config.Composer("../../.env")
		Composer()
		m := gomail.NewMessage()
		m.SetHeader("From", config.EmailFrom)
		m.SetHeader("To", "freshmanlimited@gmail.com")
		m.SetHeader("Subject", "SendTestEmail")
		m.SetBody("text/plain", "...")
		assert.NoError(t, Connection.DialAndSend(m))
	})
}
