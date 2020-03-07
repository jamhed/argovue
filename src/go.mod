module argovue

go 1.13

require (
	github.com/argoproj/argo v2.4.3+incompatible
	github.com/argoproj/pkg v0.0.0-20191031223000-02a6aac40ac4 // indirect
	github.com/aws/aws-sdk-go v1.19.11
	github.com/colinmarc/hdfs v1.1.4-0.20180805212432-9746310a4d31 // indirect
	github.com/coreos/go-oidc v2.1.0+incompatible
	github.com/fluxcd/helm-operator v1.0.0-rc5
	github.com/go-openapi/spec v0.19.5 // indirect
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.7.3
	github.com/gorilla/sessions v1.2.0
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/hashicorp/go-uuid v1.0.1 // indirect
	github.com/jcmturner/gofork v1.0.0 // indirect
	github.com/onsi/ginkgo v1.10.3 // indirect
	github.com/onsi/gomega v1.7.1 // indirect
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/valyala/fasttemplate v1.1.0 // indirect
	golang.org/x/net v0.0.0-20191207000613-e7e4b65ae663 // indirect
	golang.org/x/oauth2 v0.0.0-20191202225959-858c2ad4c8b6
	golang.org/x/sys v0.0.0-20191206220618-eeba5f6aabab // indirect
	gopkg.in/boj/redistore.v1 v1.0.0-20160128113310-fc113767cd6b
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/jcmturner/goidentity.v2 v2.0.0 // indirect
	gopkg.in/square/go-jose.v2 v2.4.0 // indirect
	k8s.io/api v0.15.7
	k8s.io/apimachinery v0.15.7
	k8s.io/client-go v12.0.0+incompatible
)

replace github.com/docker/distribution => github.com/2opremio/distribution v0.0.0-20190419185413-6c9727e5e5de

replace github.com/docker/docker => github.com/docker/docker v0.7.3-0.20190327010347-be7ac8be2ae0
