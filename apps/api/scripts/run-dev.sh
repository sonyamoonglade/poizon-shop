#!/bin/bash
export $(xargs < .env)
go build -o build/api cmd/main.go
./build/api -strict=false


