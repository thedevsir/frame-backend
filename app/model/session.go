package model

import (
	"time"

	"github.com/zebresel-com/mongodm"
)

const SessionCollection = "Session"

type Session struct {
	mongodm.DocumentBase `json:",inline" bson:",inline"`

	IP           string    `json:"ip" bson:"ip"`
	Key          string    `json:"key" bson:"key"`
	UserID       string    `json:"userId" bson:"userId"`
	UserAgent    string    `json:"userAgent" bson:"userAgent"`
	LastActivity time.Time `json:"lastActivity" bson:"lastActivity"`
	ExpireAt     time.Time `json:"expireAt" bson:"expireAt"`
}
