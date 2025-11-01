#! /bin/sh

CGO_ENABLED=0 go build -o bootstrap main.go && zip lambda-handler.zip bootstrap
