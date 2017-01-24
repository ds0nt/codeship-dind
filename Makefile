

build:
	go build .

dockerize:
	docker build -t codeship-dind .

glide-install:
	glide install --strip-vendor

.PHONY: build dockerize glide-install