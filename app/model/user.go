package model

import (
	"github.com/zebresel-com/mongodm"
)

const UserCollection = "User"

type User struct {
	mongodm.DocumentBase `json:",inline" bson:",inline"`

	Username        string `json:"username" bson:"username"`
	Password        string `json:"password" bson:"password"`
	Email           string `json:"email" bson:"email"`
	IsEmailVerified bool   `json:"isEmailVerified" bson:"isEmailVerified"`
	IsActive        bool   `json:"isActive" bson:"isActive"`
}
