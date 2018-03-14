# The stock generator

## Local development

Set it up locally (assumes Go 1.9 or above installed):

```bash
$ go get github.com/kubernauts/dok-example-us/stock-gen
```

Build and run in one terminal:

```bash
$ go install . && DOK_STOCKGEN_PORT=9876 DOK_STOCKGEN_CRASHEPOCH=20 stock-gen
```

In another terminal (requires `http` and `jq` installed):

```bash
$ while true; do \
   http localhost:9876/stockdata | jq .[].symbol,.[].value ; \
   sleep 1 ; \
  done
```

## Containerized deployment

The `stock-gen` app container image is defined in [Dockerfile](./Dockerfile) and publicly available via [quay.io/mhausenblas/stock-gen:0.3](https://quay.io/repository/mhausenblas/stock-gen).

To deploy it as an Kubernetes application, use the [provided manifest](./app.yaml) like so:

```bash
$ kubectl create -f app.yaml
$ kubectl get deploy,svc,po
```