/*
 *  Q-100 PA Server
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package temperature

import (
	"math/rand"
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

const (
	PA_SENSOR_SLAVE_ID     = "28-3c01d607e348" // pin 7 GPIO_4
	PREAMP_SENSOR_SLAVE_ID = "28-3c01d607d440" // pin 7 GPIO_4
)

func Configure() {
	//
}

func Shutdown() {
	// revert lines to input on the way out
}

// func Read() string {
// 	str := fmt.Sprintf("Pre %4.1f°C PA %4.1f°C",
// 		readPreAmp(), readPA())
// 	return str
// }

func PreAmp() float64 {
	min := 50
	max := 55
	r := rand.Intn(max-min) + min
	return float64(r)
}

func FinalPA() float64 {
	min := 40
	max := 45
	r := rand.Intn(max-min) + min
	return float64(r)
}
