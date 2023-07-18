/*
 *  Q-100 PA Server
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package temperature

import (
	"os"
	"q100paserver/logger"
	"strconv"
	"strings"
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
	// Typical file contents
	// 73 01 4b 46 7f ff 0d 10 41 : crc=41 YES
	// 73 01 4b 46 7f ff 0d 10 41 t=23187
	file := "/sys/bus/w1/devices/" + sen.slaveId + "/w1_slave"
	sen.newtemp = 0.0
	data, err := os.ReadFile(file) // 75 bytes
	if err != nil {
		logger.Error.Printf("1-Wire %s failed to read\n%v", sen.slaveId, err)
		sen.mu.Lock()
		sen.temp = sen.newtemp
		sen.mu.Unlock()
		return sen.temp
	}
	// convert bytes to string
	str := string(data)
	if !strings.Contains(str, "YES") {
		logger.Warn.Printf("1-Wire %s did not contain YES", sen.slaveId)
		sen.mu.Lock()
		sen.temp = sen.newtemp
		sen.mu.Unlock()
		return sen.temp
	}
	i := strings.LastIndex(str, "t=")
	if i == -1 {
		logger.Warn.Printf("1-Wire %s did not contain t=", sen.slaveId)
		sen.mu.Lock()
		sen.temp = sen.newtemp
		sen.mu.Unlock()
		return sen.temp
	}
	tempC, err := strconv.ParseFloat(str[i+2:len(str)-1], 64)
	if err != nil {
		logger.Warn.Printf("1-Wire %s invalid temperature\n%v", sen.slaveId, err)
		sen.mu.Lock()
		sen.temp = sen.newtemp
		sen.mu.Unlock()
		return sen.temp
	}
	sen.newtemp = tempC / 1000.0
	sen.mu.Lock()
	sen.temp = sen.newtemp
	sen.mu.Unlock()
	return sen.temp
}
