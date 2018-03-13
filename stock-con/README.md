# The stock consumer

The stock consumer (`stock-con`) consumes data from the [stock generator](../stock-gen/) and performs various calculations on different endpoints on it.

- [Local development](#local-development)
- [Endpoints](#endpoints)
  - [Basic stats](#basic-stats)
  - [Health check](#health-check)
- [Containerized deployment](#containerized-deployment)
- [Grafana dashboard](#grafana-dashboard)

## Local development

Set it up locally (assumes `npm` installed, tested with version `5.7.1`):

```bash
$ npm install express
$ npm install
```

Build and run in one terminal (note: make sure you've got the [stock generator](../stock-gen/) service running as well otherwise it won't work):

```bash
$ DOK_STOCKGEN_HOSTNAME=localhost DOK_STOCKGEN_PORT=9876 npm start
```

In another terminal (requires [http](https://httpie.org/) installed, otherwise use `curl`):

```bash
$ while true; do \
   http localhost:9898/average/NYSE:RHT ; \
   sleep 2 ; \
  done
```

## Endpoints

### Basic stats

Usage:

```
/average/$SYMBOL
```

Behavior: If `$SYMBOL` is a valid and known one a response like below, otherwise a `404`:

```bash
HTTP/1.1 200 OK
Connection: keep-alive
Content-Length: 114
Content-Type: application/json; charset=utf-8
Date: Tue, 13 Mar 2018 17:34:01 GMT
ETag: W/"72-vkOrThRnV1hPbp7xrmGe17qHb4I"
X-Powered-By: Express

{
    "current": 598.8333836078177,
    "forecast": 437.20974460845946,
    "moving_average": 599.5866511312555,
    "symbol": "NYSE:RHT"
}
```

### Health check

Usage:

```
/healthz
```

Behavior: A `200` response like below, listing the number of known symbols and the symbol that has been queried for most recently (ignoring unknown symbols):

```bash
HTTP/1.1 200 OK
Connection: keep-alive
Content-Length: 38
Content-Type: application/json; charset=utf-8
Date: Tue, 13 Mar 2018 17:34:43 GMT
ETag: W/"26-rNuCWemdLty1c1TH7oAt5mn3q2U"
X-Powered-By: Express

{
    "lastseen": "NYSE:RHT",
    "numsymbols": 4
}
```

## Containerized deployment

The `stock-con` app container image is defined in [Dockerfile](./Dockerfile) and publicly available via [quay.io/mhausenblas/stock-con](https://quay.io/repository/mhausenblas/stock-con).

To deploy it as an Kubernetes application, use the [provided manifest](./app.yaml) like so:

```bash
$ kubectl create -f app.yaml
$ kubectl get deploy,svc,po
```

## Grafana dashboard

Note: the following is WIP!

The Grafana app requires an [Ingress controller](https://github.com/kubernetes/ingress-nginx/blob/master/deploy/README.md).

```bash
$ kubectl create -f grafana-app.yaml
```

Note: see also the [NGINX Ingress user guide](https://github.com/kubernetes/ingress-nginx/tree/master/docs/user-guide) for more options.
