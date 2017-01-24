#!/bin/bash

docker build -t codeship-dind .
docker run -d --name dind-test-postgres postgres
docker run --rm -i --link dind-test-postgres:postgres codeship-dind
docker rm -f dind-test-postgres