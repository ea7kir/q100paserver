/*
 *  Q-100 PA Server
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package current

import (
	"sync"

	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
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
	kDefaultSDAPin  = rpi.J8p3
	kDefaultSCLPin  = rpi.J8p5

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
		lineSDA *gpiod.Line
		lineSCL *gpiod.Line
		mu      sync.Mutex
		newAmps float64
		amps    float64
		address int8
		shunt   float64
		maxAmps float64
	}
)

var (
	finalPA ina226Type
)

func newIna226(address int8, shunt float64, maxAmps float64) ina226Type {
	// not using gpiod yet
	// lineSDA, err := gpiod.RequestLine("gpiochip0", kDefaultSDAPin, gpiod.AsInput)
	// if err != nil {
	// 	logger.Fatal.Panicf("Request lineSDA failed: %v", err)
	// }
	// lineSCL, err := gpiod.RequestLine("gpiochip0", kDefaultSCLPin, gpiod.AsInput)
	// if err != nil {
	// 	logger.Fatal.Panicf("Request lineSCL failed: %v", err)
	// }
	// return ina226Type{lineSDA: lineSDA, lineSCL: lineSCL, address: address, shunt: shunt, maxAmps: maxAmps}
	return ina226Type{lineSDA: nil, lineSCL: nil, address: address, shunt: shunt, maxAmps: maxAmps}
}

func Configure() {
	finalPA = newIna226(kFinalPaAddrees, kFinalPaShunt, kFinalPaMaxAmps)
}

func Shutdown() {
	// revert lines to input on the way out
}

func FinalPA() float64 {
	return ampsForSensor(&finalPA)
}

func byteSwapped(w uint16) uint16 {
	// 16 bit word byte swap from TxServer
	b1 := w >> 8
	b2 := w & 0xFF
	result := b2 << 8
	result |= b1
	return result
}

const i2c_bus = 1

func i2c_open(bus uint8, address uint8) {
	//
}

func i2c_write_word_data(b uint8, w uint16) {
	//
}

func i2c_read_word_data(b uint8, w uint16) uint16 {
	var d uint16
	return d
}

func ampsForSensor(sen *ina226Type) float64 {

	sen.newAmps = 0.0

	sen.mu.Lock()
	sen.amps = sen.newAmps
	sen.mu.Unlock()
	return sen.amps
}
