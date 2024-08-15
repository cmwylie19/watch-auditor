# Watch Auditor

The Watch Auditor is a tool that ensures that measures watch misses between Pepr and the `kube-apiserver`. The [soak test](https://github.com/defenseunicorns/pepr-excellent-examples/tree/main/hello-pepr-soak-ci) is configured to delete pods from the `pepr-demo` namespace. The Watch Auditor deploys pods at an interval and checks to see if Pepr properly deleted them  and reports any failures through Prometheus metrics.

- [Configurability](#configurability)
- [Metrics](#metrics)
- [Developing](#developing)
- [Check Logs and Metrics](#check-logs-and-metrics)  

## Configurability:

The Watch Auditor is configured through args in the Deployment's container spec.  

```bash
Start the server

Usage:
  watch-auditor serve [flags]

Flags:
  -e, --every duration          Interval to check in seconds (default 30s)
  -h, --help                    help for serve
  -l, --log-level string        Log level (debug, info, error) (default "info")
  -p, --port int                Port to listen on (default 8080)
  -n, --namespace string        Namespace to check (default "pepr-demo")
```


## Metics

```bash
watch_controller_failures_total 10 

```

## Developing

Quick Restart

```bash
k3d cluster delete --all;
docker system prune -a -f 
k3d cluster create;
docker build -t auditor:dev -f Dockerfile.arm .;
k3d image import auditor:dev -c k3s-default  
k apply -f k8s
```

build/push image:

```bash
make build-arm-image
# or 
make build-push-arm-image
```

```bash
make build-amd-image
# or 
make build-push-amd-image
```


## Check Logs and Metrics
```bash
k logs -n watch-auditor -l app=watch-auditor -f

kubectl run curler -n watch-auditor --image=nginx
k exec -it -n watch-auditor curler -- curl watch-auditor:8080/metrics
```


## Test

unit tests:

```bash
make unit-test
```

integration tests:

```bash
make e2e-test
```

```bash
