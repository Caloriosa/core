package user

import (
	"core/pkg/db"
	"core/types"
	"core/types/httptypes"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
	"time"
)

const COLLECTION_USERS = "users"

type UserResource struct {
	users map[string]types.User
}

func Register(container *restful.Container) {
	u := UserResource{map[string]types.User{}}
	u.Register(container)

	if err := db.MONGO_CONNECTION.Collection(COLLECTION_USERS).Collection().
		EnsureIndex(mgo.Index{Key: []string{"login"}, Unique: true}); err != nil {
		glog.Fatalf("Error ensuring index: ", err)
	}

	if err := db.MONGO_CONNECTION.Collection(COLLECTION_USERS).Collection().
		EnsureIndex(mgo.Index{Key: []string{"email"}, Unique: true}); err != nil {
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
	ws.Route(ws.PUT("{user-id}").To(u.putUser))

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
	result := db.MONGO_CONNECTION.Collection(COLLECTION_USERS).Find(nil)

	if result.Error != nil {
		glog.Warning("Got error fetching users: ", result.Error)
		httptypes.SendGeneralError(httptypes.EMPTY_CONTENT, response)
		return
	}

	var users []types.User
	result.Query.All(&users)

	for index := range users {
		users[index].Password = ""
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

	newUser.CreatedAt = time.Now().UTC()
	newUser.Activated = false
	newUser.Role = types.RoleUser

	err = db.MONGO_CONNECTION.Collection(COLLECTION_USERS).Save(&newUser)
	if err != nil {
		glog.Warning("Error saving new user to db: ", err)
		httptypes.SendDuplicated(httptypes.EMPTY_CONTENT, response)
	} else {
		glog.Infof("Created new user: ", newUser)
		httptypes.SendCreated(httptypes.EMPTY_CONTENT, response)
	}
}

func (u *UserResource) getSelf(request *restful.Request, response *restful.Response) {

}

func (u *UserResource) putSelf(request *restful.Request, response *restful.Response) {

}

func (u *UserResource) patchSelf(request *restful.Request, response *restful.Response) {

}

func (u *UserResource) deleteSelf(request *restful.Request, response *restful.Response) {

}

func (u *UserResource) getUser(request *restful.Request, response *restful.Response) {

}

func (u *UserResource) patchUser(request *restful.Request, response *restful.Response) {

}

func (u *UserResource) putUser(request *restful.Request, response *restful.Response) {

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

}

func (u *UserResource) activateUser(request *restful.Request, response *restful.Response) {

}

func (u *UserResource) resetPassword(request *restful.Request, response *restful.Response) {

}
