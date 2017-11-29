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
	Type        string
	Description string
}

const (
	SensorTemperature = "Temperature"
	SensorHumidity    = "Humidity"
	SensorWindSpeed   = "WindSpeed"
)
