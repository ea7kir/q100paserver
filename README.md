# Q-100 PA Server
Monitors Pre-Amplifier and Power Amplier temperatures, and the Enclosure and PA fans speeds. Data is returned to q100transmitter when connected.
### REQUIRES Raspberry PI OS Lite (64-BIT) - the Bookworm Lite version

## Hardware
- Raspberry Pi 4B with 4GB RAM (minimum)
- Waveshare RPi relay board
- 2 x DS18B20 temperature sensors
- 1 x INA226 current/voltage sensor
- 4 x 12v fans
- 1 x 5v power supply
- 1 x 12v power supply
- 1 x 28v power supply
- 2 x 12/230v contactors

## Connections
- find INA266 address using `sudo i2cdetect -y 1`
- find DS18B20 addresses using `ls /sys/bus/w1/devices/`

TODO: add more details and photos

## Installing
**A keyboard and mouse are not required at any time**

### Using Raspberry Pi Imager v1.8.5:
```
CHOOSE Raspberry Pi Device: Raspberry Pi 4

CHOOSE OS: Raspberry Pi OS (other) -> Raspberry Pi OS Lite (64-bit)

CONFIGURE:
	Set hostname:			paserver
	Enable SSH
		Use password authentication
	Set username and password
		Username:			pi
		Password: 			<password>
	Set locale settings
		Time zone:			<Europe/Madrid> # or wherever you are
		Keyboard layout:	<us>
	Eject media when finished
SAVE and WRITE
```

Insert the card into the Raspberry Pi and switch on

NOTE: the Pi will reboot during the install, so please allow it to complete

### Remote login from a Mac, PC or Linux host
```
ssh pi@paserver.local

sudo apt update
sudo apt -y install git
mkdir Q100
cd Q100
git clone https://github.com/ea7kir/q100paserver.git

cd q100paserver/etc
chmod +x install.sh
./install.sh
```

### After rebooting
Login from a PC, Mac, or Linux computer
```
ssh pi@paserver.local or open VSCODE to PaServer  ~/Q100/q100paserver
```
Execute the following commands
```
cd Q100/q100paserver
go mod tidy
go build .
```
If all went well it can be executed at boot by enabling systemctl
```
sudo systemctl enable q100paserver
sudo systemctl start q100paserver
```

## License
Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see https://www.gnu.org/licenses/.
