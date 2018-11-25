package storage

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thedevsir/frame-backend/config"
	"github.com/thedevsir/frame-backend/services/errors"
	"github.com/thedevsir/frame-backend/services/test"
)

func TestStorage(t *testing.T) {

	config.Composer("../../.env")
	objName := "test_" + strconv.FormatInt(time.Now().Unix(), 10)

	t.Run("PrepareStorage", func(t *testing.T) {
		assert.NotPanics(t, Composer)
	})

	t.Run("Put", func(t *testing.T) {
		r := base64.NewDecoder(base64.StdEncoding, strings.NewReader(test.TestPNGBase64))
		assert.NoError(t, Put(objName, "avatar", r, -1, "image/png"))
	})

	t.Run("GetURL", func(t *testing.T) {
		url := GetURL(objName, "avatar")
		resp, err := http.Get(url)
		assert.NoError(t, err)
		resp.Body.Close()
		assert.Equal(t, "image/png", resp.Header.Get("Content-Type"))
	})

	t.Run("Delete", func(t *testing.T) {

		t.Run("Success", func(t *testing.T) {
			assert.NoError(t, Delete(objName, "avatar"))
		})

		t.Run("ObjectNotFound", func(t *testing.T) {
			assert.Equal(t, errors.ErrObjectNotFound, Delete(objName, "avatar"))
		})
	})
}
