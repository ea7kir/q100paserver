/*
 *  Q-100 PA Server
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package psudriver

import (
	"q100paserver/mylogger"
	"time"

	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
)

// using: https://github.com/warthog618/gpiod

const (
	// WAVESHARE RPi RELAY BOARD

	k28vRelayPin = rpi.J8p37 // pin 37 GPIO_26 (CH1 P25)
	k12vRelayPin = rpi.J8p38 // pin 38 GPIO_20 (CH2 P28)
	k5vRelayPin  = rpi.J8p40 // pin 40 GPIO_21 (CH3 P29)

	// NOTE: the opto coupleers need reverse logic

	kRelayOn  = 0
	kRelayOff = 1
)

var (
	relay28v *gpiod.Line
	relay12v *gpiod.Line
	relay5v  *gpiod.Line
	isUp     bool
)

func Configure() {
	relay28vLine, err := gpiod.RequestLine("gpiochip0", k28vRelayPin, gpiod.AsOutput(0))
	if err != nil {
		mylogger.Fatal.Fatalf("Failed to configure 28v rpi.J8p37 : %s", err)
	}
	relay28v = relay28vLine
	relay28v.SetValue(kRelayOff)
	relay12vLine, err := gpiod.RequestLine("gpiochip0", k12vRelayPin, gpiod.AsOutput(0))
	if err != nil {
		mylogger.Fatal.Fatalf("Failed to configure 12v rpi.J8p38 : %s", err)
	}
	relay12v = relay12vLine
	relay12v.SetValue(kRelayOff)
	relay5vLine, err := gpiod.RequestLine("gpiochip0", k5vRelayPin, gpiod.AsOutput(0))
	if err != nil {
		mylogger.Fatal.Fatalf("Failed to configure 5v rpi.J8p40 : %s", err)
	}
	relay5v = relay5vLine
	relay5v.SetValue(kRelayOff)
}

func Shutdown() {
	// revert lines to input on the way out
	if isUp {
		Down()
	}
	relay28v.Reconfigure(gpiod.AsInput)
	relay28v.Close()
	relay12v.Reconfigure(gpiod.AsInput)
	relay12v.Close()
	relay5v.Reconfigure(gpiod.AsInput)
	relay5v.Close()
}

func Up() {
	mylogger.Info.Printf("Power UP is starting...")
	relay5v.SetValue(kRelayOn)
	time.Sleep(200 * time.Millisecond)
	relay28v.SetValue(kRelayOn)
	time.Sleep(200 * time.Millisecond)
	relay12v.SetValue(kRelayOn)
	isUp = true
	mylogger.Info.Printf("Power UP has completed\n")
}

func Down() {
	mylogger.Info.Printf("Power DOWN is starting...\n")
	relay28v.SetValue(kRelayOff)
	time.Sleep(200 * time.Millisecond)
	relay5v.SetValue(kRelayOff)
	time.Sleep(200 * time.Millisecond)
	relay12v.SetValue(kRelayOff)
	isUp = false
	mylogger.Info.Printf("Power DOWN has completed\n")
}
