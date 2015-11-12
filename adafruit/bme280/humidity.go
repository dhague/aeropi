package bme280

import (

)

// Humidity is stored as a float32 representing relative humidity
type Humidity float32

func (humidity Humidity) percent() float32 {
	return float32(humidity)
}
