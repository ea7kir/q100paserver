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
echo "-- Updating eeprom firmware"
echo "-------------------------------"
echo

sudo rpi-update

echo
echo "-------------------------------"
echo "-- Enable 1-Wire and I2C"
echo
echo "-- Enable with sudo raspi-config"
echo
echo "-- The reboot and run install_2"
echo "-------------------------------"
echo
