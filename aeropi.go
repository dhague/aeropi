package main

import (
	"fmt"
	"time"

	log "gopkg.in/inconshreveable/log15.v2"

	"github.com/davecheney/i2c"
	"github.com/dhague/atmosphere"
	"github.com/quinte17/bme280"
)

func main() {

	log.Debug("Starting up...")

	dev, err := i2c.New(0x77, 1)
	if err != nil {
		log.Error("Error: %v", err)
	}
	bme, err := bme280.NewI2CDriver(dev)
	if err != nil {
		log.Error("Error: %v", err)
	}

	for {
		weather, err := bme.Readenv()

		density := atmosphere.AirDensity(atmosphere.TemperatureC(weather.Temp),
			atmosphere.Pressure(weather.Press*100),
			atmosphere.RelativeHumidity(weather.Hum))
		if err != nil {
			log.Error("Error: %v", err)
		} else {
			log.Debug(fmt.Sprintf("Temperature (degC) %2.2f, pressure (mbar) %4.2f, humidity %2.2f, Density %2.3f kg/m3, CorrFact %1.4f",
				weather.Temp, weather.Press, weather.Hum, float64(density), float64(density*0.8)))
		}

		<-time.After(5 * time.Second)
	}
}
