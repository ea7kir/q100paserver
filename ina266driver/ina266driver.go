/*
 *  Q-100 PA Server
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package ina266driver

// TODO: impleement an INA266 device

import (
	"sync"
	"time"
	// "github.com/warthog618/gpiod/device/rpi"
)

// TODO: install i2c-tools
// sudo apt install -y i2c-tools
// TODO: enable I2C in config

// There's an INA266 C library at https://github.com/MarioAriasGa/raspberry-pi-ina226

// INA226 current/voltage sensors
// To discover I2C devices
// $ sudo i2cdetect -y 1
// TODO: address could be 0x40, 0x41 or 0x42
// I2C pin3 GPIO2 SDA, pin 5 GPIO3 SCL
// 4k7 pull-up resistors on data lines to 3.3v

// TRY https://pkg.go.dev/github.com/Fede85/go-ina226@v0.0.0-20160609112003-36bc433ce086

const (
	kFinalPaAddrees = 0x40
	kFinalPaShunt   = 0.0021 // modified from 0,002 to get correct current reading
	kFinalPaMaxAmps = 10
	// kDefaultSDAPin  = rpi.J8p3
	// kDefaultSCLPin  = rpi.J8p5

	INA226_RESET                     uint16 = 0x8000 // UInt16
	INA226_REG_CALIBRATION           uint8  = 0x05   // UInt8
	INA226_REG_CONFIGURATION         uint8  = 0x00   // UInt8
	INA226_TIME_8MS                  uint8  = 7      // 8.244ms UInt8
	INA226_AVERAGES_16               uint8  = 2      // UInt8
	INA226_MODE_SHUNT_BUS_CONTINUOUS uint8  = 7      // UI
	INA226_REG_BUS_VOLTAGE           uint8  = 0x02   // UInt8
	INA226_REG_CURRENT               uint8  = 0x04   // UInt8
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

// func byteSwapped(w uint16) uint16 {
// 	return bits.ReverseBytes16(w)
// }

// const i2c_bus = 1

// func i2c_open(bus uint8, address uint8) {
// 	//
// }

// func i2c_write_word_data(b uint8, w uint16) {
// 	//
// }

// func i2c_read_word_data(b uint8, w uint16) uint16 {
// 	var d uint16
// 	return d
// }

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