/*
 *  Q-100 PA Server
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package power

import (
	"q100paserver/logger"
	"time"

	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
)

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
	up       bool
)

func Configure(pi int) {
	relay28vLine, err := gpiod.RequestLine("gpiochip0", rpi.J8p37, gpiod.AsOutput(0))
	if err != nil {
		panic(err)
	}
	relay28v = relay28vLine
	relay12vLine, err := gpiod.RequestLine("gpiochip0", rpi.J8p38, gpiod.AsOutput(0))
	if err != nil {
		panic(err)
	}
	relay12v = relay12vLine
	relay5vLine, err := gpiod.RequestLine("gpiochip0", rpi.J8p40, gpiod.AsOutput(0))
	if err != nil {
		panic(err)
	}
	relay5v = relay5vLine
}

func Shutdown() {
	if up {
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
	time.Sleep(time.Second)
	relay28v.SetValue(RELAY_ON)
	time.Sleep(time.Second)
	relay12v.SetValue(RELAY_ON)
	up = true
	time.Sleep(time.Second)
	logger.Info.Printf("Power UP has completed\n")
}

func Down() {
	logger.Info.Printf("Power DOWN is starting...\n")
	relay28v.SetValue(RELAY_OFF)
	time.Sleep(time.Second)
	relay5v.SetValue(RELAY_OFF)
	time.Sleep(time.Second)
	relay12v.SetValue(RELAY_OFF)
	up = false
	logger.Info.Printf("Power DOWN has completed\n")
}
