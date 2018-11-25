package mail

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thedevsir/frame-backend/services/utils"
)

func TestTemplates(t *testing.T) {

	utils.Composer("../../../.env")

	t.Run("Verify", func(t *testing.T) {
		verifyBody := Verify{
			Username:     "fakeUser",
			EmailAddress: "fakeEmail",
			Token:        "fakeToken",
		}
		assert.Equal(t, "verify", verifyBody.Name())
		assert.NotPanics(t, func() { GenerateTemplate(verifyBody.Email()) })
	})

	t.Run("Reset", func(t *testing.T) {
		resetBody := Reset{
			Username:     "fakeUser",
			EmailAddress: "fakeEmail",
			Token:        "fakeToken",
		}
		assert.Equal(t, "reset", resetBody.Name())
		assert.NotPanics(t, func() { GenerateTemplate(resetBody.Email()) })
	})
}
