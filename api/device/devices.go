package device

import (
	"core/types"
	"github.com/emicklei/go-restful"
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
	ws.Route(ws.GET("{device-id}").To(u.getDevice))
	ws.Route(ws.POST("").To(u.createDevice))
	ws.Route(ws.PATCH("{device-id}").To(u.patchDevice))
	ws.Route(ws.DELETE("{device-id}").To(u.deleteDevice))

	ws.Route(ws.GET("my").To(u.listMyDevices))

	ws.Route(ws.GET("{device-id}/sensors").To(u.listSensors))
	ws.Route(ws.POST("{device-id}/sensors").To(u.createSensor))
	ws.Route(ws.GET("{device-id}/sensors/{sensor-id}").To(u.getSensor))
	ws.Route(ws.PATCH("{device-id}/sensors/{sensor-id}").To(u.patchSensor))
	ws.Route(ws.DELETE("{device-id}/sensors/{sensor-id}").To(u.deleteSensor))

	ws.Route(ws.GET("{device-id}/measurements").To(u.listMeasurements))
	ws.Route(ws.GET("{device-id}/measurements/{timestamp}").To(u.getMeasurement))
	ws.Route(ws.DELETE("{device-id}/measurements/{timestamp}").To(u.deleteMeasurement))

	ws.Route(ws.GET("{device-id}/measurements/history").To(u.getMeasurementsHistory))

	ws.Route(ws.GET("{device-id}/token").To(u.listTokens))
	ws.Route(ws.POST("{device-id}/token").To(u.createToken))
	ws.Route(ws.DELETE("{device-id}/token/{token-id}").To(u.deleteToken))

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

}

func (d *DeviceResource) getDevice(request *restful.Request, response *restful.Response) {

}

func (d *DeviceResource) createDevice(request *restful.Request, response *restful.Response) {

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