/*
 *  Q-100 PA Server
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package ina226monitor

import (
	"sync"
	"time"
)

const (
	// INA226 current/voltage sensors
	// To discover I2C devices
	// $ sudo i2cdetect -y 1
	// Address could be 0x40, 0x41 or 0x42
	// Connect to pin3 GPIO2 SDA, pin 5 GPIO3 SCL
	// 4k7 pull-up resistors on data lines to 3.3v

	kFinalPaAddrees = 0x40
	kFinalPaShunt   = 0.002 // modified to 0.0021 from 0.002 to get correct current reading
	kFinalPaMaxAmps = 10
)

type (
	ina226Type struct {
		mu      sync.Mutex
		quit    chan bool
		volts   float64
		amps    float64
		address int8
		shunt   float64
		maxAmps float64
	}
)

var (
	finalPA *ina226Type
)

func newIna226(address int8, shunt float64, maxAmps float64) *ina226Type {
	return &ina226Type{
		mu:      sync.Mutex{},
		quit:    make(chan bool),
		volts:   0.0,
		amps:    0.0,
		address: address,
		shunt:   shunt,
		maxAmps: maxAmps,
	}
}

func Configure() {
	finalPA = newIna226(kFinalPaAddrees, kFinalPaShunt, kFinalPaMaxAmps)
	go readVoltsAmpsFor(finalPA)
}

func Shutdown() {
	finalPA.quit <- true
}

func PaVoltage() float64 {
	finalPA.mu.Lock()
	defer finalPA.mu.Unlock()
	return finalPA.volts
}

func PaCurrent() float64 {
	finalPA.mu.Lock()
	defer finalPA.mu.Unlock()
	return finalPA.amps
}

// Go routine to read voltage and current
func readVoltsAmpsFor(sensor *ina226Type) {
	var volts float64
	var amps float64
	for {
		select {
		case <-sensor.quit:
			return
		default:
		}
		volts = 0.0
		amps = 0.0

		// TODO: implementation goes here

		sensor.mu.Lock()
		sensor.volts = volts
		sensor.amps = amps
		sensor.mu.Unlock()
		time.Sleep(1 * time.Second)
	}
}
