package rest

import (
	"core/pkg/lib/tokenlib"
	"core/types"
	"core/types/httptypes"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
)

const ATTRIBUTE_AUTHED_USER = "user"
const ATTROBUTE_AUTHED_DEVICE = "device"
const ATTRIBUTE_URL_DEVICE = "url-device"
const ATTRIBUTE_URL_USER = "url-user"
const ATTRIBUTE_AUTHED_APP = "app"

func UserAuthFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	user, err := types.GetUserFromRequest(req.Request)
	if err != nil {
		httptypes.SendResponse(resp, err.Status, nil)
		glog.Info("Somebody tried accessing with unknown user token")
		return
	}

	req.SetAttribute(ATTRIBUTE_AUTHED_USER, user)

	chain.ProcessFilter(req, resp)
}

func DeviceAuthFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	device, err := types.GetDeviceFromRequest(req.Request)
	if err != nil {
		httptypes.SendResponse(resp, err.Status, nil)
		glog.Info("Somebody tried accessing with unknown device token")
		return
	}

	req.SetAttribute(ATTROBUTE_AUTHED_DEVICE, device)

	chain.ProcessFilter(req, resp)
}

func ExtractUserFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	id := req.PathParameter("user-id")
	user := types.User{}
	if err := types.GetUserById(id, &user); err != nil {
		httptypes.SendResponse(resp, err.Status, nil)
		return
	}

	req.SetAttribute(ATTRIBUTE_URL_USER, user)

	chain.ProcessFilter(req, resp)
}

func ExtractDeviceFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	id := req.PathParameter("device-id")
	device := types.Device{}
	if err := types.GetDeviceById(id, &device); err != nil {
		httptypes.SendResponse(resp, err.Status, nil)
		return
	}

	req.SetAttribute(ATTRIBUTE_URL_DEVICE, device)
	chain.ProcessFilter(req, resp)
}

func ValidateDeviceOwner(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	chain.ProcessFilter(req, resp)
}

func AppAuthFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	appsign := req.HeaderParameter(types.HEADER_APPSIGN)
	app := tokenlib.GetAppFromToken(appsign)
	if appsign == "" || app == nil {
		httptypes.SendResponse(resp, &httptypes.INVALID_SIGNATURE, nil)
		glog.Info("Somebody tried using wrong app token: ", appsign)
		return
	}

	req.SetAttribute(ATTRIBUTE_AUTHED_APP, app)

	chain.ProcessFilter(req, resp)
}
