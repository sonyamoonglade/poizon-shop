#!/bin/bash
ARCH=amd64
OS=linux
IP=$(echo $VM_IP)
USER=aalexandrovich

mkdir deploy
GOOS=$OS GOARCH=$ARCH go build -o ./deploy/api cmd/main.go
echo "building..."

cp ./scripts/run.sh ./deploy/run.sh

# remove old binary
ssh -i ~/.ssh/vadim-shop $USER@$IP "rm -rf ~/build/api/api ~/build/api/run.sh"

# transfer binary 
scp -r -i ~/.ssh/vadim-shop deploy/* $USER@$IP:./build/api/
echo "deploying build folder"

# stop existing session
ssh -i ~/.ssh/vadim-shop $USER@$IP "kill -9 \$(pidof api)"
echo "stopped running api process"

rm -rf deploy
