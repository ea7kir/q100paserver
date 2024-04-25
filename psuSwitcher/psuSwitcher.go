/*
 *  Q-100 PA Server
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package psuSwitcher

import (
	"os"
	"time"

	"github.com/ea7kir/qLog"
	"github.com/warthog618/go-gpiocdev"
	"github.com/warthog618/go-gpiocdev/device/rpi"
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
	relay28v *gpiocdev.Line
	relay12v *gpiocdev.Line
	relay5v  *gpiocdev.Line
	isUp     bool
)

func Configure() {
	relay28vLine, err := gpiocdev.RequestLine("gpiochip0", k28vRelayPin, gpiocdev.AsOutput(0))
	if err != nil {
		qLog.Fatal("Failed to configure 28v rpi.J8p37 : %s", err)
		os.Exit(1)
	}
	relay28v = relay28vLine
	relay28v.SetValue(kRelayOff)
	relay12vLine, err := gpiocdev.RequestLine("gpiochip0", k12vRelayPin, gpiocdev.AsOutput(0))
	if err != nil {
		qLog.Fatal("Failed to configure 12v rpi.J8p38 : %s", err)
		os.Exit(1)
	}
	relay12v = relay12vLine
	relay12v.SetValue(kRelayOff)
	relay5vLine, err := gpiocdev.RequestLine("gpiochip0", k5vRelayPin, gpiocdev.AsOutput(0))
	if err != nil {
		qLog.Fatal("Failed to configure 5v rpi.J8p40 : %s", err)
		os.Exit(1)
	}
	relay5v = relay5vLine
	relay5v.SetValue(kRelayOff)
}

func Shutdown() {
	// revert lines to input on the way out
	if isUp {
		Down()
	}
	relay28v.Reconfigure(gpiocdev.AsInput)
	relay28v.Close()
	relay12v.Reconfigure(gpiocdev.AsInput)
	relay12v.Close()
	relay5v.Reconfigure(gpiocdev.AsInput)
	relay5v.Close()
}

func Up() {
	qLog.Info("Power UP is starting...")
	relay5v.SetValue(kRelayOn)
	time.Sleep(200 * time.Millisecond)
	relay28v.SetValue(kRelayOn)
	time.Sleep(200 * time.Millisecond)
	relay12v.SetValue(kRelayOn)
	isUp = true
	qLog.Info("Power UP has completed\n")
}

func Down() {
	qLog.Info("Power DOWN is starting...\n")
	relay28v.SetValue(kRelayOff)
	time.Sleep(200 * time.Millisecond)
	relay5v.SetValue(kRelayOff)
	time.Sleep(200 * time.Millisecond)
	relay12v.SetValue(kRelayOff)
	isUp = false
	qLog.Info("Power DOWN has completed\n")
}
