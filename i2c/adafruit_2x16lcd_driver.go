package i2c

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/i2c"
	"time"
)

const (
	REG_RED   = 0x04
	REG_GREEN = 0x03
	REG_BLUE  = 0x02

	LCD_CLEARDISPLAY        = 0x01
	LCD_RETURNHOME          = 0x02
	LCD_ENTRYMODESET        = 0x04
	LCD_DISPLAYCONTROL      = 0x08
	LCD_CURSORSHIFT         = 0x10
	LCD_FUNCTIONSET         = 0x20
	LCD_SETCGRAMADDR        = 0x40
	LCD_SETDDRAMADDR        = 0x80
	LCD_ENTRYRIGHT          = 0x00
	LCD_ENTRYLEFT           = 0x02
	LCD_ENTRYSHIFTINCREMENT = 0x01
	LCD_ENTRYSHIFTDECREMENT = 0x00
	LCD_DISPLAYON           = 0x04
	LCD_DISPLAYOFF          = 0x00
	LCD_CURSORON            = 0x02
	LCD_CURSOROFF           = 0x00
	LCD_BLINKON             = 0x01
	LCD_BLINKOFF            = 0x00
	LCD_DISPLAYMOVE         = 0x08
	LCD_CURSORMOVE          = 0x00
	LCD_MOVERIGHT           = 0x04
	LCD_MOVELEFT            = 0x00
	LCD_2LINE               = 0x08
	LCD_CMD                 = 0x80
	LCD_DATA                = 0x40

	LCD_2NDLINEOFFSET = 0x40
)

var _ gobot.Driver = (*Adafruit2x16LcdDriver)(nil)

// Adafruit2x16LcdDriver is a driver for the Jhd1313m1 LCD display which has two i2c addreses,
// one belongs to a controller and the other controls solely the backlight.
// This module was tested with the Seed Grove LCD RGB Backlight v2.0 display which requires 5V to operate.
// http://www.seeedstudio.com/wiki/Grove_-_LCD_RGB_Backlight
type Adafruit2x16LcdDriver struct {
	name       string
	connection i2c.I2c
	lcdAddress int
	rgbAddress int
}

// NewAdafruit2x16LcdDriver creates a new driver with specified name and i2c interface.
func NewAdafruit2x16LcdDriver(a i2c.I2c, name string) *Adafruit2x16LcdDriver {
	return &Adafruit2x16LcdDriver{
		name:       name,
		connection: a,
		lcdAddress: 0x3E,
		rgbAddress: 0x62,
	}
}

// Name returns the name the JHD1313M1 Driver was given when created.
func (h *Adafruit2x16LcdDriver) Name() string { return h.name }

// Connection returns the driver connection to the device.
func (h *Adafruit2x16LcdDriver) Connection() gobot.Connection {
	return h.connection.(gobot.Connection)
}

// Start starts the backlit and the screen and initializes the states.
func (h *Adafruit2x16LcdDriver) Start() []error {
	if err := h.connection.I2cStart(h.lcdAddress); err != nil {
		return []error{err}
	}

	if err := h.connection.I2cStart(h.rgbAddress); err != nil {
		return []error{err}
	}

	<-time.After(50000 * time.Microsecond)
	payload := []byte{LCD_CMD, LCD_FUNCTIONSET | LCD_2LINE}
	if err := h.connection.I2cWrite(h.lcdAddress, payload); err != nil {
		if err := h.connection.I2cWrite(h.lcdAddress, payload); err != nil {
			return []error{err}
		}
	}

	<-time.After(100 * time.Microsecond)
	if err := h.connection.I2cWrite(h.lcdAddress, []byte{LCD_CMD, LCD_DISPLAYCONTROL | LCD_DISPLAYON}); err != nil {
		return []error{err}
	}

	<-time.After(100 * time.Microsecond)
	if err := h.Clear(); err != nil {
		return []error{err}
	}

	if err := h.connection.I2cWrite(h.lcdAddress, []byte{LCD_CMD, LCD_ENTRYMODESET | LCD_ENTRYLEFT | LCD_ENTRYSHIFTDECREMENT}); err != nil {
		return []error{err}
	}

	if err := h.setReg(0, 0); err != nil {
		return []error{err}
	}
	if err := h.setReg(1, 0); err != nil {
		return []error{err}
	}
	if err := h.setReg(0x08, 0xAA); err != nil {
		return []error{err}
	}

	if err := h.SetRGB(255, 255, 255); err != nil {
		return []error{err}
	}

	return nil
}

// SetRGB sets the Red Green Blue value of backlit.
func (h *Adafruit2x16LcdDriver) SetRGB(r, g, b int) error {
	if err := h.setReg(REG_RED, r); err != nil {
		return err
	}
	if err := h.setReg(REG_GREEN, g); err != nil {
		return err
	}
	return h.setReg(REG_BLUE, b)
}

// Clear clears the text on the lCD display.
func (h *Adafruit2x16LcdDriver) Clear() error {
	err := h.command([]byte{LCD_CLEARDISPLAY})
	return err
}

// Home sets the cursor to the origin position on the display.
func (h *Adafruit2x16LcdDriver) Home() error {
	err := h.command([]byte{LCD_RETURNHOME})
	// This wait fixes a race condition when calling home and clear back to back.
	<-time.After(2 * time.Millisecond)
	return err
}

// Write displays the passed message on the screen.
func (h *Adafruit2x16LcdDriver) Write(message string) error {
	// This wait fixes an odd bug where the clear function doesn't always work properly.
	<-time.After(1 * time.Millisecond)
	for _, val := range message {
		if val == '\n' {
			if err := h.SetPosition(16); err != nil {
				return err
			}
			continue
		}
		if err := h.connection.I2cWrite(h.lcdAddress, []byte{LCD_DATA, byte(val)}); err != nil {
			return err
		}
	}
	return nil
}

// SetPosition sets the cursor and the data display to pos.
// 0..15 are the positions in the first display line.
// 16..32 are the positions in the second display line.
func (h *Adafruit2x16LcdDriver) SetPosition(pos int) (err error) {
	if pos < 0 || pos > 31 {
		err = i2c.ErrInvalidPosition
		return
	}
	offset := byte(pos)
	if pos >= 16 {
		offset -= 16
		offset |= LCD_2NDLINEOFFSET
	}
	err = h.command([]byte{LCD_SETDDRAMADDR | offset})
	return
}

func (h *Adafruit2x16LcdDriver) Scroll(leftToRight bool) error {
	if leftToRight {
		return h.connection.I2cWrite(h.lcdAddress, []byte{LCD_CMD, LCD_CURSORSHIFT | LCD_DISPLAYMOVE | LCD_MOVELEFT})
	}

	return h.connection.I2cWrite(h.lcdAddress, []byte{LCD_CMD, LCD_CURSORSHIFT | LCD_DISPLAYMOVE | LCD_MOVERIGHT})
}

// Halt is a noop function.
func (h *Adafruit2x16LcdDriver) Halt() []error { return nil }

func (h *Adafruit2x16LcdDriver) setReg(command int, data int) error {
	if err := h.connection.I2cWrite(h.rgbAddress, []byte{byte(command), byte(data)}); err != nil {
		return err
	}
	return nil
}

func (h *Adafruit2x16LcdDriver) command(buf []byte) error {
	return h.connection.I2cWrite(h.lcdAddress, append([]byte{LCD_CMD}, buf...))
}
