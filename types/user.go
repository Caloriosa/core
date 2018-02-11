package types

import (
	"core/pkg/config"
	"core/pkg/db"
	calerror "core/pkg/error"
	"core/pkg/tools"
	"core/types/httptypes"
	"errors"
	"github.com/go-bongo/bongo"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

type UserDB struct {
	bongo.DocumentBase `bson:",inline"`
	User               User `json:",inline" bson:",inline"`
	AllowedToSave      bool `json:"-" bson:"-"`
}

const COLLECTION_USERS = "users"
const SALT_LENGTH = 32

func (u *UserDB) Save() *calerror.CalError {
	if !u.AllowedToSave {
		glog.Warning("Tried saving a copy! NOPE")
		return nil
	}

	err := db.MONGO.Save(COLLECTION_USERS, u)
	if mgo.IsDup(err) {
		return &calerror.CalError{Status: &httptypes.DUPLICATED}
	} else if err != nil {
		glog.Info("Got error saving: ", err.Error())
		return &calerror.CalError{Status: &httptypes.DATASOURCE_ERROR}
	}

	return nil
}

func (u *UserDB) Delete() *calerror.CalError {
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

func (u *UserDB) Activate() *calerror.CalError {
	u.User.Activated = true
	u.User.ActivationKey = nil
	u.User.ActivationExpiry = nil

	return u.Save()
}

func GetUserById(id string, user *UserDB) *calerror.CalError {
	err := db.MONGO.FindById(COLLECTION_USERS, &user, id)
	if err != nil {
		return &calerror.CalError{Status: &httptypes.NOT_FOUND}
	}

	return nil
}

func (u *UserDB) PrepareToSend(strict bool) *User {
	newUser := new(User)
	*newUser = u.User
	newUser.UID = u.Id.Hex()
	newUser.CreatedAt = u.Created
	newUser.ModifiedAt = u.Modified
	newUser.Password = "" // sanitize
	if strict {
		newUser.Email = ""
	}
	return newUser
}

func GetUserByLogin(login string) (*calerror.CalError, *UserDB) {
	users := []UserDB{}
	if err := db.MONGO.Get(COLLECTION_USERS, &users, bson.M{"login": login}, 0, 1); err != nil || len(users) == 0 {
		return &calerror.CalError{Status: &httptypes.NOT_FOUND}, nil
	}
	return nil, &users[0]
}

func ActivateUserByToken(token string) (*UserDB, *calerror.CalError) {
	users := []UserDB{}
	user := UserDB{}
	if err := db.MONGO.Get(COLLECTION_USERS, &users, bson.M{"activationkey": token}, 0, 1); err != nil {
		return nil, &calerror.CalError{Status: &httptypes.INVALID_TOKEN}
	}

	if len(users) == 0 {
		return nil, &calerror.CalError{Status: &httptypes.INVALID_TOKEN}
	}

	user = users[0]

	if user.User.ActivationExpiry == nil || user.User.Activated {
		glog.Info("Tried reactivating already activated user")
		return &user, nil
	}

	if user.User.ActivationExpiry.Before(time.Now()) {
		glog.Info("Tried activating expired user: ", user.User.ActivationExpiry)
		return nil, &calerror.CalError{Status: &httptypes.TOKEN_EXPIRED}
	}

	err := user.Activate()

	return &user, err
}

func GetUserFromRequest(r *http.Request) (*UserDB, *calerror.CalError) {
	token := GetTokenFromRequest(r)
	if token == nil {
		return nil, &calerror.CalError{Status: &httptypes.UNAUTHORIZED}
	}
	if token.Token.User == nil || token.Token.Type != TokenUser {
		return nil, &calerror.CalError{Status: &httptypes.INVALID_TOKEN}
	}

	user := UserDB{}
	if err := GetUserById(token.Token.User.Hex(), &user); err != nil {
		return nil, &calerror.CalError{Status: &httptypes.INVALID_TOKEN}
	}

	return &user, nil
}

func CreateUser(newUser *User) (*UserDB, *calerror.CalError) {

	if err := newUser.Validate(); err != nil {
		glog.Info("Error validating a new user: ", err, " user: ", newUser)
		return nil, &calerror.CalError{Status: &httptypes.DATA_PARSE_ERROR}
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

	db := UserDB{}
	db.User = *newUser
	db.AllowedToSave = true

	return &db, db.Save()
}
func GetUserFromTokenString(token string) (bson.ObjectId, *calerror.CalError) {
	xtokens := []Token{}
	err := db.MONGO.Get(COLLECTION_TOKENS, &xtokens, bson.M{"token": token}, 0, 1)
	if err != nil {
		return "", &calerror.CalError{Status: &httptypes.SERVER_ERROR}
	}
	if len(xtokens) == 0 {
		return "", &calerror.CalError{Status: &httptypes.INVALID_TOKEN}
	}

	return *xtokens[0].User, nil
}
