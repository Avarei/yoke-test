# Yoke Tests

In this repository I want to experiment with [yoke](https://yokecd.github.io/docs/).

Backend is just the example from the [ATC Docs](https://yokecd.github.io/docs/airtrafficcontroller/atc/)

## Getting Started

install yoke see [flake.nix](./flake.nix) or `go install github.com/yokecd/yoke/cmd/yoke@latest`


## Install Yoke ATC

### For now
make install

### In the Future
TODO:

* Deploy ArgoCD using Yoke (with yokecd plugin)
* Deploy YokeATC using ArgoCD

## Build Flight and Airway

To build all flights and airways run
```sh
make build
```

To build a specific flight and airway find it using `make list` and build it using e.g. `make build-cluster`

## Push it
to push an airway and flight to ghcr.io run `make push` or for a specific airway `make push-<target>`

## Deploy it

To deploy it to the current cluster use `make deploy` or `make deploy-<target>` to deploy a specific api

### Cluster

Create a vCluster

```sh
make deploy-cluster

kubectl apply -f - <<EOF
apiVersion: example.com/v1alpha1
kind: Cluster
metadata:
  name: sunshine
spec:
  type: vCluster
EOF
```

### MyList

Lets create a conflict were two CRs both want to manage the same resource

```sh
make deploy-mylist

kubectl apply -f - <<EOF
---
apiVersion: example.com/v1alpha1
kind: MyList
metadata:
  name: conflict-a
spec:
  items:
    - apiVersion: v1
      kind: ConfigMap
      metadata:
        name: first
      data:
        foo: bar
---
apiVersion: example.com/v1alpha1
kind: MyList
metadata:
  name: conflict-b
spec:
  items:
    - apiVersion: v1
      kind: ConfigMap
      metadata:
        name: first
      data:
        bar: baz
EOF
```