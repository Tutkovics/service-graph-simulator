# service-graph-simulator
Project to create service graph.

## Usage
```
$ kubectl apply -f service-graph-crd.yaml 
customresourcedefinition.apiextensions.k8s.io/servicegraphs.service.graph.com created

$ kubectl apply -f my-service-graph.yaml 
servicegraph.service.graph.com/my-service-mesh created

$ kubectl get servicegraphs [-o json]
```

## Latest yaml to deploy
```
sample-graphs.yaml
```
