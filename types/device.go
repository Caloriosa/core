package types

import (
	"core/pkg/db"
	"core/pkg/error"
	"core/pkg/lib/devices"
	"core/types/httptypes"
	"github.com/go-bongo/bongo"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

const COLLECTION_DEVICES = "devices"

type DeviceDB struct {
	bongo.DocumentBase `bson:",inline"`
	Device             Device `json:",inline" bson:",inline"`
}

func (d *DeviceDB) Save() *errors.CalError {
	if d.Device.Name == "" {
		d.Device.Name = deviceslib.GenerateDeviceUIDString()
	}

	if err := db.MONGO.Save(COLLECTION_DEVICES, d); err != nil {
		return &errors.CalError{Status: &httptypes.DATASOURCE_ERROR}
	}

	return nil
}

func (d *Device) ValidateNew() *errors.CalError {
	return nil
}

func (d *DeviceDB) PrepareToSend() *Device {
	device := new(Device)
	*device = d.Device
	device.UID = d.Id.Hex()
	device.CreatedAt = d.Created
	device.ModifiedAt = d.Modified
	return device
}

func CreateNewDevice(d *Device) (*DeviceDB, *errors.CalError) {
	db := DeviceDB{Device: *d}
	if err := db.Save(); err != nil {
		return nil, &errors.CalError{Status: &httptypes.DATASOURCE_ERROR}
	}

	return &db, nil
}

func GetDeviceFromRequest(r *http.Request) (*DeviceDB, *errors.CalError) {
	token := GetTokenFromRequest(r)
	if token == nil {
		return nil, &errors.CalError{Status: &httptypes.UNAUTHORIZED}
	}
	if token.Token.Device == nil || token.Token.Type != TokenDevice {
		return nil, &errors.CalError{Status: &httptypes.INVALID_TOKEN}
	}

	device := DeviceDB{}
	if err := GetDeviceById(token.Token.Device.Hex(), &device); err != nil {
		return nil, &errors.CalError{Status: &httptypes.INVALID_TOKEN}
	}

	return &device, nil
}

func GetDeviceById(id string, device *DeviceDB) *errors.CalError {
	err := db.MONGO.FindById(COLLECTION_DEVICES, &device, id)
	if err != nil {
		return &errors.CalError{Status: &httptypes.NOT_FOUND}
	}

	return nil
}

func GetUsersDevices(user *UserDB, devices *[]DeviceDB) *errors.CalError {
	err := db.MONGO.Get(COLLECTION_DEVICES, devices, bson.M{"user": user.Id.Hex()}, 0, 999)
	if err != nil {
		return &errors.CalError{Status: &httptypes.NOT_FOUND}
	}

	return nil
}
