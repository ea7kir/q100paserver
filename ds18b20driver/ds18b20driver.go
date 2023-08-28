/*
 *  Q-100 PA Server
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package ds18b20driver

import (
	"os"
	"q100paserver/mylogger"
	"strconv"
	"time"

	"strings"
	"sync"
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
)

type (
	ds18b20Type struct {
		mu      sync.Mutex
		quit    chan bool
		tempC   float64
		slaveId string
	}
)

var (
	preAmp  *ds18b20Type
	finalPA *ds18b20Type
	busBusy = sync.Mutex{} // to prevent concurrent access to the 1-wire bus
	//                        but could this deadlock on shutdown?
)

func newDs18b20(slaveId string) *ds18b20Type {
	return &ds18b20Type{
		mu:      sync.Mutex{},
		quit:    make(chan bool),
		tempC:   0.0,
		slaveId: slaveId,
	}
}

func Configure() {
	preAmp = newDs18b20(kPreampSensorAddress)
	finalPA = newDs18b20(kPaSensorAddress)
	go readTemperatureFor(preAmp)
	go readTemperatureFor(finalPA)
}

func Shutdown() {
	preAmp.quit <- true
	finalPA.quit <- true
}

func PreAmpTemperature() float64 {
	preAmp.mu.Lock()
	defer preAmp.mu.Unlock()
	return preAmp.tempC
}

func PaTemperature() float64 {
	finalPA.mu.Lock()
	defer finalPA.mu.Unlock()
	return finalPA.tempC
}

// Go routine to read temperature
//
//	Typical file contents
//	73 01 4b 46 7f ff 0d 10 41 : crc=41 YES
//	73 01 4b 46 7f ff 0d 10 41 t=23187
func readTemperatureFor(sensor *ds18b20Type) {
	var tempC float64
	file := "/sys/bus/w1/devices/" + sensor.slaveId + "/w1_slave"
	for {
		select {
		case <-sensor.quit:
			return
		default:
		}
		tempC = 0.0
		busBusy.Lock()
		time.Sleep(475 * time.Millisecond)
		data, err := os.ReadFile(file) // 75 bytes
		time.Sleep(475 * time.Millisecond)
		busBusy.Unlock()
		if err != nil {
			mylogger.Error.Printf("1-Wire %s failed to read\n%v", sensor.slaveId, err)
		}
		str := string(data)
		if !strings.Contains(str, "YES") {
			mylogger.Warn.Printf("1-Wire %s did not contain YES", sensor.slaveId)
		} else {
			subStr := strings.Split(str, "t=")
			subStr1 := strings.TrimSpace(subStr[1])
			tempC, err = strconv.ParseFloat(subStr1, 64)
			if err != nil {
				mylogger.Warn.Printf("1-Wire %s failed to create float: %s", sensor.slaveId, err)
			}
		}
		tempC /= 1000.0
		sensor.mu.Lock()
		sensor.tempC = tempC
		sensor.mu.Unlock()
		time.Sleep(5 * time.Second)
	}
}
