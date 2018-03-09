# The stock consumer

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
$ http localhost:9898/average/NYSE:RHT/100
```
