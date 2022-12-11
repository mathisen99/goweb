#!/bin/bash -e

docker image build --network="host" -f Dockerfile -t forum-img .
docker container run -p 8080:8080 --network="host" --detach --name forum-container forum-img
