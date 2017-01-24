#!/bin/sh

docker build -t codeship-dind /test-root
docker run -d --name dind-test-postgres postgres
docker run --rm -i --link dind-test-postgres:dind-test-postgres codeship-dind
docker rm -f dind-test-postgres