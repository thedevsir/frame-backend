package repository

import (
	"time"

	"github.com/satori/go.uuid"
	"github.com/zebresel-com/mongodm"
	"github.com/thedevsir/frame-backend/app/model"
	"github.com/thedevsir/frame-backend/config/database"
	"github.com/thedevsir/frame-backend/services/encrypt"
	"github.com/thedevsir/frame-backend/services/errors"
	"github.com/thedevsir/frame-backend/services/paginate"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func SessionCreate(IP, userID, userAgent string) (sid, key string, err error) {

	sessionModel := database.Connection.Model(model.SessionCollection)
	session := &model.Session{}
	sessionModel.New(session)

	// uuid use panic
	uuid := uuid.Must(uuid.NewV4(), nil).String() // nil
	hash, err := encrypt.Hash(uuid)
	if err != nil {
		return "", "", errors.ErrInternal
	}

	session.IP = IP
	session.Key = hash
	session.UserID = userID
	session.UserAgent = userAgent
	session.LastActivity = time.Now()
	session.ExpireAt = time.Now().Add(time.Hour * 24)

	err = session.Save()
	if err != nil {
		return "", "", errors.ErrInternal
	}

	return session.Id.Hex(), uuid, nil
}

func SessionFindByCredentials(Session, SID string) error {

	session, err := SessionFindByID(SID)
	if err != nil {
		return err
	}

	match := encrypt.CheckHash(Session, session.Key)
	if !match {
		return errors.ErrInvalidCredentials
	}

	return nil
}

func SessionFindByID(SID string) (*model.Session, error) {

	sessionModel := database.Connection.Model(model.SessionCollection)
	session := &model.Session{}

	err := sessionModel.FindId(bson.ObjectIdHex(SID)).Exec(session)
	_, ok := err.(*mongodm.NotFoundError)
	switch {
	case ok:
		return nil, errors.ErrSessionNotFound
	case err != nil:
		return nil, errors.ErrInternal
	default:
		return session, nil
	}
}

func SessionUpdateLastActivity(SID string) error {

	sessionModel := database.Connection.Model(model.SessionCollection)
	update := bson.M{
		"$set": bson.M{
			"lastActivity": time.Now(),
		},
	}

	err := sessionModel.UpdateId(bson.ObjectIdHex(SID), update)
	switch {
	case err == mgo.ErrNotFound:
		return errors.ErrInvalidCredentials
	case err != nil:
		return errors.ErrInternal
	default:
		return nil
	}
}

func GetUserSessions(userID string, page, limit int) (*paginate.Paginate, error) {

	sessionModel := database.Connection.Model(model.SessionCollection)
	sessions := []*model.Session{}
	result := sessionModel.Find(bson.M{"userId": userID}).
		Select(bson.M{"key": 0}).
		Sort("updatedAt").
		Skip((page - 1) * limit).
		Limit(limit)

	count, err := result.Count()
	if err != nil {
		return nil, errors.ErrInternal
	}

	err = result.Exec(&sessions)
	_, ok := err.(*mongodm.NotFoundError)
	switch {
	case ok || count == 0:
		return nil, errors.ErrSessionNotFound
	case err != nil:
		return nil, errors.ErrInternal
	}

	pagination := paginate.Generate(sessions, count, page, limit)
	return pagination, nil
}

func TerminateSession(ID string) error {

	sessionModel := database.Connection.Model(model.SessionCollection)
	err := sessionModel.RemoveId(bson.ObjectIdHex(ID))
	switch {
	case err == mgo.ErrNotFound:
		return errors.ErrSessionNotFound
	case err != nil:
		return errors.ErrInternal
	default:
		return nil
	}
}

func TerminateAllSessions(userID string) error {

	sessionModel := database.Connection.Model(model.SessionCollection)
	_, err := sessionModel.RemoveAll(bson.M{"userId": userID})
	switch {
	case err != nil:
		return errors.ErrInternal
	default:
		return nil
	}
}
