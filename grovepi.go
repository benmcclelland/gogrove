package gogrove

import (
	"encoding/binary"
	"fmt"
	"math"
	"sync"
	"time"

	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
)

const (
	digitalRead  uint8 = 1
	digitalWrite uint8 = 2
	analogRead   uint8 = 3
	analogWrite  uint8 = 4
	pinMode      uint8 = 5
	// dustSensorReadInt uint8 = 6
	ultraSonic      uint8 = 7
	firmwareVersion uint8 = 8
	// dustSensorInt uint8 = 9
	// dustSensorRead uint8 = 10
	// encoderRead uint8 = 11
	// flowRead uint8 = 12
	// flowDis uint8 = 13
	// dustSensorEn uint8 = 14
	// dustSensorDis uint8 = 15
	// encoderEn uint8 =  16
	// encoderDis uint8 = 17
	// flowEn uint8 = 18
	// irRead uint8 = 21
	// irRecvPin uint8 = 22
	dataNotAvailable uint8 = 23
	// irReadIsData uint8 = 24
	dhtTemp uint8 = 40
	// ledBarInit uint8 = 50
	// ledBarSetGreenToRed uint8 = 51
	// ledBarSetLevel uint8 = 52
	// ledBarSetLed uint8 = 53
	// ledBarToggelLed uint8 = 54
	// ledBarSetBits uint8 = 55
	// ledBarGetBits uint8 = 56
	// fourDigitInit uint8 = 70
	// fourDigitSetBright uint8 = 71
	// fourDigitRAn0s uint8 = 72
	// fourDigitRAw0s uint8 = 73
	// fourDigitSetDigit uint8 = 74
	// fourDigitSetSegment uint8 = 75
	// fourDigitSetValsWithColon uint8 = 76
	// fourDigitDisAReadNSec uint8 = 77
	// fourDigitDispOn uint8 = 78
	// fourDigitDispOff uint8 = 79
	// chainLedStorColor uint8 = 90
	// chainLedInit uint8 = 91
	// chainLedInitWithColor uint8 = 92
	// chainLedSetLedsStorPattern uint8 = 93
	// chainLedSetLedsStorModulo uint8 = 94
	// chainLedSetBarGraph uint8 = 95

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

	// BlueDHTSensor is the DHT sensor that comes with base kit (DHT11)
	BlueDHTSensor uint8 = 0
	// WhiteDHTSensor is the separate white DHT sensor (DHT22)
	WhiteDHTSensor uint8 = 1
	// DHT21Sensor DHT21
	DHT21Sensor uint8 = 2
	// AM2301Sensor AM2301
	AM2301Sensor uint8 = 3
)

// Session holds session info for interacting with GrovePi.
type Session struct {
	sync.Mutex
	d *i2c.Dev
	b i2c.BusCloser
}

// New initializes new session with default GrovePi address.
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

// NewWithAddress initializes new session with given GrovePi address.
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

// Close closes GrovePi session.
func (s *Session) Close() error { return s.b.Close() }

// SetPortMode sets port to mode, for example:
// SetPortMode(gogrove.PortA0, gogrove.ModeOutput),
// SetPortMode(gogrove.PortD3, gogrove.ModeInput).
func (s *Session) SetPortMode(port, mode uint8) error {
	s.Lock()
	defer s.Unlock()

	write := []byte{0, pinMode, port, mode, 0}
	return s.d.Tx(write, nil)
}

// GetFirmwareVersion returns the GrovePi firmware version.
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

// DigitalRead returns the value from a digital port.
// On success, this will be either 0 or 1.
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

// IsOn is shorthand for DigitalRead on a digital port,
// returning true if the port is 1
// and false if the port is 0.
// This ignores errors from DigitalRead for easier inlining.
func (s *Session) IsOn(port uint8) bool {
	state, _ := s.DigitalRead(port)
	if state == 0 {
		return false
	}
	return true
}

// DigitalWrite sets the value for the given port.
// The value must be 0 or 1.
func (s *Session) DigitalWrite(port, value uint8) error {
	s.Lock()
	defer s.Unlock()

	if value != 0 && value != 1 {
		return fmt.Errorf("invalid digital write value")
	}

	write := []byte{0, digitalWrite, port, value, 0}
	return s.d.Tx(write, nil)
}

// TurnOn is shorthand for DigitalWrite(port, 1).
func (s *Session) TurnOn(port uint8) error {
	return s.DigitalWrite(port, 1)
}

// TurnOff is shorthand for DigitalWrite(port, 0).
func (s *Session) TurnOff(port uint8) error {
	return s.DigitalWrite(port, 0)
}

// AnalogRead reads analog value from port.
// This is only valid for PortA0, PortA1, or PortA2.
// The returned value will be between 0-1023 inclusive.
func (s *Session) AnalogRead(port uint8) (uint16, error) {
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
		if read[0] != dataNotAvailable || i == 100 {
			break
		}
		time.Sleep(1 * time.Millisecond)
		if err := s.d.Tx(nil, read); err != nil {
			return 0, err
		}
		i++
	}
	if read[0] != analogRead {
		return 0, fmt.Errorf("command error response: %v", read[0])
	}

	return binary.BigEndian.Uint16(read[1:][:2]), nil
}

// AnalogWrite writes value 0-255 inclusive to given port.
// This appears to only be valid for PortD3, PortD5, and PortD6
// using PWM write.
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

// ReadDHT returns temp (C), humidity (%).
// Must pass the sensort type,
// one of: gogrove.BlueDHTSensor gogrove.WhiteDHTSensor.
func (s *Session) ReadDHT(port, sensor uint8) (float32, float32, error) {
	s.Lock()
	defer s.Unlock()

	write := []byte{0, dhtTemp, port, sensor, 0}
	read := make([]byte, 9)
	if err := s.d.Tx(write, read); err != nil {
		return 0, 0, err
	}
	i := 0
	for {
		if read[0] == dhtTemp || i == 100 {
			break
		}
		time.Sleep(10 * time.Millisecond)
		if err := s.d.Tx(nil, read); err != nil {
			return 0, 0, err
		}
		i++
	}

	if read[0] != dhtTemp {
		return 0, 0, fmt.Errorf("invalid command response: %v", read[0])
	}

	temp := float32frombytes(read[1:][:4])
	humidity := float32frombytes(read[5:][:4])

	if temp > -100.0 && temp < 150.0 && humidity >= 0.0 && humidity <= 100.0 {
		return temp, humidity, nil
	}

	return 0, 0, fmt.Errorf("value out of bounds")
}

func float32frombytes(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

// ReadUltraSonic returns distance in cm.
// Sensor spec: measuring range 2-350cm, resolution 1cm.
func (s *Session) ReadUltraSonic(port uint8) (uint16, error) {
	s.Lock()
	defer s.Unlock()

	write := []byte{0, ultraSonic, port, 0, 0}
	read := make([]byte, 3)
	if err := s.d.Tx(write, read); err != nil {
		return 0, err
	}
	i := 0
	for {
		if read[0] == ultraSonic || i == 100 {
			break
		}
		time.Sleep(10 * time.Millisecond)
		if err := s.d.Tx(nil, read); err != nil {
			return 0, err
		}
		i++
	}

	if read[0] != ultraSonic {
		return 0, fmt.Errorf("invalid command response: %v", read[0])
	}

	return binary.BigEndian.Uint16(read[1:][:2]), nil
}
