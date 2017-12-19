package tools

import (
	"core/pkg/db"
	"core/pkg/error"
	"core/types"
	"core/types/httptypes"
	"gopkg.in/mgo.v2/bson"
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
