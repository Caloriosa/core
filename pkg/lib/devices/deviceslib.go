package deviceslib

import (
	"core/pkg/db"
	"core/pkg/error"
	"core/types"
	"core/types/httptypes"
	"gopkg.in/mgo.v2/bson"
)

const COLLECTION_DEVICES = "devices"

func GetUsersDevices(user *types.User, devices *[]types.Device) *errors.CalError {
	err := db.MONGO.Get(COLLECTION_DEVICES, devices, bson.M{"user": user.Id.Hex()}, 0, 999)
	if err != nil {
		return &errors.CalError{Status: &httptypes.NOT_FOUND}
	}

	return nil
}
