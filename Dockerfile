FROM golang

# ADD codeship-dind ./codeship-dind
WORKDIR $GOPATH/src/github.com/ds0nt/codeship-dind
ADD . .
RUN go get
RUN make build
CMD codeship-dind
