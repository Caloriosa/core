package user

import (
	"core/pkg/activation"
	"core/pkg/db"
	"core/pkg/lib/devices"
	"core/pkg/lib/rest"
	"core/pkg/lib/tokenlib"
	"core/pkg/lib/user"
	"core/pkg/sanitization"
	"core/pkg/tools"
	"core/pkg/validation"
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

	if err := db.MONGO.Connection.Collection(userlib.COLLECTION_USERS).Collection().
		EnsureIndex(mgo.Index{Key: []string{"login"}, Unique: true, Collation: &collation}); err != nil {
		glog.Fatalf("Error ensuring index: ", err)
	}

	if err := db.MONGO.Connection.Collection(userlib.COLLECTION_USERS).Collection().
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
		httptypes.SendGeneralError(nil, response)
	}

	filters, err := tools.GetFilters(request.Request.URL.Query(), &types.User{})
	glog.Infof("Filters: ", filters)

	db.MONGO.Get(userlib.COLLECTION_USERS, &users, filters, page, limit)

	for index := range users {
		sanitization.UserSanitization(&users[index], true) // todo not a strict sanitization
	}

	httptypes.SendOK(users, response)
	glog.Info("Returned user list: ", users)
}

func (u *UserResource) createUser(request *restful.Request, response *restful.Response) {
	authedUser := request.Attribute(rest.ATTRIBUTE_AUTHED_USER).(*types.User)

	if authedUser.Role != types.RoleAdmin {
		httptypes.SendResponse(response, &httptypes.UNAUTHORIZED, nil)
		return
	}

	newUser := types.User{}
	err2 := request.ReadEntity(&newUser)
	if err2 != nil {
		httptypes.SendInvalidData(httptypes.EMPTY_CONTENT, response)
		glog.Warning("Error parsing a new user: ", err2)
		return
	}

	if err := userlib.CreateUser(&newUser); err != nil {
		httptypes.SendResponse(response, err.Status, nil)
		return
	}

	httptypes.SendResponse(response, &httptypes.HTTP_RESPONSE_OK, newUser)
}

func (u *UserResource) getSelf(request *restful.Request, response *restful.Response) {
	user := request.Attribute(rest.ATTRIBUTE_AUTHED_USER).(*types.User)

	sanitization.UserSanitization(user, false)

	httptypes.SendOK(user, response)
}

func (u *UserResource) patchSelf(request *restful.Request, response *restful.Response) {
	user := request.Attribute(rest.ATTRIBUTE_AUTHED_USER).(*types.User)

	sentUser := types.User{}

	decoder := json.NewDecoder(request.Request.Body)
	decoder.Decode(&sentUser) // TODO error handling

	validation.MergeChangedUser(user, &sentUser)

	if err := userlib.SaveUser(user); err != nil {
		httptypes.SendResponse(response, err.Status, nil)
		return
	}

	sanitization.UserSanitization(user, true) // todo not a strict sanitization
	httptypes.SendResponse(response, &httptypes.HTTP_RESPONSE_OK, user)
}

func (u *UserResource) deleteSelf(request *restful.Request, response *restful.Response) {
	user := request.Attribute(rest.ATTRIBUTE_AUTHED_USER).(*types.User)

	if err := userlib.DeleteUser(user); err != nil {
		httptypes.SendResponse(response, err.Status, nil)
		return
	}
}

func (u *UserResource) getUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter(rest.ATTRIBUTE_URL_USER)
	user := types.User{}
	if err := userlib.FindUserById(id, &user); err != nil {
		httptypes.SendResponse(response, err.Status, nil)
		return
	}

	sanitization.UserSanitization(&user, true) // todo not a strict sanitization

	httptypes.SendResponse(response, &httptypes.HTTP_RESPONSE_OK, user)
}

func (u *UserResource) patchUser(request *restful.Request, response *restful.Response) {
	user := request.Attribute(rest.ATTRIBUTE_URL_USER).(*types.User)
	sentUser := types.User{}

	decoder := json.NewDecoder(request.Request.Body)
	decoder.Decode(&sentUser) // TODO error handling

	validation.MergeChangedUser(user, &sentUser)

	if err := userlib.SaveUser(user); err != nil {
		httptypes.SendResponse(response, err.Status, nil)
		return
	}

	sanitization.UserSanitization(user, true) // todo not a strict sanitization
	httptypes.SendResponse(response, &httptypes.HTTP_RESPONSE_OK, user)
}

func (u *UserResource) getTokens(request *restful.Request, response *restful.Response) {
	tokens := []*types.Token{}
	user := request.Attribute(rest.ATTRIBUTE_AUTHED_USER).(*types.User)

	if err := tokenlib.GetTokensForUser(user, tokens); err != nil {
		httptypes.SendResponse(response, err.Status, nil)
		return
	}

	httptypes.SendOK(tokens, response)
}

func (u *UserResource) getToken(request *restful.Request, response *restful.Response) {
	tokenid := request.PathParameter("token-id")
	user := request.Attribute(rest.ATTRIBUTE_URL_USER).(*types.User)

	token := types.Token{}
	if err := tokenlib.GetToken(tokenid, &token); err != nil {
		httptypes.SendResponse(response, err.Status, nil)
		return
	}

	if token.User.Hex() != user.Id.Hex() {
		httptypes.SendResponse(response, &httptypes.INVALID_CREDENTIALS, nil)
		return
	}

	httptypes.SendOK(token, response)
}

func (u *UserResource) deleteToken(request *restful.Request, response *restful.Response) {
	tokenid := request.PathParameter("token-id")
	user := request.Attribute(rest.ATTRIBUTE_URL_USER).(*types.User)

	token := types.Token{}
	if err := tokenlib.GetToken(tokenid, &token); err != nil {
		httptypes.SendResponse(response, err.Status, nil)
		return
	}

	if token.User.Hex() != user.Id.Hex() {
		httptypes.SendResponse(response, &httptypes.INVALID_CREDENTIALS, nil)
		return
	}

	if err := tokenlib.DeleteToken(&token); err != nil {
		httptypes.SendResponse(response, err.Status, nil)
		return
	}

	httptypes.SendOK(token, response)
}

func (u *UserResource) getDevices(request *restful.Request, response *restful.Response) {
	devices := []types.Device{}

	user := request.Attribute(rest.ATTRIBUTE_URL_USER).(*types.User)

	if err := deviceslib.GetUsersDevices(user, &devices); err != nil {
		httptypes.SendResponse(response, err.Status, nil)
		return
	}

	httptypes.SendOK(devices, response)
}

func (u *UserResource) registerUser(request *restful.Request, response *restful.Response) {
	newUser := types.User{}
	err := request.ReadEntity(&newUser)
	if err != nil {
		httptypes.SendInvalidData(httptypes.EMPTY_CONTENT, response)
		glog.Warning("Error parsing a new user: ", err)
		return
	}

	if err := userlib.CreateUser(&newUser); err != nil {
		httptypes.SendResponse(response, err.Status, nil)
		return
	}

	newUser.Activated = false
	userlib.SaveUser(&newUser)

	httptypes.SendResponse(response, &httptypes.HTTP_RESPONSE_OK, newUser)

	go activation.SendValidationEmail(&newUser)
}

func (u *UserResource) activateUser(request *restful.Request, response *restful.Response) {
	foo := map[string]string{}

	decoder := json.NewDecoder(request.Request.Body)
	decoder.Decode(&foo) // TODO error handling

	if val, ok := foo["activation_token"]; ok && val != "" {
		if user, err := userlib.ActivateUserByActToken(val); err == nil {
			sanitization.UserSanitization(user, true)
			httptypes.SendResponse(response, &httptypes.HTTP_RESPONSE_OK, user)
			return
		}
	}

	httptypes.SendResponse(response, &httptypes.INVALID_TOKEN, nil)
}

func (u *UserResource) resetPassword(request *restful.Request, response *restful.Response) {
	// just generate a new token
	// TODO where do I get userID?
}
