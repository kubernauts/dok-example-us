# The stock consumer

## Local development

Set it up locally (assumes `npm` installed):

```bash
$ npm install express
$ npm install
```

Build and run in one terminal (note: make sure you've got the [stock generator](../stock-gen/) service running as well otherwise it won't work):

```bash
$ DOK_STOCKGEN_HOSTNAME=localhost DOK_STOCKGEN_PORT=9876 npm start
```

In another terminal (requires `http` installed):

```bash
$ while true; do \
   http localhost:9898/average/NYSE:RHT ; \
   sleep 2 ; \
  done
```

## Containerized deployment

The `stock-con` app container image is defined in [Dockerfile](./Dockerfile) and publicly available via [quay.io/mhausenblas/stock-con](https://quay.io/repository/mhausenblas/stock-con).

To deploy it as an Kubernetes application, use the [provided manifest](./app.yaml) like so:

```bash
$ kubectl create -f app.yaml
$ kubectl get deploy,svc,po
```

## Grafana dashboard

The following Grafana app requires an [Ingress controller](https://github.com/kubernetes/ingress-nginx/blob/master/deploy/README.md).

```bash
$ kubectl create -f grafana-app.yaml
```

Note: see also the [NGINX Ingress user guide](https://github.com/kubernetes/ingress-nginx/tree/master/docs/user-guide) for more options.
