package device

import (
	"core/pkg/db"
	"core/pkg/lib/rest"
	"core/pkg/tools"
	"core/types"
	"core/types/httptypes"
	"encoding/json"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
)

type DeviceResource struct {
	users map[string]types.Device
}

func Register(container *restful.Container) {
	u := DeviceResource{map[string]types.Device{}}
	u.Register(container)
}

func (u *DeviceResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/devices")
	ws.Route(ws.GET("").To(u.listDevices))
	ws.Route(ws.GET("{device-id}").Filter(rest.ExtractDeviceFilter).To(u.getDevice))
	ws.Route(ws.POST("").Filter(rest.UserAuthFilter).To(u.createDevice))
	ws.Route(ws.PATCH("{device-id}").Filter(rest.ExtractDeviceFilter).Filter(rest.UserAuthFilter).To(u.patchDevice))
	ws.Route(ws.DELETE("{device-id}").Filter(rest.ExtractDeviceFilter).Filter(rest.UserAuthFilter).To(u.deleteDevice))

	ws.Route(ws.GET("my").Filter(rest.UserAuthFilter).To(u.listMyDevices))

	ws.Route(ws.GET("{device-id}/sensors").Filter(rest.ExtractDeviceFilter).To(u.listSensors))
	ws.Route(ws.POST("{device-id}/sensors").Filter(rest.ExtractDeviceFilter).Filter(rest.UserAuthFilter).To(u.createSensor))
	ws.Route(ws.GET("{device-id}/sensors/{sensor-id}").Filter(rest.ExtractDeviceFilter).To(u.getSensor))
	ws.Route(ws.PATCH("{device-id}/sensors/{sensor-id}").Filter(rest.ExtractDeviceFilter).Filter(rest.UserAuthFilter).To(u.patchSensor))
	ws.Route(ws.DELETE("{device-id}/sensors/{sensor-id}").Filter(rest.ExtractDeviceFilter).Filter(rest.UserAuthFilter).To(u.deleteSensor))

	ws.Route(ws.GET("{device-id}/measurements").Filter(rest.ExtractDeviceFilter).To(u.listMeasurements))
	ws.Route(ws.GET("{device-id}/measurements/{timestamp}").Filter(rest.ExtractDeviceFilter).To(u.getMeasurement))
	ws.Route(ws.DELETE("{device-id}/measurements/{timestamp}").Filter(rest.ExtractDeviceFilter).Filter(rest.UserAuthFilter).To(u.deleteMeasurement))

	ws.Route(ws.GET("{device-id}/measurements/history").Filter(rest.ExtractDeviceFilter).To(u.getMeasurementsHistory))

	ws.Route(ws.GET("{device-id}/token").Filter(rest.ExtractDeviceFilter).Filter(rest.UserAuthFilter).To(u.listTokens))
	ws.Route(ws.POST("{device-id}/token").Filter(rest.ExtractDeviceFilter).Filter(rest.UserAuthFilter).To(u.createToken))
	ws.Route(ws.DELETE("{device-id}/token/{token-id}").Filter(rest.ExtractDeviceFilter).Filter(rest.UserAuthFilter).To(u.deleteToken))

	ws.Route(ws.GET("me").To(u.getSelf))

	ws.Route(ws.GET("me/sensors").To(u.listSelfSensors))
	ws.Route(ws.POST("me/sensors").To(u.createSelfSensor))
	ws.Route(ws.GET("me/sensors/{sensor-id}").To(u.getSelfSensor))
	ws.Route(ws.GET("me/sensors/{sensor-id}").To(u.patchSelfSensor))
	ws.Route(ws.DELETE("me/sensors/{sensor-id}").To(u.deleteSelfSensor))

	container.Add(ws)
}

func (d *DeviceResource) getSelf(request *restful.Request, response *restful.Response) {

}

func (d *DeviceResource) listDevices(request *restful.Request, response *restful.Response) {
	var devices []types.Device

	var page, limit int
	var err error

	page, limit, err = tools.GetPagination(request.Request.URL.Query())
	if err != nil {
		httptypes.SendGeneralError(response)
		return
	}

	filters, err := tools.GetFilters(request.Request.URL.Query(), &types.User{})
	glog.Infof("Filters: ", filters)

	db.MONGO.Get(types.COLLECTION_DEVICES, &devices, filters, page, limit)

	httptypes.SendOK(response, &devices)
	glog.Info("Returned devices list: ", devices)
}

func (d *DeviceResource) getDevice(request *restful.Request, response *restful.Response) {
	httptypes.SendOK(response, request.Attribute(rest.ATTRIBUTE_URL_DEVICE))
}

func (d *DeviceResource) createDevice(request *restful.Request, response *restful.Response) {
	user := request.Attribute("user").(*types.User)

	device := types.Device{}
	decoder := json.NewDecoder(request.Request.Body)
	decoder.Decode(&device)

	if err := device.ValidateNew(); err != nil {
		httptypes.SendError(response, err.Status)
		return
	}

	device.User = &user.Id
	device.FeaturedSensor = nil

	if err := device.Save(); err != nil {
		httptypes.SendError(response, err.Status)
		return
	}

	httptypes.SendOK(response, &device)
	glog.Info("Created a new device: ", device)

}

func (d *DeviceResource) patchDevice(request *restful.Request, response *restful.Response) {

}

func (d *DeviceResource) deleteDevice(request *restful.Request, response *restful.Response) {

}

func (d *DeviceResource) listMyDevices(request *restful.Request, response *restful.Response) {

}

func (d *DeviceResource) listSensors(request *restful.Request, response *restful.Response) {

}

func (d *DeviceResource) createSensor(request *restful.Request, response *restful.Response) {

}

func (d *DeviceResource) getSensor(request *restful.Request, response *restful.Response) {

}

func (d *DeviceResource) patchSensor(request *restful.Request, response *restful.Response) {

}

func (d *DeviceResource) deleteSensor(request *restful.Request, response *restful.Response) {

}

func (d *DeviceResource) listMeasurements(request *restful.Request, response *restful.Response) {

}

func (d *DeviceResource) getMeasurement(request *restful.Request, response *restful.Response) {

}

func (d *DeviceResource) deleteMeasurement(request *restful.Request, response *restful.Response) {

}

func (d *DeviceResource) getMeasurementsHistory(request *restful.Request, response *restful.Response) {

}

func (d *DeviceResource) listTokens(request *restful.Request, response *restful.Response) {

}

func (d *DeviceResource) createToken(request *restful.Request, response *restful.Response) {

}

func (d *DeviceResource) deleteToken(request *restful.Request, response *restful.Response) {

}

func (d *DeviceResource) listSelfSensors(request *restful.Request, response *restful.Response) {

}

func (d *DeviceResource) getSelfSensor(request *restful.Request, response *restful.Response) {

}

func (d *DeviceResource) deleteSelfSensor(request *restful.Request, response *restful.Response) {

}

func (d *DeviceResource) createSelfSensor(request *restful.Request, response *restful.Response) {

}

func (d *DeviceResource) patchSelfSensor(request *restful.Request, response *restful.Response) {

}
