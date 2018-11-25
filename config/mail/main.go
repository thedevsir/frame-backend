package mail

import (
	"github.com/thedevsir/frame-backend/config"
	gomail "gopkg.in/gomail.v2"
)

var Connection *gomail.Dialer

func Composer() {

	Connection = gomail.NewDialer(
		config.SMTPHost,
		config.SMTPPort,
		config.SMTPUsername,
		config.SMTPPassword,
	)
}
