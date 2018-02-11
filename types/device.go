package types

import (
	"core/pkg/db"
	"core/pkg/error"
	"core/pkg/lib/devices"
	"core/types/httptypes"
	"github.com/go-bongo/bongo"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

const COLLECTION_DEVICES = "devices"

type DeviceDB struct {
	bongo.DocumentBase `bson:",inline"`
	Device             Device `json:",inline" bson:",inline"`
}

type Map struct {
	Position struct {
		Latitude  float32 `json:"lat" bson:"latitude"`
		Longitude float32 `json:"lng" bson:"longitude"`
	} `json:"position" bson:"_id"`

	Devices []struct {
		UID   bson.ObjectId `json:"uid" bson:"uid"`
		Title string        `json:"title" bson:"title"`
		Name  string        `json:"name" bson:"name"`
	} `json:"devices"`
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
	user := UserDB{}
	if err := GetUserById(d.Device.User.Hex(), &user); err != nil {
		return nil // uh
	}
	*device = d.Device
	device.UID = d.Id.Hex()
	device.CreatedAt = d.Created
	device.ModifiedAt = d.Modified
	device.UserObj = user.PrepareToSend(true)
	device.User = nil
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

func GetDeviceByName(id string, device *DeviceDB) *errors.CalError {
	devices := []DeviceDB{}
	if err := db.MONGO.Get(COLLECTION_DEVICES, &devices, bson.M{"name": id}, 0, 1); err != nil || len(devices) == 0 {
		return &errors.CalError{Status: &httptypes.NOT_FOUND}
	}

	*device = devices[0]

	return nil
}

func GetUsersDevices(user *UserDB, devices *[]DeviceDB) *errors.CalError {
	err := db.MONGO.Get(COLLECTION_DEVICES, devices, bson.M{"user": user.Id.Hex()}, 0, 999)
	if err != nil {
		return &errors.CalError{Status: &httptypes.NOT_FOUND}
	}

	return nil
}

func GetMap() *[]Map {
	mapreq := []bson.M{
		{
			"$group": bson.M{
				"_id": "$position",
				"devices": bson.M{
					"$push": bson.M{
						"uid":   "$_id",
						"title": "$title",
						"name":  "$name",
					},
				},
			},
		},
	}

	results := db.MONGO.Connection.Collection(COLLECTION_DEVICES).Collection().Pipe(mapreq)
	mapped := []Map{}
	results.All(&mapped)
	glog.Info("Got this map: ", mapped)
	return &mapped
}
