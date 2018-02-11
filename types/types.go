package types

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Device struct {
	UID         string    `json:"uid" bson:"-"`
	CreatedAt   time.Time `json:"created" bson:"-"`
	ModifiedAt  time.Time `json:"modified" bson:"-"`
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Position    struct {
		Latitude  float32 `json:"lat"`
		Longitude float32 `json:"lng"`
	} `json:"position"`
	FeaturedSensor *bson.ObjectId `json:"featuredsensor"`
	Tags           []string       `json:"tags"`
	User           *bson.ObjectId `json:",omitempty"`
	UserObj        *User          `json:"user" bson:"-"`
}

type Sensor struct {
	UID         string         `json:"uid" bson:"-"`
	CreatedAt   time.Time      `json:"created" bson:"-"`
	ModifiedAt  time.Time      `json:"modified" bson:"-"`
	Device      *bson.ObjectId `json:"device"`
	Alias       string         `json:"alias"`
	Title       string         `json:"title"`
	Type        string         `json:"type"`
	Description string         `json:"description"`
}

const (
	SensorTemperature = "Temperature"
	SensorHumidity    = "Humidity"
	SensorWindSpeed   = "WindSpeed"
)

type Measurement struct {
	UID        string    `json:"uid" bson:"-"`
	CreatedAt  time.Time `json:"created" bson:"-"`
	ModifiedAt time.Time `json:"modified" bson:"-"`
	Sensor     *Sensor   `json:"sensor"`
	MeasuredAt time.Time `json:"measuredat"`
	Value      uint64    `json:"value"`
}

type Token struct {
	UID        string         `json:"uid" bson:"-"`
	CreatedAt  time.Time      `json:"created" bson:"-"`
	ModifiedAt time.Time      `json:"modified" bson:"-"`
	Token      string         `json:"token"`
	Type       string         `json:"type"`
	ExpireAt   time.Time      `json:"expireat"`
	User       *bson.ObjectId `json:"user"`
	Device     *bson.ObjectId `json:"device"`
}

const (
	TokenUser   = "user"
	TokenDevice = "device"
)

type User struct {
	UID              string     `json:"uid" bson:"-"`
	CreatedAt        time.Time  `json:"created" bson:"-"`
	ModifiedAt       time.Time  `json:"modified" bson:"-"`
	Login            string     `json:"login"`
	Password         string     `json:"password,omitempty"`
	Salt             string     `json:"-"`
	Email            string     `json:"email,omitempty"`
	Name             string     `json:"name"`
	Activated        bool       `json:"activated"`
	Role             string     `json:"role"`
	ActivationKey    *string    `json:"-"`
	ActivationExpiry *time.Time `json:"-"`
}

const (
	RoleUser  = "member"
	RoleAdmin = "admin"
)

const LENGTH_TOKEN = 32

const HEADER_APPSIGN = "X-Application"
