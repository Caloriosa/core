package tools

import (
	"core/pkg/db"
	"core/types"
	"gopkg.in/mgo.v2/bson"
)

const COLLECTION_TOKENS = "tokens"

func GetUserFromToken(token string) (bson.ObjectId, error) {
	xtokens := []types.Token{}
	err := db.MONGO.Get(COLLECTION_TOKENS, &xtokens, bson.M{"token": token}, 0, 1)
	if err != nil {
		return "", err
	}
	if len(xtokens) == 0 {
		return "", nil
	}

	return *xtokens[0].User, nil
}
