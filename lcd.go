package gogrove

import (
	"sync"
	"time"

	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
)

const (
	displayRGBAddr  uint16 = 0x62
	displayTextAddr uint16 = 0x3e

	textCmd uint8 = 0x80
)

// LCD holds session info for interacting with GrovePi LCD.
type LCD struct {
	sync.Mutex
	b i2c.BusCloser
}

// NewLCD initializes a new session with the LCD.
func NewLCD() (*LCD, error) {
	if _, err := host.Init(); err != nil {
		return nil, err
	}

	b, err := i2creg.Open("")
	if err != nil {
		return nil, err
	}

	return &LCD{b: b}, nil
}

// Close closes session with LCD.
func (l *LCD) Close() error { return l.b.Close() }

// SetRGB sets the background RGB LCD color.
func (l *LCD) SetRGB(r, g, b uint8) error {
	l.Lock()
	defer l.Unlock()

	err := l.b.Tx(displayRGBAddr, []byte{0, 0}, nil)
	if err != nil {
		return err
	}
	err = l.b.Tx(displayRGBAddr, []byte{0x1, 0}, nil)
	if err != nil {
		return err
	}
	err = l.b.Tx(displayRGBAddr, []byte{0x08, 0xaa}, nil)
	if err != nil {
		return err
	}
	err = l.b.Tx(displayRGBAddr, []byte{0x4, r}, nil)
	if err != nil {
		return err
	}
	err = l.b.Tx(displayRGBAddr, []byte{0x3, g}, nil)
	if err != nil {
		return err
	}
	err = l.b.Tx(displayRGBAddr, []byte{0x2, b}, nil)
	if err != nil {
		return err
	}

	return nil
}

// ClearText clears the display text.
func (l *LCD) ClearText() error {
	l.Lock()
	defer l.Unlock()

	err := l.b.Tx(displayTextAddr, []byte{textCmd, 0x01}, nil)
	if err != nil {
		return err
	}

	return nil
}

// SetText clears display text, and sets text.
func (l *LCD) SetText(str string) error {
	l.Lock()
	defer l.Unlock()

	// clear
	err := l.b.Tx(displayTextAddr, []byte{textCmd, 0x01}, nil)
	if err != nil {
		return err
	}

	time.Sleep(50 * time.Millisecond)

	// display on, no cursor
	err = l.b.Tx(displayTextAddr, []byte{textCmd, 0x08 | 0x04}, nil)
	if err != nil {
		return err
	}

	time.Sleep(50 * time.Millisecond)

	// 2 lines
	err = l.b.Tx(displayTextAddr, []byte{textCmd, 0x28}, nil)
	if err != nil {
		return err
	}

	time.Sleep(50 * time.Millisecond)

	count := 0
	row := 0
	for _, c := range str {
		if c == '\n' || count == 16 {
			count = 0
			row++
			if row == 2 {
				break
			}
			err = l.b.Tx(displayTextAddr, []byte{textCmd, 0xc0}, nil)
			if err != nil {
				return err
			}
			if c == '\n' {
				continue
			}
		}
		count++
		err = l.b.Tx(displayTextAddr, []byte{0x40, uint8(c)}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
