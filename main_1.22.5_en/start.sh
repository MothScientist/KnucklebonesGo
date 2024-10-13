#!/bin/env sh

docker build -t knuckle-go .
sleep 5
docker run --name knuckle-go-container -d knuckle-go
