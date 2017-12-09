package user

import (
	"core/pkg/db"
	"core/types"
	"core/types/httptypes"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
	"core/pkg/sanitization"
	"core/pkg/validation"
	"encoding/json"
	"core/pkg/tools"
)

const COLLECTION_USERS = "users"
const DEFAULT_LIMIT = 50

type UserResource struct {
	users map[string]types.User
}

func Register(container *restful.Container) {
	u := UserResource{map[string]types.User{}}
	u.Register(container)

	collation := mgo.Collation{Locale: "cs", Strength: 2}

	if err := db.MONGO.Connection.Collection(COLLECTION_USERS).Collection().

		EnsureIndex(mgo.Index{Key: []string{"login"}, Unique: true, Collation: &collation}); err != nil {
		glog.Fatalf("Error ensuring index: ", err)
	}

	if err := db.MONGO.Connection.Collection(COLLECTION_USERS).Collection().
		EnsureIndex(mgo.Index{Key: []string{"email"}, Unique: true, Collation: &collation}); err != nil {
		glog.Fatalf("Error ensuring index: ", err)
	}
}

func (u *UserResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/users")

	// users
	ws.Route(ws.GET("").To(u.listUsers))
	ws.Route(ws.POST("").To(u.createUser))

	// me
	ws.Route(ws.GET("me").To(u.getSelf))
	ws.Route(ws.PUT("me").To(u.putSelf))
	ws.Route(ws.PATCH("me").To(u.patchSelf))
	ws.Route(ws.DELETE("me").To(u.deleteSelf))

	// user id
	ws.Route(ws.GET("{user-id}").To(u.getUser))
	ws.Route(ws.PATCH("{user-id}").To(u.patchUser))

	// tokens
	ws.Route(ws.GET("users/{user-id}/tokens").To(u.getTokens))
	ws.Route(ws.GET("users/{user-id}/tokens/{token-id}").To(u.getToken))
	ws.Route(ws.DELETE("users/{user-id}/{token-id}").To(u.deleteToken))

	// devices
	ws.Route(ws.GET("users/{user-id}/devices").To(u.getDevices))

	// register
	ws.Route(ws.POST("register").To(u.registerUser))
	ws.Route(ws.POST("activate").To(u.activateUser))
	ws.Route(ws.POST("resetpass").To(u.resetPassword))

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


	db.MONGO.Get(COLLECTION_USERS, &users, filters, page, limit)

	for index := range users {
		sanitization.UserSanitization(&users[index], true) // todo not a strict sanitization
	}

	httptypes.SendOK(users, response)
	glog.Info("Returned user list: ", users)
}

func (u *UserResource) createUser(request *restful.Request, response *restful.Response) {
	newUser := types.User{}
	err := request.ReadEntity(&newUser)
	if err != nil {
		httptypes.SendInvalidData(httptypes.EMPTY_CONTENT, response)
		glog.Warning("Error parsing a new user: ", err)
		return
	}

	if err = validation.ValidateNewUser(&newUser); err != nil {
		httptypes.SendInvalidData(nil, response)
		glog.Infof("Error validating a new user: ", err)
		return
	}

	newUser.Activated = false
	newUser.Role = types.RoleUser

	err = db.MONGO.Save(COLLECTION_USERS, &newUser)
	if err != nil {
		glog.Warning("Error saving new user to db: ", err)
		httptypes.SendDuplicated(httptypes.EMPTY_CONTENT, response)
	} else {
		glog.Infof("Created new user: ", newUser)
		sanitization.UserSanitization(&newUser, true) // todo not a strict sanitization
		httptypes.SendCreated(newUser, response)
	}
}

func (u *UserResource) getSelf(request *restful.Request, response *restful.Response) {
	token := tools.GetToken(request.Request)
	if token == nil {
		httptypes.SendBadAuth(response)
		return
	}

	uid, err := tools.GetUserFromToken(token.Token)
	if uid == "" {
		httptypes.SendBadAuth(response)
		return
	}

	user := types.User{}
	err = db.MONGO.FindById(COLLECTION_USERS, &user, uid.Hex())
	if err != nil {
		httptypes.SendNotFound(nil, response)
		return
	}

	sanitization.UserSanitization(&user, false)

	httptypes.SendOK(user, response)
}

func (u *UserResource) putSelf(request *restful.Request, response *restful.Response) {

}

func (u *UserResource) patchSelf(request *restful.Request, response *restful.Response) {

}

func (u *UserResource) deleteSelf(request *restful.Request, response *restful.Response) {

}

func (u *UserResource) getUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	user := types.User{}
	err := db.MONGO.FindById(COLLECTION_USERS, &user, id)
	if err != nil {
		httptypes.SendNotFound(nil, response)
		return
	}

	sanitization.UserSanitization(&user, true) // todo not a strict sanitization

	httptypes.SendOK(user, response)
}

func (u *UserResource) patchUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")

	user := types.User{}
	sentUser := types.User{}
	err := db.MONGO.FindById(COLLECTION_USERS, &user, id)
	if err != nil {
		httptypes.SendNotFound(nil, response)
		return
	}

	decoder := json.NewDecoder(request.Request.Body)
	decoder.Decode(&sentUser) // TODO error handling

	validation.MergeChangedUser(&user, &sentUser)

	glog.Info("New user: ", user)

	db.MONGO.Save(COLLECTION_USERS, &user)
	sanitization.UserSanitization(&user, true) // todo not a strict sanitization
	httptypes.SendOK(user, response)
}

func (u *UserResource) getTokens(request *restful.Request, response *restful.Response) {

}

func (u *UserResource) getToken(request *restful.Request, response *restful.Response) {

}

func (u *UserResource) deleteToken(request *restful.Request, response *restful.Response) {

}

func (u *UserResource) getDevices(request *restful.Request, response *restful.Response) {

}

func (u *UserResource) registerUser(request *restful.Request, response *restful.Response) {
	user := types.User{}
	decoder := json.NewDecoder(request.Request.Body)
	decoder.Decode(&user) // TODO error handling
	user.Activated = false
}

func (u *UserResource) activateUser(request *restful.Request, response *restful.Response) {

}

func (u *UserResource) resetPassword(request *restful.Request, response *restful.Response) {

}
