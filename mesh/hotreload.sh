#!/bin/sh
pkill app && echo "Killing process..."
rm -f ./app && echo "Removing binary..."
echo "Re-building & restarting."
go build . && ./app