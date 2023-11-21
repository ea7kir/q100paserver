#!/bin/bash

# Install Q100 PA Server on Raspberry Pi 4
# Orignal design by Michael, EA7KIR

GOVERSION=1.21.4

echo WARNING: THIS INSTALL SCRIPT HAS NOT BEEN TESTED

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

echo "\n###################################################\n"

echo Updateing Pi OS
sudo apt update
sudo apt -y full-upgrade
sudo apt -y autoremove
sudo apt clean

echo "\n###################################################\n"

# echo Running rfkill # not sure if this dupicates config.txt
# rfkill block 0
# rfkill block 1

echo "\n###################################################\n"

echo Making changes to config.txt

sudo sh -c "echo '\n# EA7KIR Additions' >> /boot/config.txt"

echo Disbaling Wifi
sudo sh -c "echo 'dtoverlay=disable-wifi' >> /boot/config.txt"

echo Disbaling Bluetooth
sudo sh -c "echo 'dtoverlay=disable-bt' >> /boot/config.txt"

echo Enabling I2C
sudo sh -c "echo 'dtparam=i2c_arm=on' >> /boot/config.txt"

echo Enabling 1-Wire
sudo sh -c "echo 'dtoverlay=w1-gpio' >> /boot/config.txt"

echo "\n###################################################\n"

echo Making changes to .profile

sudo sh -c "echo '\n# EA7KIR Additions' >> /home/pi/.profile"

echo Adding go path to .profile
echo -e 'export PATH=$PATH:/usr/local/go/bin' >> /home/pi/.profile

echo "\n###################################################\n"

echo Installing GIT
sudo apt -y install git

echo Installing i2c-tools
sudo apt -y install i2c-tools

echo Installing Go $GOVERSION
GOFILE=go$GOVERSION.linux-arm64.tar.gz
cd /usr/local
sudo wget https://go.dev/dl/$GOFILE
# sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf $GOFILE
cd

echo Copying q100paserver.service
cd /home/pi/Q100/q100paserver/etc
sudo cp q100paserver.service /etc/systemd/system/
sudo chmod 644 /etc/systemd/system/q100paserver.service
sudo systemctl daemon-reload
cd

echo "\n###################################################\n"

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
