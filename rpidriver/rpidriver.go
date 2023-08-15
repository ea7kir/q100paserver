/*
 *  Q-100 PA Server
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package rpidriver

import (
	"os/exec"
	"q100paserver/logger"
	"strconv"
	"strings"
	"sync"
	"time"
)

type (
	rpiType struct {
		mu    sync.Mutex
		quit  chan bool
		tempC float64
	}
)

var (
	rpiCpu *rpiType
)

func newRpi() *rpiType {
	return &rpiType{
		mu:    sync.Mutex{},
		quit:  make(chan bool),
		tempC: 0.0,
	}
}

func Configure() {
	rpiCpu = newRpi()
	go readRpi(rpiCpu)
}

func Shutdown() {
	rpiCpu.quit <- true
}

/*
pi@txserver:~ $ /usr/bin/vcgencmd measure_temp
temp=53.5'C
pi@txserver:~ $
*/

func Temperature() float64 {
	rpiCpu.mu.Lock()
	defer rpiCpu.mu.Unlock()
	return rpiCpu.tempC
}

// Go routine to read raspberry pi core data
//
//	An alternative legacy way is read
//	sys/class/thermal/thermal_zone0/temp
//	51121
func readRpi(pi *rpiType) {
	var tempC float64
	var err error
	var bytes []byte
	for {
		select {
		case <-pi.quit:
			return
		default:
		}
		tempC = 0.0
		bytes, err = exec.Command("vcgencmd", "measure_temp").Output()
		if err != nil {
			logger.Error.Printf("Failed to read rpi temperature: %v", err)
		}
		str0 := string(bytes[:])
		str1 := strings.Split(str0, "=")
		str2 := strings.Split(str1[1], "'C")
		str3 := strings.TrimSpace(str2[0])
		tempC, err = strconv.ParseFloat(str3, 64)
		if err != nil {
			logger.Error.Printf("Failed to convert rpi temperature: %v", err)
		} else {
			rpiCpu.mu.Lock()
			rpiCpu.tempC = tempC
			rpiCpu.mu.Unlock()
			time.Sleep(5 * time.Second)
		}
	}
}
