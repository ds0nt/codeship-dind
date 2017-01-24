FROM golang

# ADD codeship-dind ./codeship-dind
WORKDIR $GOPATH/src/github.com/ds0nt/codeship-dind
ADD . .
RUN make glide-install
RUN make build

# should be backwards compatible
ENV DOCKER_API_VERSION=1.24

CMD codeship-dind
