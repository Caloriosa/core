package userlib

import (
	"core/pkg/db"
	"core/pkg/error"
	"core/pkg/validation"
	"core/types"
	"core/types/httptypes"
	"gopkg.in/mgo.v2/bson"
)

const COLLECTION_USERS = "users"

func CreateUser(newUser *types.User) *errors.CalError {

	if err := validation.ValidateNewUser(newUser); err != nil {
		return &errors.CalError{Status: &httptypes.DATA_INCOMPLETE}
	}

	if err := db.MONGO.Save(COLLECTION_USERS, newUser); err != nil {
		return &errors.CalError{Status: &httptypes.SERVICE_UNAVAILABLE}
	}

	return nil
}

func FindUserById(id string, user *types.User) *errors.CalError {
	err := db.MONGO.FindById(COLLECTION_USERS, &user, id)
	if err != nil {
		return &errors.CalError{Status: &httptypes.NOT_FOUND}
	}

	return nil
}

func SaveUser(user *types.User) *errors.CalError {
	err := db.MONGO.Save(COLLECTION_USERS, user)
	if err != nil {
		return &errors.CalError{Status: &httptypes.UNAVAILABLE}
	}

	return nil
}

func DeleteUser(user *types.User) *errors.CalError {
	if err := db.MONGO.Connection.Collection(COLLECTION_USERS).DeleteDocument(user); err != nil {
		return &errors.CalError{Status: &httptypes.REMOVE_FAILED}
	}

	return nil
}

func ActivateUserByActToken(token string) (*types.User, *errors.CalError) {
	user := types.User{}
	if err := db.MONGO.Get(COLLECTION_USERS, &user, bson.M{"activationkey": token}, 0, 1); err != nil {
		return nil, &errors.CalError{Status: &httptypes.INVALID_TOKEN}
	}

	user.Activated = true
	user.ActivationKey = ""
	db.MONGO.Save(COLLECTION_USERS, &user)

	return &user, nil
}
