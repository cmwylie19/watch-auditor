.DEFAULT_GOAL := docker-image

IMAGE ?= cmwylie19/watch-auditor:v0.0.1
PROD_IMAGE ?= cmwylie19/watch-auditor:prod

build-push-prod-image:
	docker buildx build --platform linux/amd64,linux/arm64 --push -t $(PROD_IMAGE) -f Dockerfile.amd .

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

deploy-dev:
	kind create cluster
	docker build -t watch-auditor:dev . -f Dockerfile.arm
	kind load docker-image watch-auditor:dev
	kubectl apply -k kustomize/overlays/dev

clean-dev:
	kind delete cluster --name kind
	docker system prune -a -f

check-metrics:
	kubectl run -it curler -n watch-auditor --image=nginx:alpine --rm -it --restart=Never  -- curl -s http://watch-auditor.watch-auditor.svc.cluster.local:8080/metrics | grep watch_controller_failures_total

check-logs:
	kubectl logs -n watch-auditor -l app=watch-auditor -f | jq
