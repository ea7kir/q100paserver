/*
 *  Q-100 PA Server
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package i2cReadWriter

import (
	"fmt"
	"os"
	"syscall"
)

type Device struct {
	bus  int
	addr uint8
	rc   *os.File
}

func NewDevice(bus int, addr uint8) (*Device, error) {
	f, err := os.OpenFile(fmt.Sprintf("/dev/i2c-%d", bus), os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	if err := ioctl(f.Fd(), I2C_SLAVE, uintptr(addr)); err != nil {
		return nil, err
	}
	d := &Device{
		bus:  bus,
		addr: addr,
		rc:   f,
	}
	return d, nil
}

func (d *Device) Write(buf []byte) (int, error) {
	return d.rc.Write(buf)
}

func (d *Device) Read(buf []byte) (int, error) {
	return d.rc.Read(buf)
}

func ioctl(fd, cmd, arg uintptr) error {
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, cmd, arg, 0, 0, 0)
	if err != 0 {
		return err
	}
	return nil
}
