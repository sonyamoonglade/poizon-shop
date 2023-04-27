#!/bin/bash
MASTER=mongo1

docker-compose -f $(pwd)/scripts/docker-compose.mongo-cluster.yml up -d
sleep 2
docker exec $MASTER /scripts/init-cluster.sh