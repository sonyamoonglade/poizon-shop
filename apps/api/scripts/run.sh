#!/bin/bash
export $(xargs < ~/build/api/.env)
mkdir -p ~/logs/api
logs_name=$(echo ~/logs/api/$(date +'%m.%d.%Y.%T').log)
touch $(echo $logs_name)
cd ~/build/api
nohup ./api -strict=false -logs-path=$logs_name -production >/dev/null 2>&1 &


