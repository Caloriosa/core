package types

import (
	"core/pkg/db"
	"core/pkg/error"
	"core/types/httptypes"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strings"
	"time"
)

const COLLECTION_TOKENS = "tokens"

func GetTokensForUser(user *User, tokens []*Token) *errors.CalError {
	if err := db.MONGO.Get(COLLECTION_TOKENS, tokens, bson.M{"type": TokenUser}, 0, 999); err != nil {
		return &errors.CalError{Status: &httptypes.DATASOURCE_ERROR}
	}

	return nil
}

func GetToken(tokenid string, token *Token) *errors.CalError {
	if err := db.MONGO.Get(COLLECTION_TOKENS, token, bson.M{"token": tokenid}, 0, 1); err != nil {
		return &errors.CalError{Status: &httptypes.DATASOURCE_ERROR}
	}

	return nil
}

func (t *Token) Delete() *errors.CalError {
	if err := db.MONGO.Connection.Collection(COLLECTION_TOKENS).DeleteDocument(t); err != nil {
		return &errors.CalError{Status: &httptypes.DATASOURCE_ERROR}
	}

	return nil
}

func GetTokenFromRequest(req *http.Request) *Token {
	mytoken := []Token{}
	glog.Info("Headers: ", req.Header)
	if tokens, ok := req.Header["Authorization"]; ok {
		glog.Info("Found tokens: ", tokens)
		if len(tokens) > 0 {
			glog.Info("X token: ", tokens[0], "token: ", mytoken)
			xtoken := strings.TrimPrefix(tokens[0], "Bearer ")
			if err := db.MONGO.Get(COLLECTION_TOKENS, &mytoken, bson.M{"token": xtoken}, 0, 1); err == nil {
				glog.Info("Extracted token: ", mytoken)
				if len(mytoken) == 0 {
					return nil
				}

				glog.Info("ExpireAt: ", mytoken[0].ExpireAt.UTC(), " now: ", time.Now().UTC(), "is ok? ", mytoken[0].ExpireAt.Before(time.Now().UTC()))

				if mytoken[0].ExpireAt.UTC().After(time.Now().UTC()) {
					return &mytoken[0]
				} else {
					db.MONGO.Connection.Collection(COLLECTION_TOKENS).DeleteDocument(&mytoken[0])
					glog.Info("Deleting expired token.")
				}
			}
		}
	}

	return nil
}
