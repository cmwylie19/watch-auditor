.DEFAULT_GOAL := docker-image

IMAGE ?= cmwylie19/watch-auditor:v0.0.1
IMAGE2 ?= watch-auditor
compile: 
	GOARCH=arm64 CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o watch-auditor main.go
	mv watch-auditor image/watch-auditor

just-build: 
	docker buildx build --platform linux/amd64 -t $(IMAGE) .

build-push-image: 
	docker buildx build --push --platform linux/amd64 -t $(IMAGE) image/
