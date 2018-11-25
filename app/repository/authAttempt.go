package repository

import (
	"strings"
	"time"

	"github.com/thedevsir/frame-backend/app/model"
	"github.com/thedevsir/frame-backend/config"
	"github.com/thedevsir/frame-backend/config/database"
	"github.com/thedevsir/frame-backend/services/errors"
	"gopkg.in/mgo.v2/bson"
)

func CheckAbuse(ip, username string) error {

	anHourLater := time.Now().Add(-1 * time.Hour)
	authAttemptModel := database.Connection.Model(model.AuthAttemptCollection)

	numIP, err := authAttemptModel.Find(
		bson.M{
			"ip": ip,
			"createdAt": bson.M{
				"$gt": anHourLater,
			},
		},
	).Count()

	if err != nil {
		return errors.ErrInternal
	}

	numIPUsername, err := authAttemptModel.Find(
		bson.M{
			"ip":       ip,
			"username": username,
			"createdAt": bson.M{
				"$gt": anHourLater,
			},
		},
	).Count()

	if err != nil {
		return errors.ErrInternal
	}

	ipLimitReached := numIP >= config.AbuseIP
	ipUsernameLimitReached := numIPUsername >= config.AbuseIPUsername
	if ipLimitReached || ipUsernameLimitReached {
		return errors.ErrAttemptsReached
	}

	return nil
}

func SubmitAttempt(IP, username string) error {

	authAttemptModel := database.Connection.Model(model.AuthAttemptCollection)
	attempt := &model.AuthAttempt{}
	authAttemptModel.New(attempt)

	attempt.IP = IP
	attempt.Username = strings.ToLower(username)

	err := attempt.Save()
	if err != nil {
		return errors.ErrInternal
	}

	return nil
}
