package encrypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBcrypt(t *testing.T) {

	var result, passwd = "", "12345"
	var err error

	t.Run("Hash", func(t *testing.T) {
		result, err = Hash(passwd)
		assert.Nil(t, err)
	})

	t.Run("CheckHash", func(t *testing.T) {
		err := CheckHash(passwd, result)
		assert.True(t, err)
	})
}
