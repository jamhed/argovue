.PHONY: all argovue ui skaffold helm

all: argovue ui

argovue:
	cd src && GOOS=linux go build

ui:
	cd ui && yarn build

helm:
	helm package helm/argovue -d docs
	helm repo index docs --url https://jamhed.github.io/argovue/

skaffold: argovue
	cp src/argovue ../argovue-skaffold/argovue

