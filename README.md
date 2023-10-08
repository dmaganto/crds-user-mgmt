# idp-k8s
Create custom CRDs to manage users, teams and applications

Install cluster
```bash
kind create cluster --name idp
```

Install CRDs
```bash
kustomize build . | kubectl apply -f -
```

Use randomize to create new resources for each kind of resource.
```bash
python3 -m venv .venv
source .venv/bin/activate
pip install jinja2
for i in {0..100};do python3 randomize.py| kubectl apply -f -;done
```

Perform queries to filter all users that belongs to some team
```bash
k get developers -o json | jq '.items[] | select(.spec.teams[] == "claims") | .metadata.name'
k get applications -o json | jq '.items[] | select(.spec.team == "front") | .metadata.name'
```

Perform queries directly to the API (pod inside the cluster using serviceaccount):
```bash
export APISERVER=https://kubernetes.default.svc 
export SERVICEACCOUNT=/var/run/secrets/kubernetes.io/serviceaccount
export NAMESPACE=$(cat ${SERVICEACCOUNT}/namespace)
export TOKEN=$(cat ${SERVICEACCOUNT}/token)
export CACERT=${SERVICEACCOUNT}/ca.crt

curl --cacert ${CACERT} --header "Authorization: Bearer ${TOKEN}" -X GET ${APISERVER}/apis/dmaganto.infra/v1alpha1/namespaces/default/applications | jq '.items[] | select(.spec.team == "front") | .metadata.name'
curl --cacert ${CACERT} --header "Authorization: Bearer ${TOKEN}" -X GET ${APISERVER}/apis/dmaganto.infra/v1alpha1/namespaces/default/developers | jq '.items[] | select(.spec.teams[] == "claims") | .metadata.name'
```

Perform queries with python script, take a look to application folder.
