package user

import (
	"core/types"
	"github.com/emicklei/go-restful"
	"net/http"
)

type UserResource struct {
	users map[string]types.User
}

func Register(container *restful.Container) {
	u := UserResource{map[string]types.User{}}
	u.Register(container)
}

func (u *UserResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/users")
	ws.Route(ws.GET("").To(u.listUsers))
	ws.Route(ws.POST("").To(u.createUser))

	container.Add(ws)
}

func (u *UserResource) listUsers(request *restful.Request, response *restful.Response) {
	response.WriteErrorString(http.StatusNotFound, "User could not be found.")
}

func (u *UserResource) createUser(request *restful.Request, response *restful.Response) {

}
