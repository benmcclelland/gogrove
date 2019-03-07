package gogrove_test

import (
	"fmt"
	"log"
	"time"

	"github.com/benmcclelland/gogrove"
)

func ExampleSession_GetFirmwareVersion() {
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
func ExampleSession_SetPortMode() {
	s, err := gogrove.New()
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	// Set port D3 to Input
	err = s.SetPortMode(gogrove.PortD3, gogrove.ModeInput)
	if err != nil {
		log.Fatal(err)
	}
}

func Example() {
	s, err := gogrove.New()
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	// blink LED on D6, errors ignored
	for i := 0; i < 5; i++ {
		s.TurnOn(gogrove.PortD6)
		time.Sleep(time.Second)
		s.TurnOff(gogrove.PortD6)
		time.Sleep(time.Second)
	}
}

func Example_greenTxtLCD() {
	l, err := gogrove.NewLCD()
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	l.ClearText()
	l.SetText("Hello,\nGo Example")
	l.SetRGB(0, 255, 0)
}

func Example_offClearLCD() {
	l, err := gogrove.NewLCD()
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	l.ClearText()
	l.SetRGB(0, 0, 0)
}
