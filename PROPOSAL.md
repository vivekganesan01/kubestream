
My rational thoughts:

My current issue:

- due to choosing flag for --api_resource as an input, i am not able to fully support some of native object query. Eg:

```
--api_resource=deployment # might not have state like pod or container
--api_resource=pod # might not have actual and desired replicas 
```

The problem with my old approach is (--api_resource) if user wants to filter pod based no unhealthy state and if accidently passed that filter into deployments then cli might 
fail to differenciate this and capture this and will might fail to act accordingly

so here is my proposal

```
# lets putup dedicated commands for api_resources like kubectl

kubestream get deployment --namespace=foo --filter=uptodate=0,name=xyz
kubestream get pod --namespace=bar --filter=state=norunning,status=crashlooping
```

so with my new proposal there is another issue, what if user wants to filter multiresource

```
# eg : get both deployment and statefulset
kubestream get deployment,statefulset,secrets ?????
```

Solution:

will create a new resource fetcher called `api_resources`

```
kubestream get api_resources --resource_type=deployment,statefulset,secrets --namespace=bar
```

note: on the above case, the custom filtering or querying is not possible.