#!/bin/bash
ARCH=amd64
OS=linux
IP=$(echo $VM_IP)
USER=aalexandrovich

mkdir deploy
GOOS=$OS GOARCH=$ARCH go build -o ./deploy/clothing cmd/main.go
cp templates.json ./deploy
cp -r videos ./deploy
echo "building..."

# remove old files
ssh -i ~/.ssh/vadim-shop $USER@$IP "rm -rf ~/build/clothing/clothing ~/build/clothing/run.sh"

cp ./scripts/run.sh ./deploy/run.sh

# transfer build folder
scp -r -i ~/.ssh/vadim-shop deploy/* $USER@$IP:./build/clothing/
echo "copying build folder"

# stop existing session
ssh -i ~/.ssh/vadim-shop $USER@$IP "kill -9 \$(pidof clothing)"
echo "stopped running process"

rm -rf deploy
