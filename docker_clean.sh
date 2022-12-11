#!/bin/bash -e

docker stop forum-container
docker system prune -f --volumes
