package auth

import (
	"core/pkg/db"
	"core/pkg/lib/user"
	"core/pkg/tools"
	"core/types"
	"core/types/httptypes"
	"encoding/json"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const COLLECTION_TOKENS = "tokens"

type AuthResource struct {
}

func Register(container *restful.Container) {
	u := AuthResource{}
	u.Register(container)

	collation := mgo.Collation{Locale: "cs", Strength: 2}

	if err := db.MONGO.Connection.Collection(COLLECTION_TOKENS).Collection().
		EnsureIndex(mgo.Index{Key: []string{"token"}, Unique: true, Collation: &collation}); err != nil {
		glog.Fatalf("Error ensuring index: ", err)
	}
}

func (u *AuthResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/auth")

	// auth
	ws.Route(ws.POST("").To(u.auth))
	ws.Route(ws.PATCH("").To(u.refresh))
	ws.Route(ws.DELETE("").To(u.logout))

	container.Add(ws)

}

func (u *AuthResource) auth(request *restful.Request, response *restful.Response) {
	user := types.User{}
	decoder := json.NewDecoder(request.Request.Body)
	decoder.Decode(&user)
	foundUser := []types.User{}

	if user.Password == "" || user.Login == "" {
		httptypes.SendBadAuth(response)
		return
	}

	err := db.MONGO.Get(userlib.COLLECTION_USERS, &foundUser, bson.M{"login": user.Login, "password": user.Password}, 0, 1)
	if err != nil {
		httptypes.SendGeneralError(nil, response)
		glog.Warning("auth: ", err)
		return
	}

	if len(foundUser) == 0 {
		httptypes.SendResponse(response, &httptypes.INVALID_CREDENTIALS, nil)
		return
	}

	// create token
	token := types.Token{}
	token.Token = tools.RandStringRunes(types.LENGTH_TOKEN)
	token.Type = types.TokenUser
	token.ExpireAt = time.Now().UTC().Add(48 * time.Hour)
	token.User = &foundUser[0].DocumentBase.Id
	token.Device = nil

	err = db.MONGO.Save(COLLECTION_TOKENS, &token)
	if err != nil {
		httptypes.SendGeneralError(nil, response)
		glog.Warning("auth2: ", err)
		return
	}

	glog.Info("Creating token: ", token)
	httptypes.SendOK(token, response)
}

func (u *AuthResource) refresh(request *restful.Request, response *restful.Response) {
	token := tools.GetToken(request.Request)
	token.ExpireAt = time.Now().UTC().Add(48 * time.Hour)
	db.MONGO.Save(COLLECTION_TOKENS, token)

	httptypes.SendOK(token, response)

}

func (u *AuthResource) logout(request *restful.Request, response *restful.Response) {
	token := tools.GetToken(request.Request)
	if token == nil {
		glog.Info("Invalid token")
		httptypes.SendResponse(response, &httptypes.INVALID_TOKEN, nil)
		return
	}
	db.MONGO.Connection.Collection(COLLECTION_TOKENS).DeleteDocument(token)

	httptypes.SendOK(nil, response)
}
