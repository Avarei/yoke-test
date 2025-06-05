APIS := cluster mylist
REGISTRY := ghcr.io
REPOSITORY := avarei/yoke-test
TAG := $(shell git describe --tags --always --dirty)

build:
	@echo "Building all flights..."
	@for api in $(APIS); do \
		GOOS=wasip1 GOARCH=wasm go build -C ./$$api -o ../.out/flight-$$api.wasm ./flight; \
		GOOS=wasip1 GOARCH=wasm go build -C ./$$api -o ../.out/airway-$$api.wasm ./airway; \
	done

build-%:
	GOOS=wasip1 GOARCH=wasm go build -C ./$* -o ../.out/flight-$*.wasm ./flight
	GOOS=wasip1 GOARCH=wasm go build -C ./$* -o ../.out/airway-$*.wasm ./airway

push: build
	@for api in $(APIS); do \
		yoke push .out/flight-$$api.wasm oci://$(REGISTRY)/$(REPOSITORY)/flight-$$api:$(TAG); \
		yoke push .out/airway-$$api.wasm oci://$(REGISTRY)/$(REPOSITORY)/airway-$$api:$(TAG); \
	done

push-%: build-%
	yoke push .out/flight-$*.wasm oci://$(REGISTRY)/$(REPOSITORY)/flight-$*:$(TAG)
	yoke push .out/airway-$*.wasm oci://$(REGISTRY)/$(REPOSITORY)/airway-$*:$(TAG)

deploy:
	@for api in $(APIS); do \
		yoke takeoff --wait=30s airway-$$api oci://$(REGISTRY)/$(REPOSITORY)/airway-$$api:$(TAG) -- --flight oci://$(REGISTRY)/$(REPOSITORY)/flight-$$api:$(TAG); \
	done

deploy-%:
	yoke takeoff --wait=30s airway-$* oci://$(REGISTRY)/$(REPOSITORY)/airway-$*:$(TAG) -- --flight oci://$(REGISTRY)/$(REPOSITORY)/flight-$*:$(TAG)

dev-deploy: build
	@for api in $(APIS); do \
		yoke takeoff --wait=30s airway-$$api oci://$(REGISTRY)/$(REPOSITORY)/airway-$$api:$(TAG) -- --flight .out/flight-$$api.wasm; \
	done
dev-deploy-%: build-%
	yoke takeoff --wait=30s airway-$* oci://$(REGISTRY)/$(REPOSITORY)/airway-$*:$(TAG) -- --flight .out/flight-$*.wasm
