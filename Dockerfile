FROM golang

# ADD codeship-dind ./codeship-dind
WORKDIR $GOPATH/src/github.com/ds0nt/codeship-dind
RUN go get github.com/Masterminds/glide
ADD glide.yaml .
ADD glide.lock .
ADD Makefile .
RUN make glide-install
ADD . .
RUN make build

# should be backwards compatible
ENV DOCKER_API_VERSION=1.24

CMD ./codeship-dind
