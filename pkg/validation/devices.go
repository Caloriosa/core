package validation

import (
	"core/pkg/error"
	"core/pkg/lib/devices"
	"core/types"
)

func ValidateNewDevice(device *types.Device) *errors.CalError {
	device.Title = deviceslib.GenerateDeviceUIDString()
	return nil
}
