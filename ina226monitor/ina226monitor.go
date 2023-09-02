/*
 *  Q-100 PA Server
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package ina226monitor

import (
	"fmt"
	"os"
	"q100paserver/ina226Driver"
	"sync"
	"time"

	"github.com/ea7kir/qLog"
)

const (
	kI2cBus            = 1
	kFinalPaI2cAddress = 0x40
	kFinalPaShunt      = 0.002 // modifiedto 0.0021 from 0,002 to get correct current reading
	kFinalPaMaxAmps    = 10
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

type (
	ina226sensorType struct {
		mu    sync.Mutex
		quit  chan bool
		volts float64
		amps  float64
		// address int8
		// shunt   float64
		// maxAmps float64
		sensor *ina226Driver.Ina226
	}
)

var (
	finalPA *ina226sensorType
)

func newIna226sensor(address int8, shunt float64, maxAmps float64) *ina226sensorType {
	return &ina226sensorType{
		mu:    sync.Mutex{},
		quit:  make(chan bool),
		volts: 0.0,
		amps:  0.0,
		// address: address,
		// shunt:   shunt,
		// maxAmps: maxAmps,
		sensor: nil,
	}
}

func Configure() {
	// for Final PA
	finalPA = newIna226sensor(kFinalPaI2cAddress, kFinalPaShunt, kFinalPaMaxAmps)
	sensor, err := ina226Driver.NewDriver(kI2cBus, kFinalPaI2cAddress)
	if err != nil {
		qLog.Fatal("failed NewDriver %s: ", err)
		os.Exit(1)
	}
	err = sensor.Configure(
		ina226Driver.INA226_SHUNT_CONV_TIME_1100US,
		ina226Driver.INA226_BUS_CONV_TIME_1100US,
		ina226Driver.INA226_AVERAGES_1,
		ina226Driver.INA226_MODE_SHUNT_BUS_CONT,
	)
	if err != nil {
		qLog.Fatal("failed Configure %s: ", err)
		os.Exit(1)
	}
	sensor.Calibrate(kFinalPaShunt, kFinalPaMaxAmps)
	finalPA.sensor = sensor
	// for any other goes here

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
func readVoltsAmpsFor(sensor *ina226sensorType) {
	for {
		select {
		case <-sensor.quit:
			return
		default:
		}

		vBus, err := sensor.sensor.ReadBusVoltage()
		if err != nil {
			qLog.Error("%s", err)
		}
		iShunt, err := sensor.sensor.ReadShuntCurrent()
		if err != nil {
			qLog.Error("%s", err)
		}
		fmt.Printf("%v volts, %v amps\n", vBus, iShunt)

		sensor.mu.Lock()
		sensor.volts = vBus
		sensor.amps = iShunt
		sensor.mu.Unlock()
		time.Sleep(1333 * time.Millisecond)
	}
}
