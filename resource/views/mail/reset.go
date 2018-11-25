package mail

import (
	"fmt"

	"github.com/matcornic/hermes"
	"github.com/thedevsir/frame-backend/config"
)

type Reset struct {
	Username     string
	EmailAddress string
	Token        string
}

func (r *Reset) Name() string {
	return "reset"
}

func (r *Reset) Email() hermes.Email {
	return hermes.Email{
		Body: hermes.Body{
			Name: r.Username,
			Intros: []string{
				"You have received this email because a password reset request for \"" + r.EmailAddress + "\" account was received.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Click the button below to reset your password:",
					Button: hermes.Button{
						Color: "#DC4D2F",
						Text:  "Reset your password",
						Link:  fmt.Sprintf(config.EmailResetLink, r.Token),
					},
				},
			},
			Outros: []string{
				"If you did not request a password reset, no further action is required on your part.",
			},
			Signature: "Thanks",
		},
	}
}
