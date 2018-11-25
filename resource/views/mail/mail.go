package mail

import (
	"fmt"
	"strconv"
	"time"

	"github.com/matcornic/hermes"
	"github.com/thedevsir/frame-backend/config"
)

func GenerateTemplate(email hermes.Email) (string, string, error) {

	getYear := strconv.Itoa(time.Now().Year())
	h := hermes.Hermes{
		Product: hermes.Product{
			Name:        config.EmailThemeName,
			Link:        config.EmailThemeLink,
			Logo:        config.EmailThemeLogo,
			Copyright:   fmt.Sprintf(config.EmailThemeCopyright, getYear),
			TroubleText: "If you’re having trouble with the button '{ACTION}', copy and paste the URL below into your web browser.",
		},
	}

	emailBody, err := h.GenerateHTML(email)
	if err != nil {
		return "", "", fmt.Errorf("internal server error")
	}

	emailText, err := h.GeneratePlainText(email)
	if err != nil {
		return "", "", fmt.Errorf("internal server error")
	}

	return emailBody, emailText, nil
}
