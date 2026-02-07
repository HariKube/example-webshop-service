#! /bin/bash -e

: ${REGISTRY_PASSWORD?= required}
: ${HARIKUBE_URL:=https://harikube.info}

export KIND_CLUSTER=kind

export SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)

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
exe kubectl create secret generic -n harikube harikube-license --from-file=docs/demo/license
exe kubectl create secret docker-registry harikube-registry-secret \
--docker-server=registry.harikube.info \
--docker-username=harikube \
--docker-password='${REGISTRY_PASSWORD}' \
--namespace=harikube
exe kubectl apply -f ${HARIKUBE_URL}/manifests/harikube-operator-release-v1.0.0.yaml
exe kubectl apply -f ${HARIKUBE_URL}/manifests/harikube-middleware-vcluster-api-release-v1.0.0.yaml
exe kubectl wait -n harikube --for=jsonpath='{.status.readyReplicas}'=1 deployment/operator-controller-manager --timeout=2m
exe kubectl wait -n harikube --for=jsonpath='{.status.readyReplicas}'=1 statefulset/harikube --timeout=5m

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
exe kubectl create namespace example-webshop-service-system
exe make install
exe "echo '
apiVersion: v1
kind: Secret
metadata:
  name: remote-example-webshop-service-system
  namespace: example-webshop-service-system
  annotations:
    kubernetes.io/service-account.name: "example-webshop-service-controller-manager"
type: kubernetes.io/service-account-token
' | kubectl apply -f -
"
exe kubectl create namespace serverless-kube-watch-trigger-system
exe kubectl apply -f https://github.com/HariKube/serverless-kube-watch-trigger/releases/download/beta-v1.0.0-7/bundle-rbac.yaml
exe "echo '
apiVersion: v1
kind: Secret
metadata:
  name: remote-serverless-kube-watch-trigger
  namespace: serverless-kube-watch-trigger-system
  annotations:
    kubernetes.io/service-account.name: "serverless-kube-watch-trigger-controller-manager"
type: kubernetes.io/service-account-token
' | kubectl apply -f -
"
exe vcluster disconnect

exe kubectl create namespace serverless-kube-watch-trigger-system
exe 'echo "
apiVersion: v1
kind: Config
clusters:
- name: remote-serverless-kube-watch-trigger-system
  cluster:
    server: harikube.harikube.svc.service.local
    certificate-authority-data: $(kubectl get secret -n harikube -l vcluster.loft.sh/namespace=serverless-kube-watch-trigger-system -o jsonpath='{.items[0].data.ca\.crt}')
contexts:
- name: my-context
  context:
    cluster: remote-serverless-kube-watch-trigger-system
    user: remote-serverless-kube-watch-trigger-system
    namespace: default
current-context: my-context
users:
- name: remote-serverless-kube-watch-trigger-system
  user:
    token: $(kubectl get secret -n harikube -l vcluster.loft.sh/namespace=serverless-kube-watch-trigger-system -o jsonpath='{.items[0].data.token}')
" | kubectl create secret generic -n serverless-kube-watch-trigger-system remote-kubeconfig --from-file=kubeconfig=/dev/stdin
'
exe '(cd $(mktemp -d) && echo '\''
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - https://github.com/HariKube/serverless-kube-watch-trigger/releases/download/beta-v1.0.0-7/bundle.yaml
patches:
  - target:
      kind: Deployment
      name: serverless-kube-watch-trigger-controller-manager
    patch: |
      apiVersion: apps/v1
      kind: Deployment
      metadata:
        name: serverless-kube-watch-trigger-controller-manager
      spec:
        template:
          spec:
            containers:
            - name: manager
              env:
              - name: KUBECONFIG
                value: "/etc/kube/config"
              volumeMounts:
              - name: kubeconfig-volume
                mountPath: "/etc/kube"
                readOnly: true
            volumes:
            - name: kubeconfig-volume
              secret:
                secretName: remote-kubeconfig
'\'' >> kustomization.yaml && kubectl kustomize .)
'

exe "TAG=snapshot-$(date +'%s') make docker-build docker-load package"

exe kubectl create namespace example-webshop-service-system
exe 'echo "
apiVersion: v1
kind: Config
clusters:
- name: remote-example-webshop-service-system
  cluster:
    server: harikube.harikube.svc.service.local
    certificate-authority-data: $(kubectl get secret -n harikube -l vcluster.loft.sh/namespace=example-webshop-service-system -o jsonpath='{.items[0].data.ca\.crt}')
contexts:
- name: my-context
  context:
    cluster: remote-example-webshop-service-system
    user: remote-example-webshop-service-system
    namespace: default
current-context: my-context
users:
- name: remote-example-webshop-service-system
  user:
    token: $(kubectl get secret -n harikube -l vcluster.loft.sh/namespace=example-webshop-service-system -o jsonpath='{.items[0].data.token}')
" | kubectl create secret generic -n example-webshop-service-system remote-kubeconfig --from-file=kubeconfig=/dev/stdin
'
exe '(cd $(mktemp -d) && cp $SCRIPT_DIR/package/bundle.yaml . && echo '\''
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - bundle.yaml
patches:
  - target:
      kind: Deployment
      name: example-webshop-service-controller-manager
    patch: |
      apiVersion: apps/v1
      kind: Deployment
      metadata:
        name: example-webshop-service-controller-manager
      spec:
        template:
          spec:
            containers:
            - name: manager
              env:
              - name: KUBECONFIG
                value: "/etc/kube/config"
              volumeMounts:
              - name: kubeconfig-volume
                mountPath: "/etc/kube"
                readOnly: true
            volumes:
            - name: kubeconfig-volume
              secret:
                secretName: remote-kubeconfig
'\'' >> kustomization.yaml && kubectl kustomize .)
'
exe kubectl apply -f ./package/config.yaml

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

pushd function
sed -i 's/# - remote/- remote/' stack.yaml
exe ../bin/faas-cli template store pull python3-http
exe ../bin/faas-cli build
exe ../bin/faas-cli push
exe ../bin/faas-cli login --password ${OPENFAASPWD} --gateway http://${KINEIP}:32767
exe ../bin/faas-cli deploy --gateway http://${KINEIP}:32767
sed -i 's/- remote/# - remote/' stack.yaml
popd

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
exe vcluster connect harikube
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