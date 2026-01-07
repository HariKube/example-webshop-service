#! /bin/bash -e

: ${REGISTRY_PASSWORD?= required}
: ${HARIKUBE_URL:=https://harikube.info}

export KIND_CLUSTER=kind

exe() {
    local display_cmd="$@"
    display_cmd="${display_cmd//\$\{REGISTRY_PASSWORD\}/****}"
    display_cmd="${display_cmd//$REGISTRY_PASSWORD/****}"
    echo -e "\n\033[1;36m$ $display_cmd\033[0m"
    [[ -n $DEBUG ]] && read -p ""
    eval "$@"
}

exe kind create cluster
exe kubectl wait --for=condition=Ready node/kind-control-plane --timeout=2m

KINEIP=$(kubectl get no kind-control-plane -o jsonpath='{.status.addresses[0].address}')

exe kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.16.3/cert-manager.yaml
exe kubectl apply -f https://github.com/prometheus-operator/prometheus-operator/releases/download/v0.77.1/stripped-down-crds.yaml
exe kubectl wait -n cert-manager --for=jsonpath='{.status.readyReplicas}'=1 deployment/cert-manager-webhook --timeout=2m

exe kubectl create namespace harikube
exe kubectl create secret docker-registry harikube-registry-secret \
--docker-server=registry.harikube.info \
--docker-username=harikube \
--docker-password='${REGISTRY_PASSWORD}' \
--namespace=harikube
exe kubectl apply -f ${HARIKUBE_URL}/manifests/harikube-operator-beta-v1.0.0-3.yaml
exe kubectl apply -f ${HARIKUBE_URL}/manifests/harikube-middleware-vcluster-workload-beta-v1.0.0-20.yaml
exe kubectl wait -n harikube --for=jsonpath='{.status.readyReplicas}'=1 deployment/operator-controller-manager --timeout=2m
exe kubectl wait -n harikube --for=jsonpath='{.status.readyReplicas}'=1 statefulset/harikube --timeout=3m

exe "echo '
apiVersion: harikube.info/v1
kind: TopologyConfig
metadata:
  name: topologyconfig-tenant
  namespace: harikube
spec:
  targetSecret: harikube/topology-config
  backends:
  - name: tenant
    endpoint: sqlite:///db/tenant.db?_journal=WAL&cache=shared
    customresource:
      group: product.webshop.harikube.info/v1
      kind: tenants
' | kubectl apply -f -
"
exe "echo '
apiVersion: harikube.info/v1
kind: TopologyConfig
metadata:
  name: topologyconfig-user
  namespace: harikube
spec:
  targetSecret: harikube/topology-config
  backends:
  - name: user
    endpoint: sqlite:///db/user.db?_journal=WAL&cache=shared
    customresource:
      group: product.webshop.harikube.info/v1
      kind: users
' | kubectl apply -f -
"
exe "echo '
apiVersion: harikube.info/v1
kind: TopologyConfig
metadata:
  name: topologyconfig-email
  namespace: harikube
spec:
  targetSecret: harikube/topology-config
  backends:
  - name: email
    endpoint: sqlite:///db/email.db?_journal=WAL&cache=shared
    customresource:
      group: product.webshop.harikube.info/v1
      kind: emails
' | kubectl apply -f -
"
sleep 2
exe "kubectl logs -n harikube -l app=harikube | grep 'Backends registered' | tail -1"

exe vcluster connect harikube

exe kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.16.3/cert-manager.yaml
exe kubectl wait -n cert-manager --for=jsonpath='{.status.readyReplicas}'=1 deployment/cert-manager-webhook --timeout=2m

exe "echo '
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: serverless-kube-watch-trigger-emails-rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: serverless-kube-watch-trigger-emails-role
subjects:
- kind: ServiceAccount
  name: serverless-kube-watch-trigger-controller-manager
  namespace: serverless-kube-watch-trigger-system
' | kubectl apply -f -
"
exe "echo '
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: serverless-kube-watch-trigger-emails-role
rules:
- apiGroups:
  - product.webshop.harikube.info
  resources:
  - emails
  verbs:
  - get
  - list
  - watch
' | kubectl apply -f -
"
exe kubectl apply -f https://github.com/HariKube/serverless-kube-watch-trigger/releases/download/beta-v1.0.0-7/bundle.yaml

exe "TAG=snapshot-$(date +'%s') make docker-build docker-load package"
exe kubectl apply -f ./package/bundle.yaml
exe kubectl apply -f ./package/config.yaml
exe kubectl wait -n example-webshop-service-system --for=jsonpath='{.status.readyReplicas}'=1 deployment/example-webshop-service-controller-manager --timeout=2m

exe helm repo add openfaas https://openfaas.github.io/faas-netes/
exe helm repo update
exe kubectl apply -f https://raw.githubusercontent.com/openfaas/faas-netes/master/namespaces.yml
exe helm upgrade openfaas --install openfaas/openfaas \
--namespace openfaas \
--set basic_auth=true \
--set functionNamespace=example-webshop-service-system \
--set serviceType=NodePort \
--set gateway.nodePort=32767
exe kubectl wait -n openfaas --for=jsonpath='{.status.readyReplicas}'=1 deployment/gateway --timeout=2m

OPENFAASPWD=$(kubectl get secret -n openfaas basic-auth -o jsonpath='{.data.basic-auth-password}'| base64 -d)
exe "(cd function ; ../bin/faas-cli template store pull python3-http)"
exe "(cd function ; ../bin/faas-cli build)"
exe "(cd function ; ../bin/faas-cli push)"
exe "(cd function ; ../bin/faas-cli login --password ${OPENFAASPWD} --gateway http://${KINEIP}:32767)"
exe "(cd function ; ../bin/faas-cli deploy --gateway http://${KINEIP}:32767)"
exe 'kubectl patch deployment email -n example-webshop-service-system --type=json -p='\''[{"op": "replace", "path": "/spec/template/spec/serviceAccountName", "value": "example-webshop-service-controller-manager"}]'\'''
exe kubectl rollout status deployment/email -n example-webshop-service-system --timeout=1m

exe "echo '
apiVersion: triggers.harikube.info/v1
kind: HTTPTrigger
metadata:
  name: example-webshop-service-email
  namespace: example-webshop-service-system
spec:
  resource:
    apiVersion: product.webshop.harikube.info/v1
    kind: Email
  eventTypes:
    - ADDED
  url:
    service:
      name: gateway
      namespace: openfaas
      portName: http
      scheme: http
      uri:
        static: /function/email
  method: POST
  body:
    contentType: application/json
    template: |
      {{ toJson . }}
  delivery:
    timeout: 10s
    retries: 3
' | kubectl apply -f -
"

exe echo kubectl logs -n example-webshop-service-system -l app.kubernetes.io/name=example-webshop-service --since=0
exe echo kubectl logs -n serverless-kube-watch-trigger-system -l app.kubernetes.io/name=serverless-kube-watch-trigger --since=0
exe echo kubectl logs -n example-webshop-service-system -l faas_function=email --since=0

sleep 10
exe "echo '
apiVersion: product.webshop.harikube.info/v1
kind: RegistrationRequest
spec:
  user:
    firstName: Richard
    lastName: Kovacs
    email: kovacsricsi@gmail.com
    phoneNumber: '+15551234567'
  password: Password123!
  tenant:
    companyName: Tech Innovators Inc. 
    country: USA
    city: New York
    address: 123 Main St, Apt 4B
    postalCode: '10001'
    taxNumber: ABC-123456789
' | kubectl create --raw /apis/api.product.webshop.harikube.info/v1/registrations/richard-kovacs -f -
"