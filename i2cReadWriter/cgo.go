/*
 *  Q-100 PA Server
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package i2cReadWriter

// #include <linux/i2c-dev.h>
import "C"

// Get I2C_SLAVE constant value from
// Linux OS I2C declaration file.
const (
	I2C_SLAVE = C.I2C_SLAVE
)
