package gogrove_test

import (
	"fmt"
	"log"
	"time"

	"github.com/benmcclelland/gogrove"
)

func ExampleGetFirmwareVersion() {
	s, err := gogrove.New()
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	vers, err := s.GetFirmwareVersion()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(vers)
}
func ExampleSetPortMode() {
	s, err := gogrove.New()
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	// Set port D6 to Input
	err = s.SetPortMode(gogrove.PortD6, gogrove.ModeInput)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleBlinkLedD6() {
	s, err := gogrove.New()
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	// errors ignored
	s.SetPortMode(gogrove.PortD6, gogrove.ModeInput)
	for i := 0; i < 5; i++ {
		s.TurnOn(gogrove.PortD6)
		time.Sleep(time.Second)
		s.TurnOff(gogrove.PortD6)
		time.Sleep(time.Second)
	}
}

func ExampleLCDGreenTxt() {
	l, err := gogrove.NewLCD()
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	l.ClearText()
	l.SetText("Hello,\nGo Example")
	l.SetRGB(0, 255, 0)
}

func ExampleLCDOffClear() {
	l, err := gogrove.NewLCD()
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	l.ClearText()
	l.SetRGB(0, 0, 0)
}
