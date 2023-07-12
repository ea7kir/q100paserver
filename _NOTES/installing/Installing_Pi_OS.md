# Installing Pi OS from the Mac

NOTE: CURRENTLY REQUIRES PIOS BULLSEYE 64-BIT LITE

## Using Raspberry Pi Imager:

```
CHOOSE OS:	Raspberry Pi OS Lite (64-bit)

CONFIGURE:
	Set hostname:			txserver
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

Insert the card and login.

## Login

Clone the repro.

```
git clone https://github.com/ea7kir/TxServer.git
```

Change permissions for the 3 install scrpits.

```
chmod +x TxServer/_Resources/install_*
```

Exeute the 3 install scrpits in order.

```
./TxServer/_Resources/install_1.sh
./TxServer/_Resources/install_2.sh
./TxServer/_Resources/install_3.sh
```
