YOKE_ATC_VERISON := 0.12.3
YOKECD_VERSION := 0.12.3

APIS := cluster mylist
REGISTRY := ghcr.io
REPOSITORY := avarei/yoke-test
TAG := $(shell git describe --tags --always --dirty)

ARGOCD_NAMESPACE := argocd

WHITE := \033[0m
BLUE := \033[36m

.PHONY: build build-% push push-% deploy deploy-% help

list: ## display a list of all specific airways
	@for api in $(APIS); do \
		printf "$(BLUE)$$api$(WHITE)\n"; \
	done

build: build-argocd-aoa ## build all airways an flights
	@echo "Building all flights..."
	@for api in $(APIS); do \
		GOOS=wasip1 GOARCH=wasm go build -C ./$$api -o ../.out/flight-$$api.wasm ./flight; \
		GOOS=wasip1 GOARCH=wasm go build -C ./$$api -o ../.out/airway-$$api.wasm ./airway; \
	done

build-%: ## build a specific airway and flight
	GOOS=wasip1 GOARCH=wasm go build -C ./$* -o ../.out/flight-$*.wasm ./flight
	GOOS=wasip1 GOARCH=wasm go build -C ./$* -o ../.out/airway-$*.wasm ./airway

build-argocd-aoa:
	GOOS=wasip1 GOARCH=wasm go build -C ./argocd-aoa -o ../.out/flight-argocd-aoa-installer.wasm ./installer
	GOOS=wasip1 GOARCH=wasm go build -C ./argocd-aoa -o ../.out/flight-argocd-aoa.wasm ./app-of-apps

push: build push-argocd-aoa ## build and push all airways and flights
	@for api in $(APIS); do \
		yoke push .out/flight-$$api.wasm oci://$(REGISTRY)/$(REPOSITORY)/flight-$$api:$(TAG); \
		yoke push .out/airway-$$api.wasm oci://$(REGISTRY)/$(REPOSITORY)/airway-$$api:$(TAG); \
	done

push-%: build-% ## build and push a specific airway and flight
	yoke push .out/flight-$*.wasm oci://$(REGISTRY)/$(REPOSITORY)/flight-$*:$(TAG)
	yoke push .out/airway-$*.wasm oci://$(REGISTRY)/$(REPOSITORY)/airway-$*:$(TAG)

push-argocd-aoa: build-argocd-aoa ## build and push the argocd app of apps flight
	yoke push .out/flight-argocd-aoa-installer.wasm oci://$(REGISTRY)/$(REPOSITORY)/flight-argocd-aoa-installer:$(TAG)
	yoke push .out/flight-argocd-aoa.wasm oci://$(REGISTRY)/$(REPOSITORY)/flight-argocd-aoa:$(TAG)

deploy: ## all airways to the cluster
	@for api in $(APIS); do \
		yoke takeoff --wait=30s airway-$$api oci://$(REGISTRY)/$(REPOSITORY)/airway-$$api:$(TAG) -- --flight oci://$(REGISTRY)/$(REPOSITORY)/flight-$$api:$(TAG); \
	done

deploy-%: ## deploy a specific airway
	yoke takeoff --wait=30s airway-$* oci://$(REGISTRY)/$(REPOSITORY)/airway-$*:$(TAG) -- --flight oci://$(REGISTRY)/$(REPOSITORY)/flight-$*:$(TAG)

deploy-app-of-apps: install-yokecd ## deploys app of apps (and argocd with yokecd)
	yoke takeoff -wait 30s --namespace $(ARGOCD_NAMESPACE) app-of-apps oci://$(REGISTRY)/$(REPOSITORY)/flight-argocd-aoa-installer:$(TAG)

undeploy: ## remove all airways and flights in this repo
	@for api in $(APIS); do \
		yoke mayday airway-$$api; \
	done

install: ## the yoke-atc into the cluster
	yoke takeoff -wait 30s --create-namespace --namespace atc atc oci://ghcr.io/yokecd/atc-installer:$(YOKE_ATC_VERISON)

install-yokecd: ## installs argocd with yokecd plugin
	@echo "version: $(YOKECD_VERSION)" | yoke takeoff --create-namespace --namespace $(ARGOCD_NAMESPACE) yokecd oci://ghcr.io/yokecd/yokecd-installer:$(YOKECD_VERSION)

uninstall: undeploy ## the yoke-atc from the cluster (cleans up airways first)
	yoke mayday --namespace atc atc

help:
	@awk 'BEGIN {FS = ":.*?#"} /^[a-zA-Z0-9%_-]+:.*?#/ { printf "  $(BLUE)%-15s$(WHITE) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
