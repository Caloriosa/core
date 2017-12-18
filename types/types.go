package types

import (
	"github.com/go-bongo/bongo"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Device struct {
	bongo.DocumentBase `bson:",inline"`
	Title              string         `json:"title"`
	Description        string         `json:"description"`
	Location           string         `json:"location"`
	FeaturedSensor     *bson.ObjectId `json:"featuredsensor"`
	Tags               []string       `json:"tags"`
	User               *bson.ObjectId `json:"user"`
}

type Sensor struct {
	bongo.DocumentBase `bson:",inline"`
	Device             *bson.ObjectId `json:"device"`
	Alias              string         `json:"alias"`
	Title              string         `json:"title"`
	Type               string         `json:"type"`
	Description        string         `json:"description"`
}

const (
	SensorTemperature = "Temperature"
	SensorHumidity    = "Humidity"
	SensorWindSpeed   = "WindSpeed"
)

type Measurement struct {
	bongo.DocumentBase `bson:",inline"`
	Sensor             *Sensor   `json:"sensor"`
	MeasuredAt         time.Time `json:"measuredat"`
}

type Token struct {
	bongo.DocumentBase `bson:",inline"`
	Token              string         `json:"token"`
	Type               string         `json:"type"`
	ExpireAt           time.Time      `json:"expireat"`
	User               *bson.ObjectId `json:"user"`
	Device             *bson.ObjectId `json:"device"`
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
	Activated          bool   `json:"activated"`
	Role               string `json:"role"`
	ActivationKey      *string `json:"-"`
}

const (
	RoleUser  = "member"
	RoleAdmin = "admin"
)

const LENGTH_TOKEN = 32