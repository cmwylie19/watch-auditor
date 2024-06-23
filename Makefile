.DEFAULT_GOAL := docker-image

IMAGE ?= cmwylie19/watch-auditor:latest

compile: 
	GOARCH=arm64 CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o watch-auditor main.go
	mv watch-auditor image/watch-auditor


build-push-image: 
	docker buildx build --push --platform linux/amd64 -t $(IMAGE) image/
