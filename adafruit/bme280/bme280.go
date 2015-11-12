package bme280

import ()

func Read() (Temperature, Pressure, Humidity) {
	// TODO: This should actually read from the BME280 using I2C
	return 20.0, 1000.0, 75.0
}
