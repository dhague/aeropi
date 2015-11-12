package bme280

import (

)

// Temperature is stored as a float32 representing degrees celsius 
type Temperature float32

func (temp Temperature) Celsius() float32 {
	return float32(temp)
}

func (temp Temperature) Fahrenheit() float32 {
	return 32 + (float32(temp) * 9.0 / 5.0)
}

func (temp Temperature) Kelvin() float32 {
	return float32(temp) + 273.15
}
