#!/bin/env sh

docker stop knuckle-go-container
docker rm knuckle-go-container --force
docker rmi knuckle-go --force
