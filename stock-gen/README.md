# The stock generator

Set it up locally:

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
   http localhost:9876/stockdata | jq .[].symbol,.[].value ; 
   sleep 1 ; 
  done
```
