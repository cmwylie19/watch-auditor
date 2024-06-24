.DEFAULT_GOAL := docker-image
# DON'T USE THIS! This is changed to build in the container and currently is not outfitted to build in linux.
# It is just a POC that is rough around the edges
IMAGE ?= cmwylie19/watch-auditor:latest
IMAGE2 ?= watch-auditor
compile: 
	GOARCH=arm64 CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o watch-auditor main.go
	mv watch-auditor image/watch-auditor

just-build: 
	docker buildx build --platform linux/amd64 -t $(IMAGE) .

build-push-image: 
	docker buildx build --push --platform linux/amd64 -t $(IMAGE) image/
