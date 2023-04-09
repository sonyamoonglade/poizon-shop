#!/bin/bash

MONGO_TEST_PORT=27019
MONGO_IMAGE="mongo:6"
CONTAINER_NAME="debug_mongo-e2e"
MONGO_URI=mongodb://localhost:$MONGO_TEST_PORT
DB_NAME="testdb"
SRC=$(echo $(pwd)/../../)
export MONGO_URI=$MONGO_URI
export DB_NAME=$DB_NAME
# run mongo
CONTAINER_ID=$(docker run --rm -d -p $MONGO_TEST_PORT:27017 --name=$CONTAINER_NAME -e MONGODB_DATABASE=$DB_NAME $MONGO_IMAGE)
# run migrations
docker run -v $SRC/migrations:/migrations --network host --rm migrate/migrate -path=/migrations/ -database $MONGO_URI/$DB_NAME up
# run tests
go test -count=1 -v ./tests/

printf "container: %s\n" "$CONTAINER_ID"
docker rm -f $CONTAINER_ID
