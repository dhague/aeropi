package main

import (
	"fmt"
	"time"
	
	"github.com/dhague/aeropi/adafruit/bme280"
	"github.com/dhague/aeropi/adafruit/lcd16x2keypad"
	"github.com/dhague/aeropi/i2c"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
//	"github.com/hybridgroup/gobot/platforms/i2c"
	"github.com/hybridgroup/gobot/platforms/raspi"
)

func main() {
	temp, pressure, humidity := bme280.Read()

	lcdContent := lcd16x2keypad.LcdArray{}
	lcdContent.SetLine(0, fmt.Sprintf("%v\xb0C %v hPa", temp, pressure))
	lcdContent.SetLine(1, fmt.Sprintf("%v%% humidity", humidity))

	lcd16x2keypad.Show(lcdContent)

	gbot := gobot.NewGobot()

	r := raspi.NewRaspiAdaptor("raspi")
	led := gpio.NewLedDriver(r, "led", "7")
	
	//TODO - use these values - not compiling otherwise :-)
	sensor := i2c.NewBME280Driver(r, "sensor")
	display := i2c.NewAdafruit2x16LcdDriver(r, "display")

	work := func() {
		gobot.Every(1*time.Second, func() {
			led.Toggle()
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{r},
		[]gobot.Device{led},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
