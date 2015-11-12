package lcd16x2keypad

import (
	"fmt"
)

func Show(lines LcdArray) {
	// TODO: This should actually write the data in the LcdArray out to the device via I2C
	fmt.Println(lines.first())
	fmt.Println(lines.second())
}
