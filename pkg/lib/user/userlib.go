package userlib

import (
	"core/pkg/config"
	"core/pkg/db"
	"core/pkg/error"
	"core/pkg/tools"
	"core/pkg/validation"
	"core/types"
	"core/types/httptypes"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

const COLLECTION_USERS = "users"
const SALT_LENGTH = 32

func CreateUser(newUser *types.User) *errors.CalError {

	if err := validation.ValidateNewUser(newUser); err != nil {
		glog.Info("Error validating a new user: ", err, " user: ", newUser)
		return &errors.CalError{Status: &httptypes.DATA_INCOMPLETE}
	}

	newUser.Role = types.RoleUser
	actkey := tools.RandStringRunes(types.LENGTH_TOKEN)
	actexp := time.Now().UTC().Add(time.Duration(config.LoadedConfig.Users.ActivationExpiry) * time.Hour)
	newUser.ActivationKey = &actkey
	newUser.ActivationExpiry = &actexp
	// re-enter password but as argon2
	tmppwd := newUser.Password
	salt := tools.RandStringRunes(SALT_LENGTH)
	newUser.Password = tools.EncodeUserPassword(tmppwd, salt)
	newUser.Salt = salt

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
	users := []types.User{}
	user := types.User{}
	if err := db.MONGO.Get(COLLECTION_USERS, &users, bson.M{"activationkey": token}, 0, 1); err != nil {
		return nil, &errors.CalError{Status: &httptypes.INVALID_TOKEN}
	}

	if len(users) == 0 {
		return nil, &errors.CalError{Status: &httptypes.INVALID_TOKEN}
	}

	user = users[0]

	if user.ActivationExpiry == nil || user.Activated {
		glog.Info("Tried reactivating already activated user")
		return &user, nil
	}

	if user.ActivationExpiry.Before(time.Now()) {
		glog.Info("Tried activating expired user: ", user.ActivationExpiry)
		return nil, &errors.CalError{Status: &httptypes.TOKEN_EXPIRED}
	}

	user.Activated = true
	user.ActivationKey = nil
	user.ActivationExpiry = nil
	glog.Info("Activated user: ", &user)
	if err := db.MONGO.Save(COLLECTION_USERS, &user); err != nil {
		glog.Info("Error saving activated user: ", err)
		return nil, &errors.CalError{Status: &httptypes.DATASOURCE_ERROR}
	}

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
