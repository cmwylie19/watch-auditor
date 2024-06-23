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
```

docker exec -it k3d-k3s-default-server-0 crictl images
```
