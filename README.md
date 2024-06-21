# kubestream
A command line utility to stream api resource information from multiple kubernetes cluster.

Usage:
---
sh
```
kubestream get --api_resource=deployment --namespace=default --group_by=all
kubestream get --api_resource=deployment --namespace=default --group_by=${CONTEXT_NAME}
kubestream get --api_resource=statefulset --namespace=default --group_by=${REGULAR_EXPRESSION}
kubestream get --api_resource=pod --namespace=default --group_by=${REGULAR_EXPRESSION} [--condition='deployment,statefulset' --filter='crashloop'] 
```