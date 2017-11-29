package types

type Device struct {
	Title          string
	Description    string
	Location       string
	FeaturedSensor *Sensor
	Tags           []string
}

type Sensor struct {
	Alias       string
	Title       string
	Type        int
	Description string
}

const (
	SensorTemperature = 0
	SensorHumidity    = 1
	SensorWindSpeed   = 2
)
