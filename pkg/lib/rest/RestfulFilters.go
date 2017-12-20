package rest

import (
	"github.com/emicklei/go-restful"
	"core/pkg/lib/user"
	"core/types/httptypes"
	"core/types"
	"core/pkg/lib/tokenlib"
	"github.com/golang/glog"
)

const ATTRIBUTE_AUTHED_USER = "user"
const ATTRIBUTE_URL_USER = "url-user"
const ATTRIBUTE_AUTHED_APP = "app"

func UserAuthFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	user, err := userlib.GetUserFromRequest(req.Request)
	if err != nil {
		httptypes.SendResponse(resp, err.Status, nil)
		glog.Info("Somebody tried accessing with unknown user token")
		return
	}

	req.SetAttribute(ATTRIBUTE_AUTHED_USER, &user)

	chain.ProcessFilter(req, resp)
}

func ExtractUserFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	id := req.PathParameter("user-id")
	user := types.User{}
	if err := userlib.FindUserById(id, &user); err != nil {
		httptypes.SendResponse(resp, err.Status, nil)
		return
	}

	req.SetAttribute(ATTRIBUTE_URL_USER, &user)

	chain.ProcessFilter(req, resp)
}

func ExtractDeviceFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	chain.ProcessFilter(req, resp)
}

func ValidateDeviceOwner(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	chain.ProcessFilter(req, resp)
}

func AppAuthFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	appsign := req.HeaderParameter(types.HEADER_APPSIGN)
	app := tokenlib.GetAppFromToken(appsign)
	if appsign == "" || app == nil {
		httptypes.SendResponse(resp, &httptypes.UNAUTHORIZED, nil)
		glog.Info("Somebody tried using wrong app token: ", appsign)
		return
	}

	req.SetAttribute(ATTRIBUTE_AUTHED_APP, app)

	chain.ProcessFilter(req, resp)
}