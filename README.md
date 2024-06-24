## Getting Started

    Kubestream is a CLI tool similar to Kubectl (client side) but streams resources metadata and information from multiple Kubernetes clusters simultaneously.

### Usage:

1. Runtime Go 1.19
2. Need to add your kubeconfig.(json|yaml) file and context name into ./config/kubeobject.yaml (you can define n number of cluster config and can group it by group_by)

*note: 
    group_by is just defined for useability purposes. you might want to group cluster test, stage, prod, or region and if you want to fetch resource-only stage then --group_by="stage" can be used. If you want to get all the cluster resources then use --group_by="all"*

*note: 
    kubestream has been designed with concurrency hence you will have to pass kubeconfig path as an individual unique path and kubeconfig.yaml file to avoid race-condition while switching context between Kubernetes clusters.*



```
kubernetes_cluster:
  - name_alias: "{correct kubectl CONTEXT_NAME}"
    kubeconfig: "{correct kubeconfig.(json|yaml) file} of the above context"
    group_by: "{GROUP_NAME}"
```

**Example:**

```
kubernetes_cluster:
  - name_alias: ******-minikube-****/******/admin
    kubeconfig: "/Users/minikube/kube-config.yaml"
    group_by: local
  - name_alias: ******-optest-****/******/admin
    kubeconfig: "/Users/******-optest-****/******/admin/kube-config.yaml"
    group_by: ibm-us-south
  - name_alias: ******-aws-eks
    kubeconfig: "/Users/******-aws-eks/kube-config.yaml"
    group_by: aws-us-east
  - name_alias: ******-aws-eks-cluster-2
    kubeconfig: "/Users/******-aws-eks-cluster-2/kube-config.yaml"
    group_by: aws-us-east
  - name_alias: ******-aws-eks-cluster-3
    kubeconfig: "/Users/******-aws-eks-cluster-3/kube-config.yaml"
    group_by: aws-us-east
....
```

3. make build
4. Run binary **./kubestream**



CLI Usage:
---
```
./kubestream --help

kubestream get --api_resource=deployment --namespace="all" --group_by="aws-us-east"
kubestream get --api_resource=daemonset --namespace="default" --group_by="local"
kubestream get --api_resource=statefulset --namespace="default" --group_by="ibm-us-south"
kubestream get --api_resource=statefulset --namespace="kube-system" --group_by="all"
```
