package types

import (
	"github.com/go-bongo/bongo"
	"time"
	"gopkg.in/mgo.v2/bson"
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
	Sensor             *Sensor `json:"sensor"`
	MeasuredAt         time.Time `json:"measuredat"`
}

type Token struct {
	bongo.DocumentBase `bson:",inline"`
	Token              string `json:"token"`
	Type               string `json:"type"`
	ExpireAt           time.Time`json:"expireat"`
	User               *bson.ObjectId
	Device             *bson.ObjectId
}

const (
	TokenUser   = "user"
	TokenDevice = "device"
)

type User struct {
	bongo.DocumentBase `bson:",inline"`
	Login              string `json:"login"`
	Password           string `json:"password,omitempty"`
	Email              string `json:"email,omitempty"`
	Name               string `json:"name"`
	Activated          bool `json:"activated"`
	Role               string `json:"role"`
}

const (
	RoleUser  = "member"
	RoleAdmin = "admin"
)
