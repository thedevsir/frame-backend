package mail

import (
	"fmt"

	"github.com/matcornic/hermes"
	"github.com/thedevsir/frame-backend/config"
)

type Verify struct {
	Username     string
	EmailAddress string
	Token        string
}

func (w *Verify) Name() string {
	return "verify"
}

func (w *Verify) Email() hermes.Email {
	return hermes.Email{
		Body: hermes.Body{
			Name: w.Username,
			Intros: []string{
				"Welcome to " + config.EmailAppName + "! We're very excited to have you on board.",
			},
			Dictionary: []hermes.Entry{
				{Key: "Username", Value: w.Username},
				{Key: "Email", Value: w.EmailAddress},
			},
			Actions: []hermes.Action{
				{
					Instructions: "To get started with " + config.EmailAppName + ", please click here:",
					Button: hermes.Button{
						Text: "Confirm your account",
						Link: fmt.Sprintf(config.EmailVerifyLink, w.Token),
					},
				},
			},
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
		},
	}
}
