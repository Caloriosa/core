package types

import (
	"core/pkg/db"
	"core/pkg/error"
	"core/types/httptypes"
	"github.com/go-bongo/bongo"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strings"
	"time"
)

const COLLECTION_TOKENS = "tokens"

type TokenDB struct {
	bongo.DocumentBase `bson:",inline"`
	Token              Token `json:",inline" bson:",inline"`
}

func GetTokensForUser(user *UserDB, tokens []*TokenDB) *errors.CalError {
	if err := db.MONGO.Get(COLLECTION_TOKENS, tokens, bson.M{"type": TokenUser, "user": user.Id.Hex()}, 0, 999); err != nil {
		return &errors.CalError{Status: &httptypes.DATASOURCE_ERROR}
	}

	return nil
}

func GetToken(tokenid string, token *TokenDB) *errors.CalError {
	if err := db.MONGO.Get(COLLECTION_TOKENS, token, bson.M{"token": tokenid}, 0, 1); err != nil {
		return &errors.CalError{Status: &httptypes.DATASOURCE_ERROR}
	}

	return nil
}

func (t *TokenDB) Delete() *errors.CalError {
	if err := db.MONGO.Connection.Collection(COLLECTION_TOKENS).DeleteDocument(t); err != nil {
		return &errors.CalError{Status: &httptypes.DATASOURCE_ERROR}
	}

	return nil
}

func (t *TokenDB) Save() *errors.CalError {
	if err := db.MONGO.Connection.Collection(COLLECTION_TOKENS).Save(t); err != nil {
		return &errors.CalError{Status: &httptypes.DATASOURCE_ERROR}
	}

	return nil
}

func (t *TokenDB) PrepareToSend() *Token {
	token := new(Token)
	*token = t.Token
	token.UID = t.Id.Hex()
	token.CreatedAt = t.Created
	token.ModifiedAt = t.Modified
	return token
}

func CreateNewToken(t *Token) (*TokenDB, *errors.CalError) {
	db := TokenDB{Token: *t}
	if err := db.Save(); err != nil {
		return nil, err
	}
	return &db, nil
}

func GetTokenFromRequest(req *http.Request) *TokenDB {
	mytoken := []TokenDB{}
	glog.Info("Headers: ", req.Header)
	if tokens, ok := req.Header["Authorization"]; ok {
		glog.Info("Found tokens: ", tokens)
		if len(tokens) > 0 {
			glog.Info("X token: ", tokens[0], "token: ", mytoken)
			xtoken := strings.TrimSpace(strings.TrimPrefix(tokens[0], "Bearer "))
			if err := db.MONGO.Get(COLLECTION_TOKENS, &mytoken, bson.M{"token": xtoken}, 0, 1); err == nil {
				glog.Info("Extracted token: ", mytoken)
				if len(mytoken) == 0 {
					return nil
				}

				glog.Info("ExpireAt: ", mytoken[0].Token.ExpireAt.UTC(), " now: ", time.Now().UTC(), "is ok? ", mytoken[0].Token.ExpireAt.Before(time.Now().UTC()))

				if mytoken[0].Token.ExpireAt.UTC().After(time.Now().UTC()) {
					return &mytoken[0]
				} else {
					mytoken[0].Delete()
					glog.Info("Deleting expired token.")
				}
			}
		}
	}

	return nil
}
