

build:
	go build .

dockerize:
	docker build -t codeship-dind .


.PHONY: build dockerize