/*
 *  Q-100 PA Server
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package fan

import (
	"q100paserver/logger"

	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
)

const (
	// FAN SENSORS
	// 1k0 pull-up resistors on sensor lines to 3.3v
	ENCLOSURE_INTAKE_FAN_GPIO  = 5  // pin 29 GPIO_5
	ENCLOSURE_EXTRACT_FAN_GPIO = 6  // pin 31 GPIO_6
	PA_INTAKE_FAN_GPIO         = 13 // pin 33 GPIO_13
	PA_EXTRACT_FAN_GPIO        = 19 // pin 35 GPIO_19

)

type (
	fanType struct {
		line *gpiod.Line
		rpm  int64
	}
)

var (
	encIntake  = fanType{} // TODO: use pointer !!!!!!!!!!
	encExtract = fanType{}
	paIntake   = fanType{}
	paExtract  = fanType{}
)

/*
WithDebounce(period time.Duration) DebounceOption // DebounceOption is of type time.Duration
const WithFallingEdge = LineEdgeFalling
const WithRisingEdge = LineEdgeRising
const WithRealtimeEventClock = LineEventClockRealtime
*/

func configured(j8Pin int) fanType {
	//const deboucePeriod = time.Millisecond

	l, err := gpiod.RequestLine("gpiochip0", j8Pin /*gpiod.WithDebounce(deboucePeriod),*/, gpiod.WithRealtimeEventClock)
	if err != nil {
		logger.Fatal.Panicf("line %v failed: %v", l, err)
	}
	return fanType{line: l}
}

func Configure() {
	encIntake = configured(rpi.J8p29)
	encExtract = configured(rpi.J8p31)
	paIntake = configured(rpi.J8p33)
	paExtract = configured(rpi.J8p35)
}

func Shutdown() {
	// revert lines to input on the way out
	encIntake.line.Reconfigure(gpiod.AsInput)
	encIntake.line.Close()
	encExtract.line.Reconfigure(gpiod.AsInput)
	encExtract.line.Close()
	paIntake.line.Reconfigure(gpiod.AsInput)
	paIntake.line.Close()
	paExtract.line.Reconfigure(gpiod.AsInput)
	paExtract.line.Close()
}

func EnclosureIntake() int64 {
	return rpmForFan(&encIntake)
}

func EnclosureExtract() int64 {
	return rpmForFan(&encIntake)
}

func FinalPAintake() int64 {
	return rpmForFan(&encIntake)
}

func FinalPAextract() int64 {
	return rpmForFan(&encIntake)
}

// The plan here is to measure the period between 2 pulses and calculate an rpm value.
// This value will be smoothed into the previous values and the new smoothed value will be returned
func rpmForFan(fan *fanType) int64 {
	// 4000 rpm equates to 8000 ppm or 133 pps
	// ie. 1 pulse every 7.5 milliseconds

	// wait for 1st pulse and record inTime (with a timeout - BUT HOW?)

	// if timeout then fan is not running, so return 0 rpm

	// wait for 2nd pulse and record inTime (with a timeout - BUT HOW?)

	// if timeout then fan is not running, so return 0 rpm

	// period = 2nd - 1st

	// period = A * previous_period + (A - 1) * period (where A is approx 0.9)

	// calculate rpm = period * SOME_CONSTANT

	// return rpm
	return fan.rpm
}
