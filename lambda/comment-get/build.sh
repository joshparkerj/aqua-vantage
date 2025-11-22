#! /bin/sh

# TODO: make staticcheck work in here pls
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bootstrap main.go && zip lambda-handler.zip bootstrap
