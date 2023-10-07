# idp-k8s
Create custom CRDs to manage users, teams and applications

Install cluster
```
kind create cluster --name idp
```

Install CRDs
```
kustomize build . | kubectl apply -f -
```

Use randomize to create new resources for each kind of resource.
```
python3 -m venv .venv
source .venv/bin/activate
pip install jinja2
for i in {0..100};do python3 randomize.py| kubectl apply -f -;done
```

Perform queries to filter all users that belongs to some team
```
k get developers -o json | jq '.items[] | select(.spec.teams[] == "claims") | .metadata.name'
k get applications -o json | jq '.items[] | select(.spec.team == "front") | .metadata.name'
```