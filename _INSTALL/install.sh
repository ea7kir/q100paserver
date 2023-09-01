#!/bin/bash

echo
echo "-------------------------------"
echo "-- Updateing Pi OS"
echo "-------------------------------"
echo

sudo apt update
sudo apt full-upgrade -y
sudo apt autoremove -y
sudo apt clean

echo
echo "-------------------------------"
echo "-- running rfkill"
echo "-------------------------------"
echo

rfkill block 0
rfkill block 1

echo
echo "-------------------------------"
echo "-- Setting .profile"
echo "-------------------------------"
echo

echo -e '\n\nexport PATH=$PATH:/usr/local/go/bin\n\n' >> /home/pi/.profile

echo
echo "-------------------------------"
echo "-- Enable 1-Wire and I2C"
echo "-------------------------------"
echo

echo "TODO: by command line"

echo
echo "-------------------------------"
echo "-- Installing GIT"
echo "-------------------------------"
echo

sudo apt install git -y

echo
echo "-------------------------------"
echo "-- Installing i2c-tools"
echo "-------------------------------"
echo

sudo apt install -y i2c-tools

echo
echo "-------------------------------"
echo "-- Installing Go"
echo "-------------------------------"
echo

GOVERSION=go1.21.0.linux-arm64.tar.gz
cd /usr/local
sudo wget https://go.dev/dl/$GOVERSION
# sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf $GOVERSION
cd

echo
echo "-------------------------------"
echo "-- Done.  Reboot in 10 seconds"
echo
echo "Clone q100paserver from within VSCODE"
echo "https://github.com/ea7kir/q100paserver.git"
echo "-------------------------------"
echo

sleep 10
sudo reboot
