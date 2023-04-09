#!/bin/bash

export $(xargs < .env)
go build -o ./build/app cmd/app/main.go
./build/app -strict=false -config-path=./config.yml
