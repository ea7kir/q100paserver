#!/bin/bash

# to copy this file from Mac to Rpi, type
# scp install_2.sh pi@txserver.local:.
# and chmod +x install_2.sh

cd
echo "Installing TxServer Part 2"

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
echo "-- Installing PYENV"
echo "-------------------------------"
echo

sudo apt install make build-essential libssl-dev zlib1g-dev \
libbz2-dev libreadline-dev libsqlite3-dev wget curl llvm \
libncursesw5-dev xz-utils tk-dev libxml2-dev libxmlsec1-dev libffi-dev liblzma-dev -y

git clone https://github.com/pyenv/pyenv.git ~/.pyenv

echo 'export PYENV_ROOT="$HOME/.pyenv"' >> ~/.bashrc

echo 'export PATH="$PYENV_ROOT/bin:$PATH"' >> ~/.bashrc

echo -e 'if command -v pyenv 1>/dev/null 2>&1; then\n eval "$(pyenv init --path)"\nfi' >> ~/.bashrc

#exec $SHELL

echo
echo "-------------------------------"
echo "-- Rebooting in 5 seconds"
echo "--"
echo "-- Then run install_3.sh"
echo "-------------------------------"
echo

sleep 5

sudo reboot
