#!/bin/bash

# Q100 PA Server for Raspberry Pi 4
# Orignal design by Michael, EA7KIR

GOVERSION=1.21.0

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

echo Updateing Pi OS
sudo apt update
sudo apt -y full-upgrade
sudo apt -y autoremove
sudo apt clean

echo Running rfkill # not sure if this dupicates config.txt
rfkill block 0
rfkill block 1

echo Making changes to config.txt
cd /boot
echo Enabling I2C
sudo sed -i 's/^#dtparam=i2c_arm=on/dtparam=i2c_arm=on/' config.txt
cd

echo Enabling 1-Wire
echo -e "\ndtoverlay=w1-gpio" >> /boot/config.txt

echo Disbaling Wifi
echo -e "\ndtoverlay=disable-wifi" >> /boot/config.txt

echo Disbaling Bluetooth
echo -e "\ndtoverlay=disable-bt" >> /boot/config.txt

echo Installing GIT
sudo apt -y install git

echo Installing i2c-tools
sudo apt -y install i2c-tools

echo Adding go path to .profile
echo -e '\n\nexport PATH=$PATH:/usr/local/go/bin\n\n' >> /home/pi/.profile

echo Installing Go $GOVERSION
GOFILE=go$GOVERSION.linux-arm64.tar.gz
cd /usr/local
sudo wget https://go.dev/dl/$GOFILE
# sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf $GOFILE
cd

echo Cloning q100receiver to /home/pi/Q100
cd
mkdir Q100
cd Q100
git clone https://github.com/ea7kir/q100paserver.git
cd

echo Copying q100receiver.service
cd /home/pi/Q100/etc
sudo cp q100paserver.service /etc/systemd/system/
sudo chmod 644 /etc/systemd/system/q100paserver.service
sudo systemctl daemon-reload
cd

echo "\n
INSTALL HAS COMPLETED
   after rebooting, build and auto exec...

   cd Q100/q100paserver
   go mod tidy
   go build .
   sudo systemctl enable q100paserver
   sudo systemctl start q100paserver

   now type sudo reboot
"
