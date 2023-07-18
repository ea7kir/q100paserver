/*
 *  Q-100 PA Server
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package fan

import (
	"q100paserver/logger"
	"sync"
	"time"

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
		line   *gpiod.Line
		lock   sync.Mutex
		newRpm int64
		rpm    int64
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
WithEventHandler(eh)
WithRisingEdge
WithMonotonicEventClock
const WithFallingEdge = LineEdgeFalling
const WithRisingEdge = LineEdgeRising
const WithRealtimeEventClock = LineEventClockRealtime
*/

func configured(j8Pin int) fanType {
	// const deboucePeriod = 3 * time.Millisecond
	// WithDebounce(deboucePeriod)
	l, err := gpiod.RequestLine("gpiochip0", j8Pin, gpiod.AsInput)
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

func rpmForFan(fan *fanType) int64 {
	// 4000 rpm equates to 8000 ppm or 133 pps
	// ie. 1 pulse every 7.5 milliseconds

	// runs once per client client request
	func(fan *fanType) {
		count := 0
		fan.newRpm = 0
		var i int
		for start := time.Now(); ; {
			// no need to checl end time quite so often, slow it down by 10 iterations
			if i%10 == 0 {
				if time.Since(start) > time.Second {
					logger.Info.Println("break")
					break
				}
			}
			i++
			// count chnges for 60 seconds
			v1, err := fan.line.Value()
			if err != nil {
				logger.Warn.Printf(" %v", err)
			}
			time.Sleep(3 * time.Millisecond)
			v2, err := fan.line.Value()
			if err != nil {
				logger.Warn.Printf(" %v", err)
			}
			if v1 != v2 {
				count += 1
			}
		}
		// calc newNpm value
		fan.newRpm = int64(count / 133)
		// nerge with rpm
		// fan.newRpm = 0.9*fan.newRpm + 0.1*fan.rpm
		fan.lock.Lock()
		fan.rpm = fan.newRpm
		fan.lock.Unlock()
	}(fan)

	return fan.rpm
}
