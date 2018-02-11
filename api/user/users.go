package user

import (
	"core/pkg/activation"
	"core/pkg/db"
	"core/pkg/lib/rest"
	"core/pkg/tools"
	"core/types"
	"core/types/httptypes"
	"encoding/json"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
)

const DEFAULT_LIMIT = 50

type UserResource struct {
	users map[string]types.User
}

func Register(container *restful.Container) {
	u := UserResource{map[string]types.User{}}
	u.Register(container)

	collation := mgo.Collation{Locale: "cs", Strength: 2}

	if err := db.MONGO.Connection.Collection(types.COLLECTION_USERS).Collection().
		EnsureIndex(mgo.Index{Key: []string{"login"}, Unique: true, Collation: &collation}); err != nil {
		glog.Fatalf("Error ensuring index: ", err)
	}

	if err := db.MONGO.Connection.Collection(types.COLLECTION_USERS).Collection().
		EnsureIndex(mgo.Index{Key: []string{"email"}, Unique: true, Collation: &collation}); err != nil {
		glog.Fatalf("Error ensuring index: ", err)
	}
}

func (u *UserResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/users")

	// users
	ws.Route(ws.GET("").To(u.listUsers))
	ws.Route(ws.POST("").Filter(rest.UserAuthFilter).To(u.createUser))

	// me
	ws.Route(ws.GET("me").Filter(rest.UserAuthFilter).To(u.getSelf))
	ws.Route(ws.PATCH("me").Filter(rest.UserAuthFilter).To(u.patchSelf))
	ws.Route(ws.DELETE("me").Filter(rest.UserAuthFilter).To(u.deleteSelf))

	// user id
	ws.Route(ws.GET("{user-id}").Filter(rest.ExtractUserFilter).To(u.getUser))
	ws.Route(ws.PATCH("{user-id}").Filter(rest.ExtractUserFilter).To(u.patchUser))

	// tokens
	ws.Route(ws.GET("users/{user-id}/tokens").Filter(rest.UserAuthFilter).Filter(rest.ExtractUserFilter).To(u.getTokens))
	ws.Route(ws.GET("users/{user-id}/tokens/{token-id}").Filter(rest.UserAuthFilter).Filter(rest.ExtractUserFilter).To(u.getToken))
	ws.Route(ws.DELETE("users/{user-id}/{token-id}").Filter(rest.UserAuthFilter).Filter(rest.ExtractUserFilter).To(u.deleteToken))

	// devices
	ws.Route(ws.GET("users/{user-id}/devices").Filter(rest.ExtractUserFilter).To(u.getDevices))

	// register
	ws.Route(ws.POST("register").Filter(rest.AppAuthFilter).To(u.registerUser))
	ws.Route(ws.POST("activate").Filter(rest.AppAuthFilter).To(u.activateUser))
	ws.Route(ws.POST("resetpass").Filter(rest.AppAuthFilter).To(u.resetPassword))

	container.Add(ws)

}

func (u *UserResource) listUsers(request *restful.Request, response *restful.Response) {
	var users []types.UserDB

	var page, limit int
	var err error

	page, limit, err = tools.GetPagination(request.Request.URL.Query())
	if err != nil {
		httptypes.SendGeneralError(response)
	}

	filters, err := tools.GetFilters(request.Request.URL.Query(), &types.User{})
	glog.Infof("Filters: ", filters)

	db.MONGO.Get(types.COLLECTION_USERS, &users, filters, page, limit)

	var usersToSend = make([]types.User, len(users))

	for index := range users {
		usersToSend[index] = *users[index].PrepareToSend(true)
	}

	httptypes.SendOK(response, &usersToSend)
	glog.Info("Returned user list: ", &usersToSend)
}

func (u *UserResource) createUser(request *restful.Request, response *restful.Response) {
	authedUser := request.Attribute(rest.ATTRIBUTE_AUTHED_USER).(*types.User)

	if authedUser.Role != types.RoleAdmin {
		httptypes.SendError(response, &httptypes.UNAUTHORIZED)
		return
	}

	newUser := types.User{}
	err2 := request.ReadEntity(&newUser)
	if err2 != nil {
		httptypes.SendError(response, nil)
		glog.Warning("Error parsing a new user: ", err2)
		return
	}

	createdUser, err := types.CreateUser(&newUser)

	if err != nil {
		httptypes.SendError(response, err.Status)
		return
	}

	httptypes.SendOK(response, createdUser.PrepareToSend(false))
}

func (u *UserResource) getSelf(request *restful.Request, response *restful.Response) {
	user := request.Attribute(rest.ATTRIBUTE_AUTHED_USER).(*types.UserDB)

	httptypes.SendOK(response, user.PrepareToSend(false))
}

func (u *UserResource) patchSelf(request *restful.Request, response *restful.Response) {
	user := request.Attribute(rest.ATTRIBUTE_AUTHED_USER).(*types.UserDB)

	sentUser := types.User{}

	decoder := json.NewDecoder(request.Request.Body)
	decoder.Decode(&sentUser) // TODO error handling

	user.User.MergeIn(&sentUser)

	if err := user.Save(); err != nil {
		httptypes.SendError(response, err.Status)
		return
	}

	httptypes.SendOK(response, user.PrepareToSend(false))
}

func (u *UserResource) deleteSelf(request *restful.Request, response *restful.Response) {
	user := request.Attribute(rest.ATTRIBUTE_AUTHED_USER).(*types.UserDB)

	if err := user.Delete(); err != nil {
		httptypes.SendError(response, err.Status)
		return
	}

	httptypes.SendOK(response, nil)
}

func (u *UserResource) getUser(request *restful.Request, response *restful.Response) {
	user := request.Attribute(rest.ATTRIBUTE_URL_USER).(*types.UserDB)

	httptypes.SendOK(response, user.PrepareToSend(true))
}

func (u *UserResource) patchUser(request *restful.Request, response *restful.Response) {
	user := request.Attribute(rest.ATTRIBUTE_URL_USER).(*types.UserDB)
	sentUser := types.User{}

	decoder := json.NewDecoder(request.Request.Body)
	decoder.Decode(&sentUser) // TODO error handling

	user.User.MergeIn(&sentUser)

	if err := user.Save(); err != nil {
		httptypes.SendError(response, err.Status)
		return
	}

	httptypes.SendOK(response, user.PrepareToSend(true))
}

func (u *UserResource) getTokens(request *restful.Request, response *restful.Response) {
	tokens := []*types.TokenDB{}
	user := request.Attribute(rest.ATTRIBUTE_AUTHED_USER).(*types.UserDB)

	if err := types.GetTokensForUser(user, tokens); err != nil {
		httptypes.SendError(response, err.Status)
		return
	}

	tokensToSend := make([]types.Token, len(tokens))
	for index := range tokens {
		tokensToSend[index] = *tokens[index].PrepareToSend()
	}

	httptypes.SendOK(response, &tokensToSend)
}

func (u *UserResource) getToken(request *restful.Request, response *restful.Response) {
	tokenid := request.PathParameter("token-id")
	user := request.Attribute(rest.ATTRIBUTE_URL_USER).(*types.UserDB)

	token := types.TokenDB{}
	if err := types.GetToken(tokenid, &token); err != nil {
		httptypes.SendError(response, err.Status)
		return
	}

	if token.Token.User.Hex() != user.Id.Hex() {
		httptypes.SendError(response, &httptypes.INVALID_CREDENTIALS)
		return
	}

	httptypes.SendOK(response, token.PrepareToSend())
}

func (u *UserResource) deleteToken(request *restful.Request, response *restful.Response) {
	tokenid := request.PathParameter("token-id")
	user := request.Attribute(rest.ATTRIBUTE_URL_USER).(*types.UserDB)

	token := types.TokenDB{}
	if err := types.GetToken(tokenid, &token); err != nil {
		httptypes.SendError(response, err.Status)
		return
	}

	if token.Token.User.Hex() != user.Id.Hex() {
		httptypes.SendError(response, &httptypes.INVALID_CREDENTIALS)
		return
	}

	if err := token.Delete(); err != nil {
		httptypes.SendError(response, err.Status)
		return
	}

	httptypes.SendOK(response, nil)
}

func (u *UserResource) getDevices(request *restful.Request, response *restful.Response) {
	devices := []types.DeviceDB{}

	user := request.Attribute(rest.ATTRIBUTE_URL_USER).(*types.UserDB)

	if err := types.GetUsersDevices(user, &devices); err != nil {
		httptypes.SendError(response, err.Status)
		return
	}

	devicesToSend := make([]types.Device, len(devices))
	for index := range devices {
		devicesToSend[index] = *devices[index].PrepareToSend()
	}

	httptypes.SendOK(response, devicesToSend)
}

func (u *UserResource) registerUser(request *restful.Request, response *restful.Response) {
	newUser := types.User{}

	if err := request.ReadEntity(&newUser); err != nil {
		httptypes.SendError(response, nil)
		glog.Warning("Error parsing a new user: ", err)
		return
	}

	createdUser, err := types.CreateUser(&newUser)

	if err != nil {
		httptypes.SendError(response, err.Status)
		return
	}

	createdUser.User.Activated = false
	createdUser.Save()

	httptypes.SendOK(response, createdUser.PrepareToSend(false))

	go activation.SendValidationEmail(&newUser)
}

func (u *UserResource) activateUser(request *restful.Request, response *restful.Response) {
	foo := map[string]string{}

	decoder := json.NewDecoder(request.Request.Body)
	decoder.Decode(&foo) // TODO error handling

	if val, ok := foo["activation_token"]; ok && val != "" {
		if user, err := types.ActivateUserByToken(val); err == nil {
			httptypes.SendOK(response, user.PrepareToSend(true))
			return
		}
	}

	httptypes.SendError(response, &httptypes.INVALID_TOKEN)
}

func (u *UserResource) resetPassword(request *restful.Request, response *restful.Response) {
	// just generate a new token
	// TODO where do I get userID?
}
