/*
 *  Q-100 PA Server
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package temperature

import (
	"sync"

	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
)

// DS18B20 TEMPERATURE SENSORS
// To enable the 1-wire bus add "dtoverlay=w1-gpio" to /boot/config.txt and reboot.
// For permissions, add "/sys/bus/w1/devices/28*/w1_slave r" to /opt/pigpio/access.
// Default connection is data line to GPIO 4 (pin 7).
// 4k7 pull-up on data line to 3V3
//
// Set the slave ID for each DS18B20 TO-92 device
// To find those available, type: cd /sys/bus/w1/devices/
// and look for directories named like: 28-3c01d607d440

// TRY https://pkg.go.dev/periph.io/x/conn/v3/onewire

// 'ls /sys/bus/w1/devices/' on my setup yeilds the floowing...
const (
	kPreampSensorAddress = "28-3c01d607d440"
	kPaSensorAddress     = "28-3c01d607e348"
	kDefault1wirePin     = rpi.J8p7
)

type (
	ds18b20Type struct {
		line    *gpiod.Line
		mu      sync.Mutex
		newtemp float64
		temp    float64
		slaveId string
	}
)

var (
	preAmp  ds18b20Type
	finalPA ds18b20Type
)

func newDs18b20(j8Pin int, slaveId string) ds18b20Type {
	// not using gpiod yet
	// line, err := gpiod.RequestLine("gpiochip0", kDefault1wirePin, gpiod.AsInput)
	// if err != nil {
	// 	logger.Fatal.Panicf("Request 1-Wire failed: %v", err)
	// }
	// return ds18b20Type{line: line, slaveId: slaveId}
	return ds18b20Type{line: nil, slaveId: slaveId}
}

func Configure() {
	preAmp = newDs18b20(kDefault1wirePin, kPreampSensorAddress)
	finalPA = newDs18b20(kDefault1wirePin, kPaSensorAddress)
}

func Shutdown() {
	// revert lines to input on the way out
}

func PreAmp() float64 {
	return tempForSensor(&preAmp)
}

func FinalPA() float64 {
	return tempForSensor(&finalPA)
}

func tempForSensor(sen *ds18b20Type) float64 {

	sen.newtemp = 0.0

	sen.mu.Lock()
	sen.temp = sen.newtemp
	sen.mu.Unlock()
	return sen.temp
}
