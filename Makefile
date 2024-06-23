# Makefile for building the Admission Controller server + docker image.

.DEFAULT_GOAL := docker-image

IMAGE ?= cmwylie19/watch-auditor:latest

compile: 
	GOARCH=arm64 CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o watch-auditor main.go
	mv watch-auditor image/watch-auditor


build-image: 
	docker build -t $(IMAGE) image/


push-image:
	docker push $(IMAGE)
