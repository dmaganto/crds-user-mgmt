# idp-k8s-crossplane
Use crossplane compositions to get Identity Provider



Install cluster
```
kind create cluster --name crossplane
helm repo add crossplane-stable https://charts.crossplane.io/stable
helm dep update
helm install crossplane --namespace crossplane-system --create-namespace crossplane-stable/crossplane
```