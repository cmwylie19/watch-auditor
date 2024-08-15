.DEFAULT_GOAL := docker-image

IMAGE ?= cmwylie19/watch-auditor:v0.0.1
IMAGE2 ?= watch-auditor

build-arm-image: 
	docker build -t $(IMAGE) -f Dockerfile.arm .

build-amd-image:
	docker build -t $(IMAGE) -f Dockerfile.amd .

build-push-arm-image: 
	docker buildx build --push -t $(IMAGE) -f Dockerfile.arm .

build-push-amd-image: 
	docker buildx build --push -t $(IMAGE) -f Dockerfile.amd .

unit-test:
	go test -v ./... -tags='!e2e'

e2e-test:
	ginkgo -v --tags='e2e' ./e2e

