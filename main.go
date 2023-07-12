/*
 *  Q-100 PA Server
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package main

import (
	"fmt"
	"q100paserver/current"
	"q100paserver/fan"
	"q100paserver/logger"
	"q100paserver/power"
	"q100paserver/temperature"
	"time"
)

func configureDevices() {
	pi := 0
	power.Configure(pi)
	current.Configure(pi)
	temperature.Configure(pi)
	fan.Configure(pi)
}

func readDevices() string {
	str := fmt.Sprintf("%v %v %v",
		temperature.Read(),
		current.Read(),
		fan.Read())
	return str
}

func shutdownDevices() {
	power.Shutdown()
	fan.Shutdown()
	current.Shutdown()
	temperature.Shutdown()

}

func handler() {
	connected := true
	power.Up()
	for {
		if !connected {
			break
		}
		str := readDevices()
		logger.Info.Printf("\n\tSEND: %v\n", str)
		time.Sleep(2 * time.Second)
	}
	power.Down()
}

// TODO: add signal to cancel

func runServer() {
	connected := true
	if connected {
		handler()
	}
}

func main() {
	logger.Info.Printf("Q-100 PA Server will start...")
	configureDevices()
	runServer()
	shutdownDevices()
	logger.Info.Printf("Q-100 PA Server has stopped")
	// TODO: shutdown or reboot Rasberry Pi
}
