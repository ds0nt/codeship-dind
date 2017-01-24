FROM golang

# ADD codeship-dind ./codeship-dind
ADD . .
RUN make build
CMD codeship-dind
