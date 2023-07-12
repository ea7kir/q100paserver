/*
 *  Q-100 PA Server
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package temperature

import "fmt"

func Configure(pi int) {
	//
}

func Shutdown() {
	//
}

func Read() string {
	str := fmt.Sprintf("Pre %4.1f°C PA %4.1f°C",
		readPreAmp(), readPA())
	return str
}

func readPreAmp() float64 {
	return 53.1
}

func readPA() float64 {
	return 43.8
}
