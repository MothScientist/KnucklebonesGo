#!/bin/env sh

docker build -t knuckle-go .
docker run --name knuckle-go-container -d knuckle-go
docker exec -it knuckle-go-container sh
