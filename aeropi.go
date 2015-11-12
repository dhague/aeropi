package main 

import (
	"fmt"
	"github.com/dhague/aeropi/adafruit/lcd16x2keypad"
	"github.com/dhague/aeropi/adafruit/bme280"
)

func main() {
	temp := bme280.Temperature(21)
	pressure := bme280.Pressure(1004)
	humidity := bme280.Humidity(80)
	
	line1 := lcd16x2keypad.Line(fmt.Sprintf("%v\xb0C %v hPa", temp, pressure))
	line2 := lcd16x2keypad.Line(fmt.Sprintf("%v%% humidity", humidity))
	lcd16x2keypad.Display(line1, line2)
}

