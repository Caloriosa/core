package tools

import (
	"core/pkg/db"
	"core/pkg/error"
	"core/pkg/lib/user"
	"core/types"
	"core/types/httptypes"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

const COLLECTION_TOKENS = "tokens"

func GetUserFromTokenString(token string) (bson.ObjectId, *errors.CalError) {
	xtokens := []types.Token{}
	err := db.MONGO.Get(COLLECTION_TOKENS, &xtokens, bson.M{"token": token}, 0, 1)
	if err != nil {
		return "", &errors.CalError{Status: &httptypes.UNKNOWN}
	}
	if len(xtokens) == 0 {
		return "", &errors.CalError{Status: &httptypes.INVALID_TOKEN}
	}

	return *xtokens[0].User, nil
}

func GetUserFromRequest(r *http.Request) (*types.User, *errors.CalError) {
	token := GetToken(r)
	if token == nil {
		return nil, &errors.CalError{Status: &httptypes.DATA_INCOMPLETE}
	}

	uid, err := GetUserFromTokenString(token.Token)
	if err != nil {
		return nil, &errors.CalError{Status: &httptypes.DATASOURCE_ERROR}
	}
	if uid == "" {
		return nil, &errors.CalError{Status: &httptypes.INVALID_TOKEN}
	}

	user := types.User{}
	if err := userlib.FindUserById(uid.Hex(), &user); err != nil {
		return nil, &errors.CalError{Status: &httptypes.INVALID_TOKEN}
	}

	return &user, nil
}
