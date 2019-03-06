package gogrove

import (
	"fmt"
	"sync"
	"time"

	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
)

const (
	digitalRead     uint8 = 1
	digitalWrite    uint8 = 2
	analogRead      uint8 = 3
	analogWrite     uint8 = 4
	pinMode         uint8 = 5
	firmwareVersion uint8 = 8

	dataValid        uint8 = 3
	dataNotAvailable uint8 = 23

	// PortA0 is the value for GrovePi A0
	PortA0 uint8 = 0
	// PortA1 is the value for GrovePi A1
	PortA1 uint8 = 1
	// PortA2 is the value for GrovePi A2
	PortA2 uint8 = 2
	// PortD2 is the value for GrovePi D2
	PortD2 uint8 = 2
	// PortD3 is the value for GrovePi D3
	PortD3 uint8 = 3
	// PortD4 is the value for GrovePi D4
	PortD4 uint8 = 4
	// PortD5 is the value for GrovePi D5
	PortD5 uint8 = 5
	// PortD6 is the value for GrovePi D6
	PortD6 uint8 = 6
	// PortD7 is the value for GrovePi D7
	PortD7 uint8 = 7
	// PortD8 is the value for GrovePi D8
	PortD8 uint8 = 8

	// ModeInput is used for SetPortMode to input
	ModeInput uint8 = 0
	// ModeOutput is used for SetPortMode to output
	ModeOutput uint8 = 1
)

// Session holds session info for interacting with GrovePi
type Session struct {
	sync.Mutex
	d *i2c.Dev
	b i2c.BusCloser
}

// New initializes new session with default GrovePi address
func New() (*Session, error) {
	if _, err := host.Init(); err != nil {
		return nil, err
	}

	b, err := i2creg.Open("")
	if err != nil {
		return nil, err
	}

	return &Session{
		d: &i2c.Dev{Addr: 0x4, Bus: b},
		b: b,
	}, nil
}

// NewWithAddress initializes new session with given GrovePi address
func NewWithAddress(address uint16) (*Session, error) {
	if _, err := host.Init(); err != nil {
		return nil, err
	}

	b, err := i2creg.Open("")
	if err != nil {
		return nil, err
	}

	return &Session{
		d: &i2c.Dev{Addr: address, Bus: b},
		b: b,
	}, nil
}

// Close closes GrovePi session
func (s *Session) Close() error { return s.b.Close() }

// SetPortMode sets port to mode, for example:
// SetPortMode(gogrove.PortA0, gogrove.ModeOutput)
// SetPortMode(gogrove.PortD3, gogrove.ModeInput)
func (s *Session) SetPortMode(port, mode uint8) error {
	s.Lock()
	defer s.Unlock()

	write := []byte{0, pinMode, port, mode, 0}
	return s.d.Tx(write, nil)
}

// GetFirmwareVersion returns the GrovePi firmware version
func (s *Session) GetFirmwareVersion() (string, error) {
	s.Lock()
	defer s.Unlock()

	write := []byte{0, firmwareVersion, 0, 0, 0}
	read := make([]byte, 5)
	if err := s.d.Tx(write, read); err != nil {
		return "", err
	}
	if read[0] != firmwareVersion {
		return "", fmt.Errorf("command error response: %v", read[0])
	}
	return fmt.Sprintf("%v.%v.%v", read[1], read[2], read[3]), nil
}

// DigitalRead return the value from a digital port
// on success, this will be either 0 or 1
func (s *Session) DigitalRead(port uint8) (uint8, error) {
	s.Lock()
	defer s.Unlock()

	write := []byte{0, digitalRead, port, 0, 0}
	read := make([]byte, 5)
	if err := s.d.Tx(write, read); err != nil {
		return 0, err
	}
	return read[1], nil
}

// IsOn is shorthand for DigitalRead on a digital port
// returning true if the port is 1
// and false if the port is 0
// this ignores errors from DigitalRead for easier inlining
func (s *Session) IsOn(port uint8) bool {
	state, _ := s.DigitalRead(port)
	if state == 0 {
		return false
	}
	return true
}

// DigitalWrite sets the value for the given port
// The value must be 0 or 1
func (s *Session) DigitalWrite(port, value uint8) error {
	s.Lock()
	defer s.Unlock()

	if value != 0 && value != 1 {
		return fmt.Errorf("invalid digital write value")
	}

	write := []byte{0, digitalWrite, port, value, 0}
	return s.d.Tx(write, nil)
}

// TurnOn is shorthand for DigitalWrite(port, 1)
func (s *Session) TurnOn(port uint8) error {
	return s.DigitalWrite(port, 1)
}

// TurnOff is shorthand for DigitalWrite(port, 0)
func (s *Session) TurnOff(port uint8) error {
	return s.DigitalWrite(port, 0)
}

// AnalogRead reads analog value from port
// this is only valid for PortA0, PortA1, or PortA2
// the returned value will be between 0-1023 inclusive
func (s *Session) AnalogRead(port uint8) (int, error) {
	s.Lock()
	defer s.Unlock()

	if !(port == PortA0 || port == PortA1 || port == PortA2) {
		return 0, fmt.Errorf("invalid port for analog read")
	}
	write := []byte{0, analogRead, port, 0, 0}
	read := make([]byte, 5)
	if err := s.d.Tx(write, read); err != nil {
		return 0, err
	}
	i := 0
	for {
		if read[0] != dataNotAvailable || i == 5 {
			break
		}
		time.Sleep(2 * time.Millisecond)
		if err := s.d.Tx(nil, read); err != nil {
			return 0, err
		}
		i++
	}
	// response of 3 seems to be "data valid"?
	if read[0] != dataValid {
		return 0, fmt.Errorf("command error response: %v", read[0])
	}

	return int(read[1])*256 + int(read[2]), nil
}

// AnalogWrite writes value 0-255 inclusive to given port
// This appears to only be valid for PortD3, PortD5, and PortD6
// using PWM write
func (s *Session) AnalogWrite(port, value uint8) error {
	s.Lock()
	defer s.Unlock()

	// D3, D5, D6 PWM writes
	if !(port == PortD3 || port == PortD5 || port == PortD6) {
		return fmt.Errorf("invalid port for analog write")
	}
	write := []byte{0, analogWrite, port, value, 0}
	return s.d.Tx(write, nil)
}
