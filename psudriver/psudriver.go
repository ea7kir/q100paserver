/*
 *  Q-100 PA Server
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package psudriver

import (
	"q100paserver/logger"
	"time"

	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
)

// using: https://github.com/warthog618/gpiod

const (
	// WAVESHARE RPi RELAY BOARD
	RELAY_28v_GPIO = 26 // pin 37 GPIO_26 (CH1 P25)
	RELAY_12v_GPIO = 20 // pin 38 GPIO_20 (CH2 P28)
	RELAY_5v_GPIO  = 21 // pin 40 GPIO_21 (CH3 P29)
	// NOTE: the opto coupleers need reverse logic
	RELAY_ON  = 0
	RELAY_OFF = 1
)

var (
	relay28v *gpiod.Line
	relay12v *gpiod.Line
	relay5v  *gpiod.Line
	isUp     bool
)

func Configure() {
	relay28vLine, err := gpiod.RequestLine("gpiochip0", rpi.J8p37, gpiod.AsOutput(0))
	if err != nil {
		logger.Fatal.Fatalf("Failed to configure 28v rpi.J8p37 : %s", err)
	}
	relay28v = relay28vLine
	relay28v.SetValue(RELAY_OFF)
	relay12vLine, err := gpiod.RequestLine("gpiochip0", rpi.J8p38, gpiod.AsOutput(0))
	if err != nil {
		logger.Fatal.Fatalf("Failed to configure 12v rpi.J8p38 : %s", err)
	}
	relay12v = relay12vLine
	relay12v.SetValue(RELAY_OFF)
	relay5vLine, err := gpiod.RequestLine("gpiochip0", rpi.J8p40, gpiod.AsOutput(0))
	if err != nil {
		logger.Fatal.Fatalf("Failed to configure 5v rpi.J8p40 : %s", err)
	}
	relay5v = relay5vLine
	relay5v.SetValue(RELAY_OFF)
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
	logger.Info.Printf("Power UP is starting...")
	relay5v.SetValue(RELAY_ON)
	time.Sleep(200 * time.Millisecond)
	relay28v.SetValue(RELAY_ON)
	time.Sleep(200 * time.Millisecond)
	relay12v.SetValue(RELAY_ON)
	isUp = true
	logger.Info.Printf("Power UP has completed\n")
}

func Down() {
	logger.Info.Printf("Power DOWN is starting...\n")
	relay28v.SetValue(RELAY_OFF)
	time.Sleep(200 * time.Millisecond)
	relay5v.SetValue(RELAY_OFF)
	time.Sleep(200 * time.Millisecond)
	relay12v.SetValue(RELAY_OFF)
	isUp = false
	logger.Info.Printf("Power DOWN has completed\n")
}
