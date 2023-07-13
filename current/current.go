/*
 *  Q-100 PA Server
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package current

import "fmt"

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
	PA_CURRENT_ADDRESS   = 0x40
	PA_CURRENT_SHUNT_OHM = 0.0021 // modified from 0,002 to get correct current reading
	PA_CURRENT_MAX_AMP   = 10
)

func Configure(pi int) {
	//
}

func Shutdown() {
	//
}

func Read() string {
	str := fmt.Sprintf("%3.1f amp",
		readPaCurrent())
	return str
}

func readPaCurrent() float64 {
	return 1.3
}
