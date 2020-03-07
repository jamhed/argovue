#!/bin/bash 

mkdir -p $(pwd)/go/src
export GOPATH=$(pwd)/go
go get k8s.io/code-generator
go get k8s.io/apimachinery

go get github.com/fluxcd/helm-operator

ln -s $(pwd)/src $GOPATH/src/argovue

$GOPATH/src/k8s.io/code-generator/generate-groups.sh all argovue/client argovue/apis "argovue.io:v1"

rm -f src/src
#find go -exec chmod +w {} \;
#rm -rf go
