/*
 *  Q-100 PA Server
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package fandriver

import (
	"q100paserver/mylogger"
	"sync"
	"time"

	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
)

const (
	// FAN SENSORS with 1k0 pull-up resistors on sensor lines to 3.3v

	kEncIntakePin  = rpi.J8p29 // pin 29 GPIO_5
	kEncExtractPin = rpi.J8p31 // pin 31 GPIO_6
	kPaIntakePin   = rpi.J8p33 // pin 33 GPIO_13
	kPaExtractPin  = rpi.J8p35 // pin 35 GPIO_19
)

type (
	fanType struct {
		line *gpiod.Line
		mu   sync.Mutex
		quit chan bool
		rpm  int64
	}
)

var (
	encIntake  *fanType
	encExtract *fanType
	paIntake   *fanType
	paExtract  *fanType
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

func newFan(j8Pin int) *fanType {
	// const deboucePeriod = 3 * time.Millisecond
	// WithDebounce(deboucePeriod)
	l, err := gpiod.RequestLine("gpiochip0", j8Pin, gpiod.AsInput)
	if err != nil {
		mylogger.Fatal.Fatalf("line %v failed: %v", l, err)
	}
	return &fanType{
		line: l,
		mu:   sync.Mutex{},
		// ch:     make(chan int64),
		quit: make(chan bool),
		// newRpm: 0,
		rpm: 0,
	}
}

func Configure() {
	encIntake = newFan(kEncIntakePin)
	encExtract = newFan(kEncExtractPin)
	paIntake = newFan(kPaIntakePin)
	paExtract = newFan(kPaExtractPin)
	go rpmForFan(encIntake)
	go rpmForFan(encExtract)
	go rpmForFan(paIntake)
	go rpmForFan(paExtract)
}

func Shutdown() {
	// kill the go routines
	encIntake.quit <- true
	encExtract.quit <- true
	paIntake.quit <- true
	paExtract.quit <- true
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
	encIntake.mu.Lock()
	defer encIntake.mu.Unlock()
	return encIntake.rpm
}

func EnclosureExtract() int64 {
	encExtract.mu.Lock()
	defer encExtract.mu.Unlock()
	return encExtract.rpm
}

func FinalPAintake() int64 {
	paIntake.mu.Lock()
	defer paIntake.mu.Unlock()
	return paIntake.rpm
}

func FinalPAextract() int64 {
	paExtract.mu.Lock()
	defer paExtract.mu.Unlock()
	return paExtract.rpm
}

// Go Routine calculates rpm for each fan
//
//	4000 rpm equates to 8000 ppm or 133 pps
//	ie. 1 pulse every 7.5 milliseconds
func rpmForFan(fan *fanType) {
	// 4000 rpm equates to 8000 ppm or 133 pps
	// ie. 1 pulse every 7.5 milliseconds
	var newRpm int64
	for {
		newRpm = 0
		const loopTime = 1003 * time.Millisecond
		var i int
		for start := time.Now(); ; {
			// no need to checl end time quite so often, slow it down by 10 iterations
			if i%10 == 0 {
				select {
				case <-fan.quit:
					return
				default:
				}
				if time.Since(start) > loopTime {
					break
				}
			}
			i++
			v1, err := fan.line.Value()
			if err != nil {
				mylogger.Warn.Printf(" %v", err)
			}
			time.Sleep(3 * time.Millisecond)
			v2, err := fan.line.Value()
			if err != nil {
				mylogger.Warn.Printf(" %v", err)
			}
			if v1 != v2 {
				newRpm += 30
			}
		}
		fan.mu.Lock()
		fan.rpm = newRpm
		fan.mu.Unlock()
		time.Sleep(1 * time.Second)
	}
}
