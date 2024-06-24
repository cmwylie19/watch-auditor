# Watch Auditor

Configurable period check to ensure watch is working containing two modes: Watch and Audit. Meant to be a very small utility.

## Configurability:

* Intervals in which checks are done
* Auditing and Enforcing
* Rolls Pepr Watch Controller pods in UDS Core on failed audit



## Overview 

Deploys a short lived Pod in NeuVector namespace which receives an Envoy sidecar. The main container will exit after 5 seconds. It will check if the pod is alive at 10 seconds. If the Pod is still alive after 10 seconds this indicates that the watcher has not caught the sidecar and it rolls and/or reports on the watcher pods. 



## Developing

import dev image into k3d

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



