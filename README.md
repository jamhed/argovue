# ArgoVue

This is work in progress, check [TODO](TODO.md), also please consult [INSTALL](INSTALL.md) and [DEVEL](DEVEL.md) files.

## Motivation

Provide UI for custom objects (argo workflows) with actions with authentication (OIDC) and authorization (group membership),
and expose services uniformly with authenticating reverse proxy.

## Use case

Provide per project (namespace) UI to run and manage argo workflows and expose pre-defined services with uniform access management.

## Pre-requisites

OIDC server is required for the application to work. It could be either external OIDC provider (Okta, Auth0), or included
dex.

## Usage

Make workflow visible for group `authors`:

```sh
kubectl -n $NAMESPACE label workflow/$NAME oidc.argovue.io/group=authors
```
