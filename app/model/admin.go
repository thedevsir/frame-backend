package model

import (
	"time"

	"github.com/zebresel-com/mongodm"
)

const AdminCollection = "Admin"

type Admin struct {
	mongodm.DocumentBase `json:",inline" bson:",inline"`

	Username     string    `json:"username" bson:"username"`
	Password     string    `json:"password" bson:"password"`
	Session      string    `json:"session" bson:"session"`
	LoginAt      time.Time `json:"loginAt" bson:"loginAt"`
	LastActivity time.Time `json:"lastActivity" bson:"lastActivity"`
	IsActive     bool      `json:"isActive" bson:"isActive"`
}
