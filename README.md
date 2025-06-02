# Yoke Tests

In this repository I want to experiment with [yoke](https://yokecd.github.io/docs/).

Backend is just the example from the [ATC Docs](https://yokecd.github.io/docs/airtrafficcontroller/atc/)

## Getting Started

install yoke see [flake.nix](./flake.nix) or `go install github.com/yokecd/yoke/cmd/yoke@latest`


## Install Yoke ATC

### For now
yoke takeoff -wait 30s --create-namespace --namespace atc atc oci://ghcr.io/yokecd/atc-installer:latest

### In the Future
TODO:

* Deploy ArgoCD using Yoke (with yokecd plugin)
* Deploy YokeATC using ArgoCD

## Build Flight
```sh
GOOS=wasip1 GOARCH=wasm go build -C ./cluster -o ../.out/flight-cluster.wasm ./flight
yoke push .out/flight-cluster-v1alpha1.wasm oci://ghcr.io/avarei/yoke-test/flight-cluster:v0.0.0-dirty
```

## Build Airway
```sh
# don't forget to update the flight reference
GOOS=wasip1 GOARCH=wasm go build -C cluster -o ../.out/airway-cluster.wasm ./airway
```

## Deploy it
```sh
yoke takeoff --wait 30s airway-cluster .out/airway-cluster.wasm -- --flight=oci://ghcr.io/avarei/yoke-test/flight-cluster:v0.0.0-dirty
kubectl apply -f - <<EOF
apiVersion: example.com/v1alpha1
kind: Cluster
metadata:
  name: sunshine
spec:
  type: vCluster
EOF
```
