# Install

In order to install and use ArgoVue following needs to be installed to Kubernetes cluster:

## Argo

```
kubectl create namespace argo
kubectl apply -n argo -f https://raw.githubusercontent.com/argoproj/argo/stable/manifests/install.yaml
```

See their [install](https://github.com/argoproj/argo) page for more details.

## Flux

```
kubectl apply -f https://raw.githubusercontent.com/fluxcd/helm-operator/master/deploy/flux-helm-release-crd.yaml
helm repo add fluxcd https://charts.fluxcd.io
kubectl create namespace fluxcd
helm upgrade -i helm-operator fluxcd/helm-operator --namespace fluxcd --set helm.versions=v3
```

See their [install](https://github.com/fluxcd/helm-operator/blob/master/chart/helm-operator/README.md) page for more details.

## ArgoVue

Provide values (ingress hosts, tls, oidc provider/dex config), see [values.yaml](https://raw.githubusercontent.com/jamhed/charts/master/argovue/values.yaml)

Install ArgoVue CRDs:

```
kubectl apply -f https://raw.githubusercontent.com/jamhed/argovue/master/kube/catalogue.yaml
kubectl apply -f https://raw.githubusercontent.com/jamhed/argovue/master/kube/config.yaml
```

Install ArgoVue Helm Chart:

```
helm repo add jamhed https://jamhed.github.io/charts/
kubectl create namespace fluxcd
helm update -i -f values.yaml --namespace argovue argovue argovue
```

More details are in helm chart [definition](https://github.com/jamhed/charts/tree/master/argovue).
