#!/bin/bash
export $(xargs < .env)
go build -o build/household cmd/main.go
./build/household -strict=false -config-path ./config.yml


