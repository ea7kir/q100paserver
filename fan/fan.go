/*
 *  Q-100 PA Server
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package fan

import (
	"fmt"
	"math/rand"

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

var (
	encIntake  *gpiod.Line
	encExtract *gpiod.Line
	paIntake   *gpiod.Line
	paExtract  *gpiod.Line
)

func Configure(pi int) {
	encIntakeLine, err := gpiod.RequestLine("gpiochip0", rpi.J8p29, gpiod.AsInput)
	if err != nil {
		panic(err)
	}
	encIntake = encIntakeLine
	encExtractLine, err := gpiod.RequestLine("gpiochip0", rpi.J8p31, gpiod.AsInput)
	if err != nil {
		panic(err)
	}
	encExtract = encExtractLine
	paIntakeLine, err := gpiod.RequestLine("gpiochip0", rpi.J8p33, gpiod.AsInput)
	if err != nil {
		panic(err)
	}
	paIntake = paIntakeLine
	paExtractLine, err := gpiod.RequestLine("gpiochip0", rpi.J8p35, gpiod.AsInput)
	if err != nil {
		panic(err)
	}
	paExtract = paExtractLine
}

func Shutdown() {
	// revert lines to input on the way out
	encIntake.Reconfigure(gpiod.AsInput)
	encIntake.Close()
	encExtract.Reconfigure(gpiod.AsInput)
	encExtract.Close()
	paIntake.Reconfigure(gpiod.AsInput)
	paIntake.Close()
	paExtract.Reconfigure(gpiod.AsInput)
	paExtract.Close()
}

func Read() string {
	str := fmt.Sprintf("Ein %v Eout %v PAin %v Pout %v",
		readEncIntake(), readEncExtract(), readPaIntake(), readPaExtract())
	return str
}

func readEncIntake() int {
	min := 4000
	max := 4999
	r := rand.Intn(max-min) + min
	return r
}

func readEncExtract() int {
	min := 4000
	max := 4999
	r := rand.Intn(max-min) + min
	return r
}

func readPaIntake() int {
	min := 4000
	max := 4999
	r := rand.Intn(max-min) + min
	return r
}

func readPaExtract() int {
	min := 4000
	max := 4999
	r := rand.Intn(max-min) + min
	return r
}
