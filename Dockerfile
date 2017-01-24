FROM golang

# ADD codeship-dind ./codeship-dind
RUN make build
CMD codeship-dind
