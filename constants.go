package main

// TODO: install i2c-tools
// sudo apt install -y i2c-tools
// TODO: enable I2C in config

// There"s an INA266 C library at https://github.com/MarioAriasGa/raspberry-pi-ina226

const (
	SERVER_PORT = 8765

	// INA226 current/voltage sensors
	// To discover I2C devices
	// $ sudo i2cdetect -y 1
	// TODO: address could be 0x40, 0x41 or 0x42
	// I2C pin3 GPIO2 SDA, pin 5 GPIO3 SCL
	// 4k7 pull-up resistors on data lines to 3.3v
	PA_CURRENT_ADDRESS   = 0x40
	PA_CURRENT_SHUNT_OHM = 0.0021 // modified from 0,002 to get correct current reading
	PA_CURRENT_MAX_AMP   = 10

	// FAN SENSORS
	// 1k0 pull-up resistors on sensor lines to 3.3v
	ENCLOSURE_INTAKE_FAN_GPIO  = 5  // pin 29 GPIO_5
	ENCLOSURE_EXTRACT_FAN_GPIO = 6  // pin 31 GPIO_6
	PA_INTAKE_FAN_GPIO         = 13 // pin 33 GPIO_13
	PA_EXTRACT_FAN_GPIO        = 19 // pin 35 GPIO_19

	// DS18B20 TEMPERATURE SENSORS
	// To enable the 1-wire bus add "dtoverlay=w1-gpio" to /boot/config.txt and reboot.
	// For permissions, add "/sys/bus/w1/devices/28*/w1_slave r" to /opt/pigpio/access.
	// Default connection is data line to GPIO 4 (pin 7).
	// 4k7 pull-up on data line to 3V3
	//
	// Sset the slave ID for each DS18B20 TO-92 device
	// To find those available, type: cd /sys/bus/w1/devices/
	// and look for directories named like: 28-3c01d607d440
	PA_SENSOR_SLAVE_ID     = "28-3c01d607e348" // pin 7 GPIO_4
	PREAMP_SENSOR_SLAVE_ID = "28-3c01d607d440" // pin 7 GPIO_4

	// WAVESHARE RPi RELAY BOARD
	RELAY_28v_GPIO = 26 // pin 37 GPIO_26 (CH1 P25)
	RELAY_12v_GPIO = 20 // pin 38 GPIO_20 (CH2 P28)
	RELAY_5v_GPIO  = 21 // pin 40 GPIO_21 (CH3 P29)
	// NOTE: the opto coupleers need reverse logic
	RELAY_ON  = 0
	RELAY_OFF = 1
)
