#!/bin/env sh

docker stop knuckle-go-container
docker rm knuckle-go-container --force --if-exists
if [[ "$(docker images -q my-custom-image 2> /dev/null)" != "" ]]; then
    docker rmi my-custom-image --force
else
    echo "Образ не существует"
fi

