#!/bin/bash

# Install Q100 PA Server on Raspberry Pi 4
# Orignal design by Michael, EA7KIR

# updated to GO 1.22 on March 2 2024

GOVERSION=1.22.0

whoami | grep -q pi
if [ $? != 0 ]; then
  echo Install must be performed as user pi
  exit
fi

hostname | grep -q paserver
if [ $? != 0 ]; then
  echo Install must be performed on host paserver
  exit
fi

while true; do
    read -p "Install q100paserver using Go version $GOVERSION (y/n)? " answer
    case ${answer:0:1} in
        y|Y ) break;;
        n|N ) exit;;
        * ) echo "Please answer yes or no.";;
    esac
done

###################################################

echo Update Pi OS
sudo apt update
sudo apt -y full-upgrade
sudo apt -y autoremove
sudo apt clean

###################################################

# echo Running rfkill # not sure if this dupicates config.txt
# rfkill block 0
# rfkill block 1

###################################################

echo Making changes to config.txt

sudo sh -c "echo '\n# EA7KIR Additions' >> /boot/config.txt"

echo Disable Wifi
sudo sh -c "echo 'dtoverlay=disable-wifi' >> /boot/config.txt"

echo Disable Bluetooth
sudo sh -c "echo 'dtoverlay=disable-bt' >> /boot/config.txt"

echo Enable 1-Wire
sudo sh -c "echo 'dtoverlay=w1-gpio' >> /boot/config.txt"

echo Enable I2C
# sudo sh -c "echo 'dtparam=i2c_arm=on' >> /boot/config.txt"
sudo sed -i 's/#dtparam=i2c_arm=on/dtparam=i2c_arm=on/g' /boot/config.txt

###################################################

echo Making changes to .profile

sudo sh -c "echo '\n# EA7KIR Additions' >> /home/pi/.profile"

echo Adding go path to .profile
echo -e 'export PATH=$PATH:/usr/local/go/bin' >> /home/pi/.profile

###################################################

echo Installing i2c-tools
sudo apt -y install i2c-tools

echo Installing Go $GOVERSION
GOFILE=go$GOVERSION.linux-arm64.tar.gz
cd /usr/local
sudo wget https://go.dev/dl/$GOFILE
sudo tar -C /usr/local -xzf $GOFILE
cd

echo Copying q100paserver.service
cd /home/pi/Q100/q100paserver/etc
sudo cp q100paserver.service /etc/systemd/system/
sudo chmod 644 /etc/systemd/system/q100paserver.service
sudo systemctl daemon-reload
cd

###################################################

chmod -x /home/pi/Q100/etc/install.sh # to prevent it from being run a second time

echo "
INSTALL HAS COMPLETED
   after rebooting...

   cd Q100/q100paserver
   go mod tidy
   go build .
   sudo systemctl enable q100paserver
   sudo systemctl start q100paserver

"

while true; do
    read -p "I have read the above, so continue (y/n)? " answer
    case ${answer:0:1} in
        y|Y ) break;;
        n|N ) exit;;
        * ) echo "Please answer yes or no.";;
    esac
done

sudo reboot
