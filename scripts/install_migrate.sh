#!/bin/bash

# allow specifying different destination directory
DIR="${DIR:-"/usr/local/bin"}"

# map different architecture variations to the available binaries
ARCH=$(uname -m)
case $ARCH in
    x86_64) ARCH=amd64 ;;
    armv6*) ARCH=armv6 ;;
    armv7*) ARCH=armv7 ;;
    aarch64*) ARCH=arm64 ;;
esac

echo "arch: ${ARCH}"

OS=$(uname -s)
case $OS in
    Linux) OS=linux ;;
    Darwin) OS=darwin ;;
esac

echo "os: ${OS}"

# prepare the download URL
GITHUB_LATEST_VERSION=$(curl -L -s -H 'Accept: application/json' https://github.com/golang-migrate/migrate/releases/latest | sed -e 's/.*"tag_name":"\([^"]*\)".*/\1/')
GITHUB_FILE="migrate.${OS}-${ARCH}.tar.gz"
GITHUB_URL="https://github.com/golang-migrate/migrate/releases/download/${GITHUB_LATEST_VERSION}/${GITHUB_FILE}"

# install/update the local binary
curl -L -o migrate.tar.gz $GITHUB_URL
tar -xzvf migrate.tar.gz
mv "migrate.${OS}-${ARCH}" "migrate"
mv -f "migrate" "$DIR"
