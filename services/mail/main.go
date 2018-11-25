package mail

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/thedevsir/frame-backend/config"
	Mail "github.com/thedevsir/frame-backend/config/mail"
	"github.com/thedevsir/frame-backend/resource/views/mail"
	gomail "gopkg.in/gomail.v2"
)

type EmailToken struct {
	Action   string `json:"action"`
	UserID   string `json:"userId"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func MakeEmailToken(action, userID, username, email string, secret []byte) (string, error) {

	claims := EmailToken{
		action,
		userID,
		email,
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	signedClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := signedClaims.SignedString(secret)
	if err != nil {
		return "", err
	}

	return signedString, nil
}

func SendVerficationMail(username, email, token string) error {

	verifyBody := &mail.Verify{
		Username:     username,
		EmailAddress: email,
		Token:        token,
	}

	emailBody, emailText, err := mail.GenerateTemplate(verifyBody.Email())
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", config.EmailFrom)
	m.SetHeader("To", verifyBody.EmailAddress)
	m.SetHeader("Subject", "Confirm your account")
	m.SetBody("text/plain", emailText)
	m.AddAlternative("text/html", emailBody)

	if err := Mail.Connection.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func SendResetMail(username, email, token string) error {

	resetBody := mail.Reset{
		Username:     username,
		EmailAddress: email,
		Token:        token,
	}

	emailBody, emailText, err := mail.GenerateTemplate(resetBody.Email())
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", config.EmailFrom)
	m.SetHeader("To", resetBody.EmailAddress)
	m.SetHeader("Subject", "Reset your password")
	m.SetBody("text/plain", emailText)
	m.AddAlternative("text/html", emailBody)

	if err := Mail.Connection.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
