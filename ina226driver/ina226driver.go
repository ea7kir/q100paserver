/*
 *  Q-100 PA Server
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package ina226driver

import (
	"fmt"
	"log"
	ina226 "q100paserver/ina266"
	"sync"
	"time"
)

// TODO: install i2c-tools
// sudo apt install -y i2c-tools
// TODO: enable I2C in config

// INA226 current/voltage sensors
// To discover I2C devices
// $ sudo i2cdetect -y 1
// TODO: address could be 0x40, 0x41 or 0x42
// I2C pin3 GPIO2 SDA, pin 5 GPIO3 SCL
// 4k7 pull-up resistors on data lines to 3.3v

// Configuration values
const (
	kFinalPaAddrees = 0x40
	kFinalPaShunt   = 0.0021 // modified from 0,002 to get correct current reading
	kFinalPaMaxAmps = 10
)

/* PYTHON CONSTANTS
const (
	INA226_RESET                     uint16 = 0x8000 // UInt16
	INA226_REG_CALIBRATION           uint8  = 0x05   // UInt8
	INA226_REG_CONFIGURATION         uint8  = 0x00   // UInt8
	INA226_TIME_8MS                  uint8  = 7      // 8.244ms UInt8
	INA226_AVERAGES_16               uint8  = 2      // UInt8
	INA226_MODE_SHUNT_BUS_CONTINUOUS uint8  = 7      // UI
	INA226_REG_BUS_VOLTAGE           uint8  = 0x02   // UInt8
	INA226_REG_CURRENT               uint8  = 0x04   // UInt8
)
*/

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

/*******************************************************************
* the main
*******************************************************************/

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func xmain() {

	// INA226 A0 and A1 tied to GND, address is 0x40
	// device connected to I2C bus 1
	currentSensor, err := ina226.New(1, kFinalPaAddrees)
	checkError(err)

	// err = currentSensor.Reset()
	// checkError(err)

	err = currentSensor.Configure(ina226.INA226_SHUNT_CONV_TIME_1100US, ina226.INA226_BUS_CONV_TIME_1100US, ina226.INA226_AVERAGES_1, ina226.INA226_MODE_SHUNT_BUS_CONT)
	checkError(err)

	// measure the power rail voltage
	vBus, err := currentSensor.ReadBusVoltage()
	checkError(err)
	fmt.Printf("Bus voltage: %2.5f V\n", vBus)
	// measure the voltage drop across the shunt resistor
	vShunt, err := currentSensor.ReadShuntVoltage()
	checkError(err)
	fmt.Printf("Shunt voltage: %2.5f V\n", vShunt)

	// If you want to read the current directly you must calibrate the sensor first
	// providing the Shunt resistor value (expressed in ohm) and
	// the maximum Expected current (expressed in Ampere).
	// This values are required to set the resolution of the readings
	currentSensor.Calibrate(kFinalPaShunt, kFinalPaMaxAmps) // TODO: DOUBLE CHECK THESE VALUES

	// read back the resolution
	iResolution, err := currentSensor.CurrentResolution()
	checkError(err)
	fmt.Printf("Current resolution: %2.5f A/bit\n", iResolution)

	iRegister, err := currentSensor.ReadShuntCurrentRegister()
	checkError(err)
	fmt.Println("Shunt current register:", iRegister)

	iShunt, err := currentSensor.ReadShuntCurrent()
	checkError(err)
	fmt.Printf("Shunt current: %2.5f A\n", iShunt)

	// stop acquisition
	//err = currentSensor.Configure(ina226.INA226_MODE_POWER_DOWN)
	//checkError(err)
}
