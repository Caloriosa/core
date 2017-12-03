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

	container.Add(ws)
}
