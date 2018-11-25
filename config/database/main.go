package database

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/zebresel-com/mongodm"
	"github.com/thedevsir/frame-backend/app/model"
	"github.com/thedevsir/frame-backend/services/utils"
	"gopkg.in/mgo.v2"
)

var Connection *mongodm.Connection

type Composer struct {
	Locals   string
	Addrs    []string
	Database string
	Username string
	Password string
	Source   string
}

func (c Composer) Shoot() {

	file, err := ioutil.ReadFile(c.Locals)

	if err != nil {
		panic(err)
	}

	var localMap map[string]map[string]string
	json.Unmarshal(file, &localMap)

	dbConfig := &mongodm.Config{
		DialInfo: &mgo.DialInfo{
			Addrs:    c.Addrs,
			Timeout:  3 * time.Second,
			Database: c.Database,
			Username: c.Username,
			Password: c.Password,
			Source:   c.Source,
		},
		Locals: localMap["en-US"],
	}

	Connection, err = mongodm.Connect(dbConfig)

	if err != nil {
		panic(err)
	}

	models := map[string]mongodm.IDocumentBase{
		"authAttempts": &model.AuthAttempt{},
		"sessions":     &model.Session{},
		"users":        &model.User{},
		"admin":        &model.Admin{},
	}

	for k, v := range models {
		Connection.Register(v, k)
	}

	collections, err := Connection.Session.DB(c.Database).CollectionNames()
	if err != nil {
		panic(err)
	}

	// Indexes
	if !utils.Contains(collections, "sessions") {

		index := mgo.Index{
			Key:         []string{"expireAt"},
			ExpireAfter: time.Hour * 24,
		}

		err = Connection.Model(model.SessionCollection).EnsureIndex(index)
		if err != nil {
			panic(err)
		}
	}
}
