#!/bin/env sh

docker build -t knuckle-go .
sleep 10
docker run --name knuckle-go-container -d knuckle-go
sleep 10
docker exec -it knuckle-go-container sh
