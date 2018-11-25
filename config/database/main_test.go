package database

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thedevsir/frame-backend/config"
	"github.com/thedevsir/frame-backend/services/utils"
)

func TestDatabaseConnection(t *testing.T) {

	config.Composer("../../.env")

	t.Run("LookingForLocalFile", func(t *testing.T) {
		db := &Composer{
			Locals:   "locals.json",
			Addrs:    []string{"fakeAddrs"},
			Database: "",
			Username: "",
			Password: "",
			Source:   "",
		}
		assert.Panics(t, func() { db.Shoot() })
	})

	t.Run("WrongConnection", func(t *testing.T) {
		db := &Composer{
			Locals:   "../../resource/locals/locals.json",
			Addrs:    []string{"fakeAddrs"},
			Database: "",
			Username: "",
			Password: "",
			Source:   "",
		}
		assert.Panics(t, func() { db.Shoot() })
	})

	t.Run("SuccessConnectionWithModel", func(t *testing.T) {
		db := &Composer{
			Locals:   "../../resource/locals/locals.json",
			Addrs:    strings.Split(config.DBAddress, ","),
			Database: config.DBName,
			Username: config.DBUsername,
			Password: config.DBPassword,
			Source:   config.DBSource,
		}
		assert.NotPanics(t, func() {
			db.Shoot()
		})
	})
}
