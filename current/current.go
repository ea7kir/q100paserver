/*
 *  Q-100 PA Server
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package current

import "fmt"

func Configure(pi int) {
	//
}

func Shutdown() {
	//
}

func Read() string {
	str := fmt.Sprintf("%3.1f amp",
		readPaCurrent())
	return str
}

func readPaCurrent() float64 {
	return 1.3
}
