#!/bin/bash

CONTAINER_NAME="auth"

if [ "$(docker ps -aq -f name=$CONTAINER_NAME)" ]; then
    if [ "$(docker ps -q -f name=$CONTAINER_NAME)" ]; then
        echo "The container '$CONTAINER_NAME' is already running."
    else
        echo "Starting the existing container '$CONTAINER_NAME'."
        docker start $CONTAINER_NAME
    fi
else
    echo "Creating and running a new container '$CONTAINER_NAME'."
    docker run -d --name $CONTAINER_NAME $CONTAINER_NAME:latest
fi