package types

import (
	"github.com/go-bongo/bongo"
	"time"
)

type Device struct {
	bongo.DocumentBase `bson:",inline"`
	Title              string
	Description        string
	Location           string
	FeaturedSensor     *Sensor
	Tags               []string
	CreatedAt          time.Time
}

type Sensor struct {
	bongo.DocumentBase `bson:",inline"`
	Device             *Device
	Alias              string
	Title              string
	Type               string
	Description        string
	CreatedAt          time.Time
}

const (
	SensorTemperature = "Temperature"
	SensorHumidity    = "Humidity"
	SensorWindSpeed   = "WindSpeed"
)

type Measurement struct {
	bongo.DocumentBase `bson:",inline"`
	Sensor             *Sensor
	MeasuredAt         time.Time
}

type Token struct {
	bongo.DocumentBase `bson:",inline"`
	Token              string
	Type               string
	CreatedAt          time.Time
	ExpireAt           time.Time
	User               *User
	Device             *Device
}

const (
	TokenUser   = "User"
	TokenDevice = "Device"
)

type User struct {
	bongo.DocumentBase `bson:",inline"`
	Login              string
	Password           string
	Email              string
	Name               string
	CreatedAt          time.Time
	Activated          bool
	Role               string
}

const (
	RoleUser  = "User"
	RoleAdmin = "Admin"
)
