#!/bin/sh


docker run -d --name dind-test-postgres postgres
docker ps
docker run --rm -i --link dind-test-postgres:dind-test-postgres gocode
docker ps
docker rm -f dind-test-postgres
docker ps
