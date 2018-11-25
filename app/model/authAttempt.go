package model

import (
	"github.com/zebresel-com/mongodm"
)

const AuthAttemptCollection = "AuthAttempt"

type AuthAttempt struct {
	mongodm.DocumentBase `json:",inline" bson:",inline"`

	IP       string `json:"ip" bson:"ip"`
	Username string `json:"username" bson:"username"`
}
