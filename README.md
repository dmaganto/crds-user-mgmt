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

```
k get all
NAME                                       TEAM   SLACK            AGE
application.dmaganto.infra/claim-service      back   back-service     33s
application.dmaganto.infra/document-service   back   claims-service   2m56s

NAME                                   FULLNAME         ROLETYPE    EMAIL             AGE
developer.dmaganto.infra/daniel.maganto   Daniel Maganto   devops      dani@dani.es      3m1s
developer.dmaganto.infra/test.user        Test User        developer   test@test.infra   2s

NAME                     SLACK
team.dmaganto.infra/back    back-team
team.dmaganto.infra/claim   claim-team
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

Perform queries from external:

```bash
TOKEN=$(k create token -n default test-python)
APISERVER="https://127.0.0.1:6443"
curl -k --header "Authorization: Bearer ${TOKEN}" -X GET ${APISERVER}/apis/dmaganto.infra/v1alpha1/namespaces/default/applications | jq '.items[] | select(.spec.team == "front") | .metadata.name'
```

Perform queries with python script inside the cluster, take a look to [application](application) folder.

To get all updates that are happening in the cluster for example for developers, you can take a look on the [watcher](watcher) folder
