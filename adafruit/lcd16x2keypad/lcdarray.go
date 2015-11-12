package lcd16x2keypad

import (

)

type LcdArray [2][16] byte

func (disp *LcdArray) SetLine(index int, content string) {
    copy(disp[index][:], content)
}

func (disp LcdArray) first() string {
	return string(disp[0][:])
}

func (disp LcdArray) second() string {
	return string(disp[1][:])
}
