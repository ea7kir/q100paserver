#!/bin/bash

echo
echo "-------------------------------"
echo "-- Installing GIT"
echo "-------------------------------"
echo

sudo apt install git -y

git config --global user.name "ea7kir"
git config --global user.email "mikenaylorspain@icloud.com"
git config --global init.defaultBranch main

echo
echo "-------------------------------"
echo "-- Installing i2c-tools python3-smbus"
echo "-------------------------------"
echo

sudo apt install -y i2c-tools python3-smbus

echo
echo "-------------------------------"
echo "-- Installing Go"
echo
echo "-- this will take some time..."
echo "-------------------------------"
echo

echo "TODO: update .profile"

sudo wget https://go.dev/dl/go1.20.6.linux-arm64.tar.gz
sudo mv go1.20.6.linux-arm64.tar.gz /usr/local/
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.20.6.linux-arm64.tar.gz

echo
echo "-------------------------------"
echo "-- Done"
echo "-------------------------------"
echo

echo "Clone q100paserver from within VSCODE"
echo "using: https://github.com/ea7kir/q100paserver.git"
echo
echo "To run q100paserver, type: ./q100paserver"
