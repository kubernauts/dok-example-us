# Developing on Kubernetes

This repository contains two microservices that make up an application:

- One microservice, the so called [stock generator](stock-gen/) is written in Go and serves randomized stock data randomly via the HTTP endpoint `stockdata/`.
- A second microservice written in Node.js that consumes above stock data and performs 1. aggregation like average over the past year of a certain stock, 2. provides recommendation like sell/hold/buy, and 3. transform the data into plots (over time graphical representation).
