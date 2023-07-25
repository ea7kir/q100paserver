## Hardware Connections

TODO: list pin connections and refer to drawings

## Installing Pi OS

NOTE: CURRENTLY REQUIRES PI OS BULLSEYE 64-BIT (LITE VERSION)

### Using Raspberry Pi Imager:

```
CHOOSE OS: Raspberry Pi OS (other) -> Raspberry Pi OS Lite (64-bit)

CONFIGURE:
	Set hostname:			q100paserver
	Enable SSH
		Use password authentication
	Set username and password
		Username:			pi
		Password: 			<password>
	Set locale settings
		Time zone:			<Europe/Madrid>
		Keyboard layout:	<us>
	Eject media when finished
SAVE and WRITE
```

Insert the card into the Raspberry Pi and switch on

WARNING: the Pi may reboot during the install, so please allow it to complete

## Remote login from a Mac, PC or Linux host

```
ssh pi@q100paserver.local
```

## More to follow goes here
