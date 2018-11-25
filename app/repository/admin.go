package repository

import (
	"fmt"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/zebresel-com/mongodm"
	"github.com/thedevsir/frame-backend/app/model"
	"github.com/thedevsir/frame-backend/config/database"
	"github.com/thedevsir/frame-backend/services/encrypt"
	"github.com/thedevsir/frame-backend/services/errors"
	"github.com/thedevsir/frame-backend/services/paginate"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func CreateAdmin(username, password string) (adminID string, err error) {

	if _, err := CheckAdminUsername(username); err != nil {
		return "", err
	}

	adminModel := database.Connection.Model(model.AdminCollection)
	admin := &model.Admin{}
	adminModel.New(admin)
	hash, err := encrypt.Hash(password)
	if err != nil {
		return "", errors.ErrInternal
	}

	admin.Username = username
	admin.Password = hash
	admin.IsActive = true

	err = admin.Save()
	if err != nil {
		return "", errors.ErrInternal
	}

	return admin.Id.Hex(), nil
}

func CheckAdminUsername(username string) (*model.Admin, error) {

	adminModel := database.Connection.Model(model.AdminCollection)
	admin := &model.Admin{}
	err := adminModel.FindOne(bson.M{"username": strings.ToLower(username)}).Exec(admin)
	_, ok := err.(*mongodm.NotFoundError)
	switch {
	case ok:
		return nil, nil
	case err != nil:
		return nil, errors.ErrInternal
	default:
		return admin, errors.ErrUsernameExists
	}
}

func AdminChangeUsername(adminID, username string) error {

	if _, err := CheckAdminUsername(username); err != nil {
		return err
	}

	adminModel := database.Connection.Model(model.AdminCollection)
	update := bson.M{
		"$set": bson.M{
			"username": username,
		},
	}

	err := adminModel.UpdateId(bson.ObjectIdHex(adminID), update)
	switch {
	case err == mgo.ErrNotFound:
		return errors.ErrAdminNotFound
	case err != nil:
		return errors.ErrInternal
	default:
		return nil
	}
}

func AdminChangePassword(adminID, password string) error {

	hash, err := encrypt.Hash(password)
	if err != nil {
		return errors.ErrInternal
	}

	adminModel := database.Connection.Model(model.AdminCollection)
	update := bson.M{
		"$set": bson.M{
			"password": hash,
		},
	}

	err = adminModel.UpdateId(bson.ObjectIdHex(adminID), update)
	switch {
	case err == mgo.ErrNotFound:
		return errors.ErrAdminNotFound
	case err != nil:
		return errors.ErrInternal
	default:
		return nil
	}
}

func FindAdminByCredentials(username, password string) (*model.Admin, error) {

	adminModel := database.Connection.Model(model.AdminCollection)
	admin := &model.Admin{}
	err := adminModel.FindOne(bson.M{"username": strings.ToLower(username), "isActive": true}).Exec(admin)
	_, ok := err.(*mongodm.NotFoundError)
	switch {
	case ok:
		return nil, errors.ErrAdminNotFound
	case err != nil:
		return nil, errors.ErrInternal
	default:
		match := encrypt.CheckHash(password, admin.Password)
		if !match {
			return nil, errors.ErrInvalidCredentials
		}
	}
	return admin, nil
}

func SetAdminSession(adminID string) (key string, err error) {

	uuid := uuid.Must(uuid.NewV4(), nil).String() // nil
	hash, err := encrypt.Hash(uuid)
	if err != nil {
		return "", errors.ErrInternal
	}

	adminModel := database.Connection.Model(model.AdminCollection)
	update := bson.M{
		"$set": bson.M{
			"session":   hash,
			"lastLogin": time.Now(),
		},
	}

	err = adminModel.Update(bson.M{"_id": bson.ObjectIdHex(adminID), "isActive": true}, update)
	switch {
	case err == mgo.ErrNotFound:
		return "", errors.ErrAdminNotFound
	case err != nil:
		return "", errors.ErrInternal
	default:
		return uuid, nil
	}
}

func TerminateAdminSession(adminID string) error {

	adminModel := database.Connection.Model(model.AdminCollection)
	update := bson.M{
		"$set": bson.M{
			"session": "",
		},
	}

	err := adminModel.Update(bson.M{"_id": bson.ObjectIdHex(adminID), "isActive": true}, update)
	switch {
	case err == mgo.ErrNotFound:
		return errors.ErrAdminNotFound
	case err != nil:
		return errors.ErrInternal
	default:
		return nil
	}
}

func CheckAdminSession(adminID, session string) error {

	adminModel := database.Connection.Model(model.AdminCollection)
	admin := &model.Admin{}

	err := adminModel.FindOne(bson.M{"_id": bson.ObjectIdHex(adminID), "isActive": true}).Exec(admin)
	_, ok := err.(*mongodm.NotFoundError)
	switch {
	case ok:
		return errors.ErrAdminNotFound
	case err != nil:
		return errors.ErrInternal
	default:
		match := encrypt.CheckHash(session, admin.Session)
		if !match {
			return errors.ErrInvalidCredentials
		}
	}
	return nil
}

func AdminSessionUpdateLastActivity(adminID string) error {

	adminModel := database.Connection.Model(model.AdminCollection)
	update := bson.M{
		"$set": bson.M{
			"lastActivity": time.Now(),
		},
	}

	fmt.Print(11112, adminID)
	err := adminModel.Update(bson.M{"_id": bson.ObjectIdHex(adminID), "isActive": true}, update)
	switch {
	case err == mgo.ErrNotFound:
		return errors.ErrAdminNotFound
	case err != nil:
		return errors.ErrInternal
	default:
		return nil
	}
}

func GetAdmins(page, limit int) (*paginate.Paginate, error) {

	adminModel := database.Connection.Model(model.AdminCollection)
	admins := []*model.Admin{}
	result := adminModel.Find(nil).
		Sort("createdAt").
		Skip((page - 1) * limit).
		Limit(limit)

	count, err := result.Count()
	if err != nil {
		return nil, errors.ErrInternal
	}

	err = result.Exec(&admins)
	_, ok := err.(*mongodm.NotFoundError)
	switch {
	case ok || count == 0:
		return nil, errors.ErrAdminNotFound
	case err != nil:
		return nil, errors.ErrInternal
	}

	pagination := paginate.Generate(admins, count, page, limit)
	return pagination, nil
}

func GetAdminByID(adminID string) (*model.Admin, error) {

	adminModel := database.Connection.Model(model.AdminCollection)
	admin := &model.Admin{}
	result := adminModel.FindId(bson.ObjectIdHex(adminID))

	err := result.Exec(admin)
	_, ok := err.(*mongodm.NotFoundError)
	switch {
	case ok:
		return nil, errors.ErrAdminNotFound
	case err != nil:
		return nil, errors.ErrInternal
	default:
		return admin, nil
	}
}

func ChangeAdminStatus(adminID string, status bool) error {

	adminModel := database.Connection.Model(model.AdminCollection)
	newStruct := bson.M{"isActive": status}
	if !status {
		newStruct["session"] = ""
	}
	update := bson.M{
		"$set": newStruct,
	}

	err := adminModel.UpdateId(bson.ObjectIdHex(adminID), update)
	switch {
	case err == mgo.ErrNotFound:
		return errors.ErrAdminNotFound
	case err != nil:
		return errors.ErrInternal
	default:
		return nil
	}
}
