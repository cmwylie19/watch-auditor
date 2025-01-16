# Watch Auditor

The Watch Auditor is a tool that ensures that measures watch misses between Pepr and the `kube-apiserver`. The [soak test](https://github.com/defenseunicorns/pepr-excellent-examples/tree/main/hello-pepr-soak-ci) is configured to delete pods from the `pepr-demo` namespace. The Watch Auditor deploys pods at an interval and checks to see if Pepr properly deleted them  and reports any failures through Prometheus metrics.

- [Configurability](#configurability)
- [Metrics](#metrics)
- [Developing](#developing)

**NOTE** Watch Auditor needs amd64 image for GitHub Actions. 
```bash
docker buildx create --use --name prod-builder --driver docker-container
make build-push-prod-image
```

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

Deploy Dev Environment:

```bash
make deploy-dev
```
Check Dev Logs:

```bash
make check-logs
```

Check Dev Metrics:

```bash
make check-metrics
```

Clean up Dev Environment:

```bash
make clean-dev
```

unit tests:

```bash
make unit-test
```

integration tests:

```bash
make e2e-test
```

