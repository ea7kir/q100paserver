/*
 *  Q-100 PA Server
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package ds18b20monitor

import (
	"log"
	"os"
	"strconv"
	"time"

	"strings"
	"sync"
)

const (
	// DS18B20 TEMPERATURE SENSORS
	// To enable the 1-wire bus add "dtoverlay=w1-gpio" to /boot/config.txt and reboot.
	// For permissions, add "/sys/bus/w1/devices/28*/w1_slave r" to /opt/pigpio/access.
	// Default connection is data line to GPIO 4 (pin 7).
	// 4k7 pull-up on data line to 3V3
	//
	// Set the slave ID for each DS18B20 TO-92 device
	// To find those available, type: cd /sys/bus/w1/devices/
	// and look for directories named like: 28-3c01d607d440
	//
	// 'ls /sys/bus/w1/devices/' on my setup yeilds the floowing...

	kPaSensorSlaveId     = "28-3c01d607e348" // pin 7 GPIO_4
	kPreAmpSensorSlaveId = "28-3c01d607d440" // pin 7 GPIO_4)
)

type (
	ds18b20Type struct {
		mu      sync.Mutex
		tempC   float64
		slaveId string
	}
)

var (
	preAmp      *ds18b20Type
	finalPA     *ds18b20Type
	sensors     []*ds18b20Type
	stopChannel = make(chan struct{})
)

func newDs18b20(slaveId string) *ds18b20Type {
	return &ds18b20Type{
		mu:      sync.Mutex{},
		tempC:   0.0,
		slaveId: slaveId,
	}
}

func Configure() {
	preAmp = newDs18b20(kPreAmpSensorSlaveId)
	finalPA = newDs18b20(kPaSensorSlaveId)
	sensors = append(sensors, preAmp, finalPA)
	go readSensors(sensors, stopChannel)
}

func Shutdown() {
	close(stopChannel)
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
func readSensors(sensorList []*ds18b20Type, done chan struct{}) {
	var tempC float64
	var err error
	var data []byte
	for {
		for i := 0; i < len(sensorList); i++ {
			select {
			case <-done:
				return
			default:
			}
			sensor := sensorList[i]
			file := "/sys/bus/w1/devices/" + sensor.slaveId + "/w1_slave"
			tempC = 0.0
			data, err = os.ReadFile(file) // 75 bytes
			if err != nil {
				log.Printf("ERROR 1-Wire %s failed to read\n%s", sensor.slaveId, err)
				continue
			}
			str := string(data)
			if !strings.Contains(str, "YES") {
				// This flood the log file, so ignore it for now
				// log.Printf("WARN 1-Wire %s did not contain YES", sensor.slaveId)
				continue
			} //else {
			subStr := strings.Split(str, "t=")
			subStr1 := strings.TrimSpace(subStr[1])
			tempC, err = strconv.ParseFloat(subStr1, 64)
			if err != nil {
				log.Printf("WARN 1-Wire %s failed to create float: %s", sensor.slaveId, err)
				continue
			}
			//}
			tempC /= 1000.0
			sensor.mu.Lock()
			sensor.tempC = tempC
			sensor.mu.Unlock()
			time.Sleep(2 * time.Second)
		}
	}
}
