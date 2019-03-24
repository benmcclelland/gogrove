package gogrove

import (
	"context"
	"fmt"
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

// SetText clears display text, and sets text, accepts newlines and
// auto wraps long line.
func (l *LCD) SetText(str string) error {
	l.Lock()
	defer l.Unlock()

	// clear
	err := l.b.Tx(displayTextAddr, []byte{textCmd, 0x01}, nil)
	if err != nil {
		return err
	}

	// display on, no cursor
	err = l.b.Tx(displayTextAddr, []byte{textCmd, 0x08 | 0x04}, nil)
	if err != nil {
		return err
	}

	// 2 lines
	err = l.b.Tx(displayTextAddr, []byte{textCmd, 0x28}, nil)
	if err != nil {
		return err
	}

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

// SetText2 clears display text, and sets text, accepts 2 lines
// truncating after 16 chars
func (l *LCD) SetText2(line1, line2 string) error {
	l.Lock()
	defer l.Unlock()

	// clear
	err := l.b.Tx(displayTextAddr, []byte{textCmd, 0x01}, nil)
	if err != nil {
		return err
	}
	// display on, no cursor
	err = l.b.Tx(displayTextAddr, []byte{textCmd, 0x08 | 0x04}, nil)
	if err != nil {
		return err
	}
	// 2 lines
	err = l.b.Tx(displayTextAddr, []byte{textCmd, 0x28}, nil)
	if err != nil {
		return err
	}

	err = l.displayLine(line1)
	if err != nil {
		return err
	}

	err = l.b.Tx(displayTextAddr, []byte{textCmd, 0xc0}, nil)
	if err != nil {
		return err
	}

	err = l.displayLine(line2)
	if err != nil {
		return err
	}

	return nil
}

func (l *LCD) displayLine(str string) error {
	count := 0
	for _, c := range str {
		if c == '\n' || count == 16 {
			break
		}
		count++
		err := l.b.Tx(displayTextAddr, []byte{0x40, uint8(c)}, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

// ScrollText clears display text, and sets text, accepts 2 lines
// scrolling both lines if longer than 16.  Pass a context to cancel to
// terminate function, otherwise scolls forever.
func (l *LCD) ScrollText(ctx context.Context, line1, line2 string) error {
	l.Lock()
	defer l.Unlock()

	i := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		default:

			// clear
			err := l.b.Tx(displayTextAddr, []byte{textCmd, 0x01}, nil)
			if err != nil {
				return err
			}

			// display on, no cursor
			err = l.b.Tx(displayTextAddr, []byte{textCmd, 0x08 | 0x04}, nil)
			if err != nil {
				return err
			}
			// 2 lines
			err = l.b.Tx(displayTextAddr, []byte{textCmd, 0x28}, nil)
			if err != nil {
				return err
			}

			stbegin := 0
			frame := 16
			if len(line1) > 16 {
				stbegin = i % len(line1)
			}
			if len(line1)-stbegin < 16 {
				frame = len(line1) - stbegin
			}
			fmt.Println(line1[stbegin:][:frame])
			err = l.displayLine(line1[stbegin:][:frame])
			if err != nil {
				return err
			}

			err = l.b.Tx(displayTextAddr, []byte{textCmd, 0xc0}, nil)
			if err != nil {
				return err
			}

			stbegin = 0
			frame = 16
			if len(line2) > 16 {
				stbegin = i % len(line2)
			}
			if len(line2)-stbegin < 16 {
				frame = len(line2) - stbegin
			}
			fmt.Println(line2[stbegin:][:frame])
			err = l.displayLine(line2[stbegin:][:frame])
			if err != nil {
				return err
			}
			fmt.Println()

			time.Sleep(500 * time.Millisecond)
			i++
		}
	}
}
