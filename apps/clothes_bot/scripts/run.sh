#!/bin/bash

export $(xargs < ~/build/clothing/.env)
mkdir -p ~/logs/clothing
logs_name=$(echo ~/logs/clothing/$(date +'%m.%d.%Y.%T').log)
touch $(echo $logs_name)
cd ~/build/clothing
nohup ./clothing -strict=false -config-path=/configs/config.yml -logs-path=$logs_name -production >/dev/null 2>&1 &


