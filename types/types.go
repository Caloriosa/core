package types

import (
	"time"
)

type Device struct {
	ObjectID	   string
	Title          string
	Description    string
	Location       string
	FeaturedSensor *Sensor
	Tags           []string
	CreatedAt      time.Time
}

type Sensor struct {
	ObjectID	string
	Device		*Device
	Alias       string
	Title       string
	Type        string
	Description string
	CreatedAt   time.Time
}

const (
	SensorTemperature = "Temperature"
	SensorHumidity    = "Humidity"
	SensorWindSpeed   = "WindSpeed"
)

type Measurement struct {
	ObjectID	string
	Sensor     *Sensor
	MeasuredAt time.Time
}

type Token struct {
	ObjectID  string
	Token     string
	Type      string
	CreatedAt time.Time
	ExpireAt  time.Time
	User      *User
	Device    *Device
}

const (
	TokenUser   = "User"
	TokenDevice = "Device"
)

type User struct {
	ObjectID  string
	Login     string
	Password  string
	Email     string
	Name      string
	CreatedAt time.Time
	Activated bool
}
