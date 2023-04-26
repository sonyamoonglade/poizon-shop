#!/bin/bash
export $(xargs < ~/build/household/.env)
mkdir -p ~/logs/household
logs_name=$(echo ~/logs/household/$(date +'%m.%d.%Y.%T').log)
touch $(echo $logs_name)
cd ~/build/household
nohup ./household -strict=false -config-path=/configs/config.yml -logs-path=$logs_name -production >/dev/null 2>&1 &


