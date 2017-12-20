package tokenlib

import (
	"core/pkg/db"
	"core/pkg/error"
	"core/pkg/tools"
	"core/types"
	"core/types/httptypes"
	"gopkg.in/mgo.v2/bson"
	"core/pkg/config"
)

func GetAppFromToken(token string) *string {
	for _, app := range config.LoadedConfig.AppTokens {
		if app.Token == token {
			return &app.App
		}
	}
	return nil
}

func GetTokensForUser(user *types.User, tokens []*types.Token) *errors.CalError {

	if err := db.MONGO.Get(tools.COLLECTION_TOKENS, tokens, bson.M{"type": types.TokenUser}, 0, 999); err != nil {
		return &errors.CalError{Status: &httptypes.DATASOURCE_ERROR}
	}

	return nil
}

func GetToken(tokenid string, token *types.Token) *errors.CalError {
	if err := db.MONGO.Get(tools.COLLECTION_TOKENS, token, bson.M{"token": tokenid}, 0, 1); err != nil {
		return &errors.CalError{Status: &httptypes.DATASOURCE_ERROR}
	}

	return nil
}

func DeleteToken(token *types.Token) *errors.CalError {
	if err := db.MONGO.Connection.Collection(tools.COLLECTION_TOKENS).DeleteDocument(token); err != nil {
		return &errors.CalError{Status: &httptypes.DATASOURCE_ERROR}
	}

	return nil
}
