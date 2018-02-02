package user

import (
	"core/pkg/activation"
	"core/pkg/db"
	"core/pkg/lib/rest"
	"core/pkg/sanitization"
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
	var users []types.User

	var page, limit int
	var err error

	page, limit, err = tools.GetPagination(request.Request.URL.Query())
	if err != nil {
		httptypes.SendGeneralError(response)
	}

	filters, err := tools.GetFilters(request.Request.URL.Query(), &types.User{})
	glog.Infof("Filters: ", filters)

	db.MONGO.Get(types.COLLECTION_USERS, &users, filters, page, limit)

	for index := range users {
		sanitization.UserSanitization(&users[index], true) // todo not a strict sanitization
	}

	httptypes.SendOK(response, &users)
	glog.Info("Returned user list: ", users)
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

	if err := types.CreateUser(&newUser); err != nil {
		httptypes.SendError(response, err.Status)
		return
	}

	httptypes.SendOK(response, newUser)
}

func (u *UserResource) getSelf(request *restful.Request, response *restful.Response) {
	user := request.Attribute(rest.ATTRIBUTE_AUTHED_USER).(*types.User)

	sanitization.UserSanitization(user, false)

	httptypes.SendOK(response, &user)
}

func (u *UserResource) patchSelf(request *restful.Request, response *restful.Response) {
	user := request.Attribute(rest.ATTRIBUTE_AUTHED_USER).(*types.User)

	sentUser := types.User{}

	decoder := json.NewDecoder(request.Request.Body)
	decoder.Decode(&sentUser) // TODO error handling

	user.MergeIn(&sentUser)

	if err := user.Save(); err != nil {
		httptypes.SendError(response, err.Status)
		return
	}

	sanitization.UserSanitization(user, true) // todo not a strict sanitization
	httptypes.SendOK(response, user)
}

func (u *UserResource) deleteSelf(request *restful.Request, response *restful.Response) {
	user := request.Attribute(rest.ATTRIBUTE_AUTHED_USER).(*types.User)

	if err := user.Delete(); err != nil {
		httptypes.SendError(response, err.Status)
		return
	}
}

func (u *UserResource) getUser(request *restful.Request, response *restful.Response) {
	user := request.Attribute(rest.ATTRIBUTE_URL_USER).(*types.User)

	sanitization.UserSanitization(user, true) // todo not a strict sanitization

	httptypes.SendOK(response, user)
}

func (u *UserResource) patchUser(request *restful.Request, response *restful.Response) {
	user := request.Attribute(rest.ATTRIBUTE_URL_USER).(*types.User)
	sentUser := types.User{}

	decoder := json.NewDecoder(request.Request.Body)
	decoder.Decode(&sentUser) // TODO error handling

	user.MergeIn(&sentUser)

	if err := user.Save(); err != nil {
		httptypes.SendError(response, err.Status)
		return
	}

	sanitization.UserSanitization(user, true) // todo not a strict sanitization
	httptypes.SendOK(response, user)
}

func (u *UserResource) getTokens(request *restful.Request, response *restful.Response) {
	tokens := []*types.Token{}
	user := request.Attribute(rest.ATTRIBUTE_AUTHED_USER).(*types.User)

	if err := types.GetTokensForUser(user, tokens); err != nil {
		httptypes.SendError(response, err.Status)
		return
	}

	httptypes.SendOK(response, &tokens)
}

func (u *UserResource) getToken(request *restful.Request, response *restful.Response) {
	tokenid := request.PathParameter("token-id")
	user := request.Attribute(rest.ATTRIBUTE_URL_USER).(*types.User)

	token := types.Token{}
	if err := types.GetToken(tokenid, &token); err != nil {
		httptypes.SendError(response, err.Status)
		return
	}

	if token.User.Hex() != user.Id.Hex() {
		httptypes.SendError(response, &httptypes.INVALID_CREDENTIALS)
		return
	}

	httptypes.SendOK(response, token)
}

func (u *UserResource) deleteToken(request *restful.Request, response *restful.Response) {
	tokenid := request.PathParameter("token-id")
	user := request.Attribute(rest.ATTRIBUTE_URL_USER).(*types.User)

	token := types.Token{}
	if err := types.GetToken(tokenid, &token); err != nil {
		httptypes.SendError(response, err.Status)
		return
	}

	if token.User.Hex() != user.Id.Hex() {
		httptypes.SendError(response, &httptypes.INVALID_CREDENTIALS)
		return
	}

	if err := token.Delete(); err != nil {
		httptypes.SendError(response, err.Status)
		return
	}

	httptypes.SendOK(response, &token)
}

func (u *UserResource) getDevices(request *restful.Request, response *restful.Response) {
	devices := []types.Device{}

	user := request.Attribute(rest.ATTRIBUTE_URL_USER).(*types.User)

	if err := types.GetUsersDevices(user, &devices); err != nil {
		httptypes.SendError(response, err.Status)
		return
	}

	httptypes.SendOK(response, devices)
}

func (u *UserResource) registerUser(request *restful.Request, response *restful.Response) {
	newUser := types.User{}
	err := request.ReadEntity(&newUser)
	if err != nil {
		httptypes.SendError(response, nil)
		glog.Warning("Error parsing a new user: ", err)
		return
	}

	if err := types.CreateUser(&newUser); err != nil {
		httptypes.SendError(response, err.Status)
		return
	}

	newUser.Activated = false
	newUser.Save()

	httptypes.SendOK(response, &newUser)

	go activation.SendValidationEmail(&newUser)
}

func (u *UserResource) activateUser(request *restful.Request, response *restful.Response) {
	foo := map[string]string{}

	decoder := json.NewDecoder(request.Request.Body)
	decoder.Decode(&foo) // TODO error handling

	if val, ok := foo["activation_token"]; ok && val != "" {
		if user, err := types.ActivateUserByToken(val); err == nil {
			sanitization.UserSanitization(user, true)
			httptypes.SendOK(response, &user)
			return
		}
	}

	httptypes.SendError(response, &httptypes.INVALID_TOKEN)
}

func (u *UserResource) resetPassword(request *restful.Request, response *restful.Response) {
	// just generate a new token
	// TODO where do I get userID?
}
