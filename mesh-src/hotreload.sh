#!/bin/bash
sed -re 's/time.Unix\(0, ([0-9]*)\)/time.Unix\(0, '"$(($(date +%s%N)))"'\)/' -i main.go
GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -o binary -ldflags '-w' -mod=vendor
cp binary ../../../mesh/bin/binary