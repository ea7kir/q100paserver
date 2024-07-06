/*
 *  Q-100 PA Server
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package fanMonitor

import (
	"log"
	"sync"
	"time"

	"github.com/warthog618/go-gpiocdev"
	"github.com/warthog618/go-gpiocdev/device/rpi"
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
		line *gpiocdev.Line
		mu   sync.Mutex
		rpm  int64
	}
)

var (
	encIntake   *fanType
	encExtract  *fanType
	paIntake    *fanType
	paExtract   *fanType
	fans        []*fanType
	stopChannel = make(chan struct{})
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
	l, err := gpiocdev.RequestLine("gpiochip0", j8Pin, gpiocdev.AsInput)
	if err != nil {
		log.Fatalf("FATAL line %v failed: %s", l, err)
	}
	return &fanType{
		mu:   sync.Mutex{},
		line: l,
		rpm:  0,
	}
}

func Configure() {
	encIntake = newFan(kEncIntakePin)
	encExtract = newFan(kEncExtractPin)
	paIntake = newFan(kPaIntakePin)
	paExtract = newFan(kPaExtractPin)
	fans = append(fans, encIntake, encExtract, paIntake, paExtract)
	go readFanSpeeds(fans, stopChannel)
}

func Shutdown() {
	close(stopChannel)
	// revert lines to input on the way out
	encIntake.line.Reconfigure(gpiocdev.AsInput)
	encIntake.line.Close()
	encExtract.line.Reconfigure(gpiocdev.AsInput)
	encExtract.line.Close()
	paIntake.line.Reconfigure(gpiocdev.AsInput)
	paIntake.line.Close()
	paExtract.line.Reconfigure(gpiocdev.AsInput)
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
func readFanSpeeds(fanList []*fanType, done chan struct{}) {
	// 4000 rpm equates to 8000 ppm or 133 pps
	// ie. 1 pulse every 7.5 milliseconds
	var newRpm int64
	var fan *fanType
	const sampleTime = 1003 * time.Millisecond
	timer := time.NewTimer(sampleTime)
	timer.Stop()
	for {
	Loop:
		for i := 0; i < len(fanList); i++ {
			fan = fanList[i]
			newRpm = 0
			timer.Reset(1003 * time.Millisecond)
			for {
				select {
				case <-done:
					return
				case <-timer.C:
					fan.mu.Lock()
					fan.rpm = newRpm
					fan.mu.Unlock()
					time.Sleep(250 * time.Millisecond)
					continue Loop
				default:
				}
				v1, err := fan.line.Value()
				if err != nil {
					log.Printf("WARN  %s", err)
				}
				time.Sleep(3 * time.Millisecond)
				v2, err := fan.line.Value()
				if err != nil {
					log.Printf("WARN  %s", err)
				}
				if v1 != v2 {
					newRpm += 30
				}
			}
		}
	}
}
