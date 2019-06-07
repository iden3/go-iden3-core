#!/bin/sh
set -ex

VER=1.0.17-stable
wget https://download.libsodium.org/libsodium/releases/libsodium-${VER}.tar.gz
tar -xzvf libsodium-${VER}.tar.gz
cd libsodium-stable
./configure --prefix=/usr
make && make check
sudo make install
