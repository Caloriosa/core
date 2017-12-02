package types

import (
	"time"
)

type Device struct {
	Title          string
	Description    string
	Location       string
	FeaturedSensor *Sensor
	Tags           []string
	CreatedAt      Time
}

type Sensor struct {
	Alias       string
	Title       string
	Type        string
	Description string
	CreatedAt   Time
}

const (
	SensorTemperature = "Temperature"
	SensorHumidity    = "Humidity"
	SensorWindSpeed   = "WindSpeed"
)

type Measurement struct {
	Sensor     *Sensor
	MeasuredAt Time
}

type Token struct {
	Token     string
	Type      string
	CreatedAt Time
	ExpireAt  Time
	User      *User
	Device    *Device
}

const (
	TokenUser   = "User"
	TokenDevice = "Device"
)

type User struct {
	Login     string
	Password  string
	Email     string
	Name      string
	CreatedAt Time
}
