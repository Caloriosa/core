package types

import (
	"core/pkg/db"
	"core/pkg/error"
	"core/pkg/lib/devices"
	"core/types/httptypes"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

const COLLECTION_DEVICES = "devices"

func (d *Device) Save() *errors.CalError {
	if err := db.MONGO.Save(COLLECTION_DEVICES, d); err != nil {
		return &errors.CalError{Status: &httptypes.DATASOURCE_ERROR}
	}

	return nil
}

func (d *Device) ValidateNew() *errors.CalError {
	d.Name = deviceslib.GenerateDeviceUIDString()
	return nil
}

func GetDeviceFromRequest(r *http.Request) (*Device, *errors.CalError) {
	token := GetTokenFromRequest(r)
	if token == nil {
		return nil, &errors.CalError{Status: &httptypes.UNAUTHORIZED}
	}
	if token.Device == nil || token.Type != TokenDevice {
		return nil, &errors.CalError{Status: &httptypes.INVALID_TOKEN}
	}

	device := Device{}
	if err := GetDeviceById(token.Device.Hex(), &device); err != nil {
		return nil, &errors.CalError{Status: &httptypes.INVALID_TOKEN}
	}

	return &device, nil
}

func GetDeviceById(id string, device *Device) *errors.CalError {
	err := db.MONGO.FindById(COLLECTION_DEVICES, &device, id)
	if err != nil {
		return &errors.CalError{Status: &httptypes.NOT_FOUND}
	}

	return nil
}

func GetUsersDevices(user *User, devices *[]Device) *errors.CalError {
	err := db.MONGO.Get(COLLECTION_DEVICES, devices, bson.M{"user": user.Id.Hex()}, 0, 999)
	if err != nil {
		return &errors.CalError{Status: &httptypes.NOT_FOUND}
	}

	return nil
}
