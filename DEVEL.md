# Development

Deploy Kubernetes objects to configure `argovue` from `kube` folder, and:

```sh
skaffold dev --port-forward
```

After successful deployment point your browser to `http://localhost:8080/ui/`.

## Flux sync

```
kubectl -n fluxcd port-forward deployment/helm-operator 3030:3030 &
curl -XPOST http://localhost:3030/api/v1/sync-git
```
