package validation

import (
	"core/pkg/error"
	"core/types"
	"core/pkg/lib/devices"
)

func ValidateNewDevice(device *types.Device) *errors.CalError {
	device.Title = deviceslib.GenerateDeviceUIDString()
	return nil
}
