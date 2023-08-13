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
		cmd  string
		arg  string
		mu   sync.Mutex
		quit chan bool
		temp float64
	}
)

var (
	rpiCpu *rpiType
)

func newRpi() *rpiType {
	return &rpiType{
		cmd:  "vcgencmd",
		arg:  "measure_temp",
		mu:   sync.Mutex{},
		quit: make(chan bool),
		temp: 0.0,
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
	return rpiCpu.temp
}

// Go routine to read raspberry pi core data
func readRpi(pi *rpiType) {
	for {
		select {
		case <-pi.quit:
			return
		default:
		}
		bytes, err := exec.Command(pi.cmd, pi.arg).Output()
		if err != nil {
			logger.Error.Printf("Failed to read rpi temperature: %v", err)
		}
		str0 := string(bytes[:])
		str1 := strings.Split(str0, "=")
		str2 := strings.Split(str1[1], "'C")
		str3 := str2[0]
		temp, err := strconv.ParseFloat(str3, 64)
		if err != nil {
			logger.Error.Printf("Failed to convert rpi temperature: %v", err)
		} else {
			rpiCpu.mu.Lock()
			rpiCpu.temp = temp
			rpiCpu.mu.Unlock()
			time.Sleep(2 * time.Second)
		}
	}
}
