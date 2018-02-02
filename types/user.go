package types

import (
	"core/pkg/config"
	"core/pkg/db"
	calerror "core/pkg/error"
	"core/pkg/tools"
	"core/types/httptypes"
	"errors"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

const COLLECTION_USERS = "users"
const SALT_LENGTH = 32

func (u *User) Save() *calerror.CalError {
	err := db.MONGO.Save(COLLECTION_USERS, u)
	if mgo.IsDup(err) {
		return &calerror.CalError{Status: &httptypes.DUPLICATED}
	} else {
		glog.Info("Got error saving: ", err.Error())
		return &calerror.CalError{Status: &httptypes.DATASOURCE_ERROR}
	}

	return nil
}

func (u *User) Delete() *calerror.CalError {
	if err := db.MONGO.Connection.Collection(COLLECTION_USERS).DeleteDocument(u); err != nil {
		return &calerror.CalError{Status: &httptypes.REMOVE_FAILED}
	}

	return nil
}

func (u *User) Validate() error {
	if u.Email == "" {
		return errors.New("Missing e-mail")
	}
	if u.Password == "" {
		return errors.New("Missing password")
	}
	if u.Login == "" {
		return errors.New("Missing login")
	}
	return nil // TODO validate email
}

func (u *User) MergeIn(with *User) {
	/*if with.Login != "" {
		user.Login = with.Login
	}*/

	if with.Password != "" {
		u.Password = tools.EncodeUserPassword(u.Salt, with.Password)
	}

	if with.Email != "" {
		u.Email = with.Email
	}

	if with.Name != "" {
		u.Name = with.Name
	}
}

func (u *User) Activate() *calerror.CalError {
	u.Activated = true
	u.ActivationKey = nil
	u.ActivationExpiry = nil

	return u.Save()
}

func GetUserById(id string, user *User) *calerror.CalError {
	err := db.MONGO.FindById(COLLECTION_USERS, &user, id)
	if err != nil {
		return &calerror.CalError{Status: &httptypes.NOT_FOUND}
	}

	return nil
}

func GetUserByLogin(login string) (*calerror.CalError, *User) {
	users := []User{}
	if err := db.MONGO.Get(COLLECTION_USERS, &users, bson.M{"login": login}, 0, 1); err != nil || len(users) == 0{
		return &calerror.CalError{Status: &httptypes.NOT_FOUND}, nil
	}
	return nil, &users[0]
}

func ActivateUserByToken(token string) (*User, *calerror.CalError) {
	users := []User{}
	user := User{}
	if err := db.MONGO.Get(COLLECTION_USERS, &users, bson.M{"activationkey": token}, 0, 1); err != nil {
		return nil, &calerror.CalError{Status: &httptypes.INVALID_TOKEN}
	}

	if len(users) == 0 {
		return nil, &calerror.CalError{Status: &httptypes.INVALID_TOKEN}
	}

	user = users[0]

	if user.ActivationExpiry == nil || user.Activated {
		glog.Info("Tried reactivating already activated user")
		return &user, nil
	}

	if user.ActivationExpiry.Before(time.Now()) {
		glog.Info("Tried activating expired user: ", user.ActivationExpiry)
		return nil, &calerror.CalError{Status: &httptypes.TOKEN_EXPIRED}
	}

	return &user, user.Activate()
}

func GetUserFromRequest(r *http.Request) (*User, *calerror.CalError) {
	token := GetTokenFromRequest(r)
	if token == nil {
		return nil, &calerror.CalError{Status: &httptypes.UNAUTHORIZED}
	}
	if token.User == nil || token.Type != TokenUser {
		return nil, &calerror.CalError{Status: &httptypes.INVALID_TOKEN}
	}

	user := User{}
	if err := GetUserById(token.User.Hex(), &user); err != nil {
		return nil, &calerror.CalError{Status: &httptypes.INVALID_TOKEN}
	}

	return &user, nil
}

func CreateUser(newUser *User) *calerror.CalError {

	if err := newUser.Validate(); err != nil {
		glog.Info("Error validating a new user: ", err, " user: ", newUser)
		return &calerror.CalError{Status: &httptypes.DATA_INCOMPLETE}
	}

	newUser.Role = RoleUser
	actkey := tools.RandStringRunes(LENGTH_TOKEN)
	actexp := time.Now().UTC().Add(time.Duration(config.LoadedConfig.Users.ActivationExpiry) * time.Hour)
	newUser.ActivationKey = &actkey
	newUser.ActivationExpiry = &actexp
	// re-enter password but as argon2
	tmppwd := newUser.Password
	salt := tools.RandStringRunes(SALT_LENGTH)
	newUser.Password = tools.EncodeUserPassword(tmppwd, salt)
	newUser.Salt = salt

	return newUser.Save()
}
func GetUserFromTokenString(token string) (bson.ObjectId, *calerror.CalError) {
	xtokens := []Token{}
	err := db.MONGO.Get(COLLECTION_TOKENS, &xtokens, bson.M{"token": token}, 0, 1)
	if err != nil {
		return "", &calerror.CalError{Status: &httptypes.UNKNOWN}
	}
	if len(xtokens) == 0 {
		return "", &calerror.CalError{Status: &httptypes.INVALID_TOKEN}
	}

	return *xtokens[0].User, nil
}
