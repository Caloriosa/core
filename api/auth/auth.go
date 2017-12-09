package auth

import (
	"github.com/emicklei/go-restful"
	"core/types"
	"core/pkg/db"
	user2 "core/api/user"
	"encoding/json"
	"core/types/httptypes"
	"github.com/golang/glog"
	"time"
	"core/pkg/tools"
	"gopkg.in/mgo.v2/bson"
)

const TOKEN_COLLECTION = "tokens"

type AuthResource struct {

}

func Register(container *restful.Container) {
	u := AuthResource{}
	u.Register(container)
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

	err := db.MONGO.Get(user2.COLLECTION_USERS, &foundUser, bson.M{"login": user.Login, "password": user.Password}, 0, 1)
	if err != nil {
		httptypes.SendGeneralError(nil, response)
		glog.Warning("auth: ", err)
		return
	}

	if len(foundUser) == 0 {
		httptypes.SendBadAuth(response)
		return
	}

	// create token
	token := types.Token{}
	token.Token = "123" // todo haha unique
	token.Type = types.TokenUser
	token.ExpireAt = time.Now().UTC()
	token.ExpireAt.Add(48 * time.Hour) // 2 days
	token.User = &foundUser[0].DocumentBase.Id
	token.Device = nil

	err = db.MONGO.Save(tools.COLLECTION_TOKENS, &token)
	if err != nil {
		httptypes.SendGeneralError(nil, response)
		glog.Warning("auth2: ", err)
		return
	}

	httptypes.SendOK(token, response)
}

func (u *AuthResource) refresh(request *restful.Request, response *restful.Response) {
	token := tools.GetToken(request.Request)
	token.ExpireAt = time.Now().UTC()
	token.ExpireAt.Add(48 * time.Hour)
	db.MONGO.Save(tools.COLLECTION_TOKENS, token)

	httptypes.SendOK(token, response)

}

func (u *AuthResource) logout(request *restful.Request, response *restful.Response) {
	token := tools.GetToken(request.Request)
	db.MONGO.Connection.Collection(tools.COLLECTION_TOKENS).DeleteDocument(token)

	httptypes.SendOK(nil, response)
}