#!/bin/bash
ARCH=amd64
OS=linux
IP=$(echo $VM_IP)
USER=aalexandrovich

mkdir deploy
GOOS=$OS GOARCH=$ARCH go build -o ./deploy/household cmd/main.go
echo "building..."

cp ./scripts/run.sh ./deploy/run.sh

#backup
ssh -i ~/.ssh/vadim-shop $USER@$IP "./backup.sh"
echo "backing up old files..."

# remove old binary
ssh -i ~/.ssh/vadim-shop $USER@$IP "rm -rf ~/build/household/household ~/build/household/run.sh"

# transfer binary 
scp -r -i ~/.ssh/vadim-shop deploy/* $USER@$IP:./build/household/
echo "deploying build folder"

# stop existing session
ssh -i ~/.ssh/vadim-shop $USER@$IP "kill -9 \$(pidof household)"
echo "stopped running household process"

rm -rf deploy
