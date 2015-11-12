package bme280

import (

)

// Pressure is stored as a float32 representing hPa
type Pressure float32

func (pressure Pressure) Hectopascals() float32 {
	return float32(pressure)
}

func (pressure Pressure) Millibars() float32 {
	return float32(pressure)
}

func (pressure Pressure) MmHg() float32 {
	return float32(pressure) * 0.750061561303
}
