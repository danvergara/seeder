#!/bin/bash

# allow specifying different destination directory
DIR="${DIR:-"/usr/local/bin"}"

# get the OS
OS=$(uname -s)
case $OS in
    Linux) OS=linux ;;
    Darwin) OS=darwin ;;
esac

# map different architecture variations to the available binaries
ARCH=$(uname -m)
case $ARCH in
    i386|i686|x86_64) ARCH=amd64 ;;
    armv6*) ARCH=armv6 ;;
    armv7*) ARCH=armv7 ;;
    aarch64*) ARCH=arm64 ;;
esac

# prepare the download URL
GITHUB_LATEST_VERSION=$(curl -L -s -H 'Accept: application/json' https://github.com/danvergara/seeder/releases/latest | sed -e 's/.*"tag_name":"\([^"]*\)".*/\1/')
GITHUB_FILE="seeder_${GITHUB_LATEST_VERSION//v/}_${OS}_${ARCH}.tar.gz"
GITHUB_URL="https://github.com/danvergara/seeder/releases/download/${GITHUB_LATEST_VERSION}/${GITHUB_FILE}"

echo $GITHUB_FILE
echo $GITHUB_LATEST_VERSION
echo $GITHUB_URL

# install/update the local binary
curl -L -o seeder.tar.gz $GITHUB_URL
tar xzvf seeder.tar.gz
sudo mv -f seeder "$DIR"
rm seeder.tar.gz
