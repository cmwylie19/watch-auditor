# Watch Auditor

Configurable audit period check to ensure Watch Controller is working containing two modes: 
- Enforcing - Reports Watch failures through Prometheus metrics and can roll the Watch Controller pods in UDS Core on failed audit.
- Audit - Reports Watch failures through Prometheus metrics.

Why? Because it is _light-weight_ and easy.

## Configurability:

The Watch Auditor is configured through args in the Deployment's container spec.  

```bash
Start the server

Usage:
  watch-auditor serve [flags]

Flags:
  -e, --every duration          Interval to check in seconds (default 30s) (default 30s)
  -f, --failure-threshold int   Failure threshold to roll watch controller pod (default 3)
  -h, --help                    help for serve
  -l, --log-level string        Log level (debug, info, error) (default "info")
  -m, --mode string             Mode to run in (audit, enforcing) (default "enforcing")
  -p, --port int                Port to listen on (default: 8080) (default 8080)
```



## Overview 

UDS Core is a microservice architecture that uses Kubernetes Watch to monitor/enact the state of the cluster. The Watch Auditor is a tool that ensures that the Watchers are functioning correctly. The Watch Auditor periodically checks the state of the Watcher and reports any failures through Prometheus metrics. The Watch Auditor can be configured to run in two modes: audit and enforcing. In audit mode, the Watch Auditor reports any failures but does not take any action. In enforcing mode, the Watch Auditor reports any failures and can take action to restart the Watcher. The Watch Auditor is configured through command line arguments. The Watch Auditor could be critical component of the UDS Core architecture as it ensures that the Watchers are functioning correctly and that the state of the cluster is being monitored correctly.

**Background:**
Kubernetes watch works correctly through the KFC Informer pattern most of the time, it is possible in extreme cases it is susceptible to failure. The Watch Auditor optimizes around those failure scenarios by reporting failures and taking action to restart the Watcher which solves the problem. 

**How it works:**
UDS Core uses NeuVector to run security tooling as a `CronJob`. The neuvector namespace is labeled with `istio-injection=enabled` which means that all pods in the namespace will have an Envoy sidecar injected. The Watch Auditor is deployed in the `watch-auditor` namespace. It creates "test" pods in the `neuvector` namespace that quickly die after 5 seconds, however, the Envoy sidecar should still be running. The job of the watcher is to notice the Envoy and kill it, which terminates the pod. If the pod continues lingering, then the Watcher has failed. The Watch Auditor checks if the pod is still alive after 10 seconds. If the pod is still alive after 10 seconds, this indicates that the Watcher has not caught the sidecar and it rolls and/or reports on the watcher pods.


## Metics

```bash
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 1.3333e-05
go_gc_duration_seconds{quantile="0.25"} 0.000138626
go_gc_duration_seconds{quantile="0.5"} 0.000263833
go_gc_duration_seconds{quantile="0.75"} 0.000436083
go_gc_duration_seconds{quantile="1"} 0.000436083
go_gc_duration_seconds_sum 0.000851875
go_gc_duration_seconds_count 4
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 10
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.22.4"} 1
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 2.559544e+06
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 6.783152e+06
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 9577
# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 33417
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 3.060888e+06
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 2.559544e+06
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 2.801664e+06
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 4.964352e+06
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 17454
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 1.736704e+06
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 7.766016e+06
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 1.7192354189262578e+09
# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter
go_memstats_lookups_total 0
# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 50871
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 7200
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 15600
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 112160
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 130560
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 5.429552e+06
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 1.459999e+06
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 622592
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 622592
# HELP go_memstats_sys_bytes Number of bytes obtained from system.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 1.3065232e+07
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 8
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 0.19
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds 1.048576e+06
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 10
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 2.9958144e+07
# HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
# TYPE process_start_time_seconds gauge
process_start_time_seconds 1.71923514813e+09
# HELP process_virtual_memory_bytes Virtual memory size in bytes.
# TYPE process_virtual_memory_bytes gauge
process_virtual_memory_bytes 1.289691136e+09
# HELP process_virtual_memory_max_bytes Maximum amount of virtual memory available in bytes.
# TYPE process_virtual_memory_max_bytes gauge
process_virtual_memory_max_bytes 1.8446744073709552e+19
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 2
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
# HELP watch_controller_deletions_total The total number of watch controller deletions
# TYPE watch_controller_deletions_total counter
watch_controller_deletions_total 3 # IMPORTANT!!!
# HELP watch_controller_failures_total The total number of watch controller failures
# TYPE watch_controller_failures_total counter
watch_controller_failures_total 10 # IMPORTANT!!!

```

## Developing


```bash
There are 10 types of people in this world, those who understand binary and those who don't ¯\_(ツ)_/¯.
```

Reallll quick restart/reset of simulated environment:

```bash
k3d cluster delete --all;
docker system prune -a -f 
k3d cluster create;
docker build -t watch-auditor:dev .;
k3d image import watch-auditor:dev -c k3s-default  
k apply -f k8s

istioctl install --set profile=demo -y
k create ns pepr-system
k create ns neuvector
k label ns neuvector istio-injection=enabled
k apply -f -<<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: pepr-deploy
    pepr.dev/controller: watcher
  name: pepr-deploy
  namespace: pepr-system
spec:
  replicas: 1
  selector:
    matchLabels:
      app: pepr-deploy
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: pepr-deploy
        pepr.dev/controller: watcher
    spec:
      containers:
      - image: nginx
        name: nginx
        resources: {}
status: {}
EOF
k run curler -n watch-auditor --image=nginx
```

Check the k3d images  
```bash
docker exec -it k3d-k3s-default-server-0 crictl images
```

Logs and metrics 
```bash
k logs -n watch-auditor -l app=watch-auditor -f


k exec -it -n watch-auditor curler -- curl watch-auditor:8080/metrics
```



