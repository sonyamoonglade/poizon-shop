#!/bin/bash

export $(xargs < .env)
go build -o ./build/app cmd/main.go
./build/app -strict=false -config-path=./config.yml
