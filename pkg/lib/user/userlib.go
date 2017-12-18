package userlib

import (
	"core/pkg/db"
	"core/pkg/error"
	"core/pkg/validation"
	"core/types"
	"core/types/httptypes"
	"gopkg.in/mgo.v2/bson"
	"github.com/golang/glog"
	"core/pkg/tools"
	"net/http"
	"gopkg.in/mgo.v2"
)

const COLLECTION_USERS = "users"

func CreateUser(newUser *types.User) *errors.CalError {

	if err := validation.ValidateNewUser(newUser); err != nil {
		return &errors.CalError{Status: &httptypes.DATA_INCOMPLETE}
	}

	newUser.Role = types.RoleUser
	actkey := tools.RandStringRunes(types.LENGTH_TOKEN)
	newUser.ActivationKey = &actkey

	if err := db.MONGO.Save(COLLECTION_USERS, newUser); err != nil {
		if mgo.IsDup(err) {
			return &errors.CalError{Status: &httptypes.DUPLICATED}
		} else {
			glog.Info("Got error saving: ", err.Error())
			return &errors.CalError{Status: &httptypes.DATASOURCE_ERROR}
		}
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
		glog.Warning("Error saving user: ", err)
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
	user.ActivationKey = nil
	db.MONGO.Save(COLLECTION_USERS, &user)

	return &user, nil
}

func GetUserFromRequest(r *http.Request) (*types.User, *errors.CalError) {
	token := tools.GetToken(r)
	if token == nil {
		return nil, &errors.CalError{Status: &httptypes.UNAUTHORIZED}
	}

	uid, err := tools.GetUserFromTokenString(token.Token)
	if err != nil {
		return nil, &errors.CalError{Status: &httptypes.DATASOURCE_ERROR}
	}
	if uid == "" {
		return nil, &errors.CalError{Status: &httptypes.INVALID_TOKEN}
	}

	user := types.User{}
	if err := FindUserById(uid.Hex(), &user); err != nil {
		return nil, &errors.CalError{Status: &httptypes.INVALID_TOKEN}
	}

	return &user, nil
}
