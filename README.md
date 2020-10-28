# service-graph-simulator
Project to create service graph.

2+1 components:
- GO service (node_source)
<br>From this source we will bulid Docker image
- K8s operator (operator)
<br>Use: operator framework, to create Kubernetes operator to manage our CR and deploy the given service-mesh.
- Custom Resource (yaml)
<br>Sample .yaml file to deploy