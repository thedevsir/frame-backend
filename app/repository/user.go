package repository

import (
	"strings"

	"github.com/zebresel-com/mongodm"
	"github.com/thedevsir/frame-backend/app/model"
	"github.com/thedevsir/frame-backend/config/database"
	"github.com/thedevsir/frame-backend/services/encrypt"
	"github.com/thedevsir/frame-backend/services/errors"
	"github.com/thedevsir/frame-backend/services/paginate"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func FindUserByCredentials(username, password string) (*model.User, error) {

	userModel := database.Connection.Model(model.UserCollection)
	user := &model.User{}
	findStruct := bson.M{"isActive": true}

	if strings.Index(username, "@") > -1 {
		findStruct["email"] = strings.ToLower(username)
	} else {
		findStruct["username"] = strings.ToLower(username)
	}

	err := userModel.FindOne(findStruct).Exec(user)
	_, ok := err.(*mongodm.NotFoundError)
	switch {
	case ok:
		return nil, errors.ErrUserNotFound
	case err != nil:
		return nil, errors.ErrInternal
	default:
		match := encrypt.CheckHash(password, user.Password)
		if !match {
			return nil, errors.ErrInvalidCredentials
		}
	}

	return user, nil
}

func ChangePassword(userID, password string, admin bool) error {

	userModel := database.Connection.Model(model.UserCollection)
	hash, err := encrypt.Hash(password)
	if err != nil {
		return errors.ErrInternal
	}

	findStruct := bson.M{"_id": bson.ObjectIdHex(userID)}
	if !admin {
		findStruct["isActive"] = true
	}
	update := bson.M{
		"$set": bson.M{
			"password": hash,
		},
	}

	err = userModel.Update(findStruct, update)
	switch {
	case err == mgo.ErrNotFound:
		return errors.ErrUserNotFound
	case err != nil:
		return errors.ErrInternal
	default:
		return nil
	}
}

func UserActivation(userID string) error {

	userModel := database.Connection.Model(model.UserCollection)
	update := bson.M{
		"$set": bson.M{
			"isEmailVerified": true,
		},
	}
	err := userModel.Update(bson.M{"_id": bson.ObjectIdHex(userID), "isActive": true}, update)
	switch {
	case err == mgo.ErrNotFound:
		return errors.ErrUserNotFound
	case err != nil:
		return errors.ErrInternal
	default:
		return nil
	}
}

func CheckUsername(username string) (*model.User, error) {

	userModel := database.Connection.Model(model.UserCollection)
	user := &model.User{}
	err := userModel.FindOne(bson.M{"username": strings.ToLower(username)}).Exec(user)
	_, ok := err.(*mongodm.NotFoundError)
	switch {
	case ok:
		return nil, nil
	case err != nil:
		return nil, errors.ErrInternal
	default:
		return user, errors.ErrUsernameExists
	}
}

func CheckEmail(email string) (*model.User, error) {

	userModel := database.Connection.Model(model.UserCollection)
	user := &model.User{}
	err := userModel.FindOne(bson.M{"email": strings.ToLower(email)}).Exec(user)
	_, ok := err.(*mongodm.NotFoundError)
	switch {
	case ok:
		return nil, nil
	case err != nil:
		return nil, errors.ErrInternal
	default:
		return user, errors.ErrEmailExists
	}
}

func CreateUser(username, password, email string) (*model.User, error) {

	userModel := database.Connection.Model(model.UserCollection)
	user := &model.User{}
	userModel.New(user)
	hash, err := encrypt.Hash(password)
	if err != nil {
		return nil, errors.ErrInternal
	}

	user.Username = strings.ToLower(username)
	user.Password = hash
	user.Email = strings.ToLower(email)
	user.IsEmailVerified = false
	user.IsActive = true

	err = user.Save()
	if err != nil {
		return nil, errors.ErrInternal
	}

	return user, nil
}

func ChangeUsername(userID, username string, admin bool) error {

	userModel := database.Connection.Model(model.UserCollection)
	findStruct := bson.M{"_id": bson.ObjectIdHex(userID)}
	if !admin {
		findStruct["isActive"] = true
	}
	update := bson.M{
		"$set": bson.M{
			"username": username,
		},
	}

	err := userModel.Update(findStruct, update)
	switch {
	case err == mgo.ErrNotFound:
		return errors.ErrUserNotFound
	case err != nil:
		return errors.ErrInternal
	default:
		return nil
	}
}

func GetAccountInfo(userID string) (*model.User, error) {

	userModel := database.Connection.Model(model.UserCollection)
	user := &model.User{}
	err := userModel.FindOne(bson.M{"_id": bson.ObjectIdHex(userID)}).Select(bson.M{"password": 0}).Exec(user)
	_, ok := err.(*mongodm.NotFoundError)
	switch {
	case ok:
		return nil, errors.ErrUserNotFound
	case err != nil:
		return nil, errors.ErrInternal
	default:
		return user, nil
	}
}

func ChangeEmail(userID, email string, admin bool) error {

	userModel := database.Connection.Model(model.UserCollection)
	findStruct := bson.M{"_id": bson.ObjectIdHex(userID)}
	if !admin {
		findStruct["isActive"] = true
	}
	update := bson.M{
		"$set": bson.M{
			"IsEmailVerified": false,
			"email":           email,
		},
	}

	err := userModel.Update(findStruct, update)
	switch {
	case err == mgo.ErrNotFound:
		return errors.ErrUserNotFound
	case err != nil:
		return errors.ErrInternal
	default:
		return nil
	}
}

func GetUsers(page, limit int) (*paginate.Paginate, error) {

	userModel := database.Connection.Model(model.UserCollection)
	users := []*model.User{}
	result := userModel.Find(nil).
		Sort("createdAt").
		Skip((page - 1) * limit).
		Limit(limit)

	count, err := result.Count()
	if err != nil {
		return nil, errors.ErrInternal
	}

	err = result.Exec(&users)
	_, ok := err.(*mongodm.NotFoundError)
	switch {
	case ok || count == 0:
		return nil, errors.ErrUserNotFound
	case err != nil:
		return nil, err
	}

	pagination := paginate.Generate(users, count, page, limit)
	return pagination, nil
}

func GetUserByID(userID string) (*model.User, error) {

	userModel := database.Connection.Model(model.UserCollection)
	user := &model.User{}
	result := userModel.FindOne(bson.M{"_id": bson.ObjectIdHex(userID), "isActive": true}).
		Select(bson.M{"username": 1, "email": 1})

	err := result.Exec(user)
	_, ok := err.(*mongodm.NotFoundError)
	switch {
	case ok:
		return nil, errors.ErrUserNotFound
	case err != nil:
		return nil, errors.ErrInternal
	default:
		return user, nil
	}
}

func GetUserByIDFromAdmin(userID string) (*model.User, error) {

	userModel := database.Connection.Model(model.UserCollection)
	user := &model.User{}
	result := userModel.FindId(bson.ObjectIdHex(userID))

	err := result.Exec(user)
	_, ok := err.(*mongodm.NotFoundError)
	switch {
	case ok:
		return nil, errors.ErrUserNotFound
	case err != nil:
		return nil, errors.ErrInternal
	default:
		return user, nil
	}
}

func ChangeUserStatus(userID string, status bool) error {

	userModel := database.Connection.Model(model.UserCollection)
	update := bson.M{
		"$set": bson.M{
			"isActive": status,
		},
	}

	err := userModel.UpdateId(bson.ObjectIdHex(userID), update)
	switch {
	case err == mgo.ErrNotFound:
		return errors.ErrUserNotFound
	case err != nil:
		return errors.ErrInternal
	default:
		return TerminateAllSessions(userID)
	}
}
