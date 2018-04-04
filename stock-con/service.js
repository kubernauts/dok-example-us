'use strict';

const express = require('express');
const avg = require('moving-average');
const http = require("http");
const https = require("https");
var HashMap = require('hashmap');

const SERVICE_PORT = 9898;
const app = express();

var DOK_STOCKGEN_HOSTNAME = "stock-gen";
var DOK_STOCKGEN_PORT = 80;
var period = 60 * 1000; // 60 seconds 
var sym2Avg = new HashMap();
var lastseensym = "";

// the average endpoint
app.get('/average/:symbol', function (req, res) {
    var symbol = req.params.symbol;
    console.info('\n===\nTargeting symbol [' + symbol + '] over the past ' + period/1000 + ' seconds');
    httpGetJSON(DOK_STOCKGEN_HOSTNAME, DOK_STOCKGEN_PORT, '/stockdata', function (e, stocks) {
        if (stocks == null){
            res.status(500).end();
            return
        }
        for (var i = 0, len = stocks.length; i < len; i++) {
            var stock = stocks[i];
            var ma = sym2Avg.get(stock.symbol);
            if (symbol == stock.symbol) {
                lastseensym = symbol;
                ma.push(Date.now(), stock.value);
                console.info('Stock: ' + stock.symbol + ' @ ' + stock.value);
                console.info('-> moving average:', ma.movingAverage());
                console.info('-> forecast:', ma.forecast());
                var result = {
                    symbol: stock.symbol,
                    current: stock.value,
                    moving_average: ma.movingAverage(),
                    forecast: ma.forecast()
                }
                res.json(result);
                res.end();
                return
            }
        }
        res.status(404).end();
    });
});

// the health check endpoint
app.get('/healthz', function (req, res) {
    var result = {
        numsymbols: sym2Avg.size,
        lastseen: lastseensym,
        test: "hello Paris"
    }
    res.json(result);
    res.end();
});

function httpGetJSON(host, port, path, callback) {
    return http.get({
        host: host,
        port: port,
        path: path,
        json: true
    }, function (response) {
        var body = '';
        response.setEncoding('utf8');
        response.on('data', function (d) {
            body += d;
        });
        response.on('end', function () {
            try {
                var data = JSON.parse(body);
            } catch (err) {
                console.error('Can\'t parse JSON payload: ', err);
                return callback(err);
            }
            callback(null, data);
        });
    }).on('error', function (err) {
        console.error('Can\'t GET data from ' + host + ':' + port + path + ' due to:', err.message);
        callback(err);
    });
}

function init(){
    if (process.env.DOK_STOCKGEN_HOSTNAME) {
        DOK_STOCKGEN_HOSTNAME = process.env.DOK_STOCKGEN_HOSTNAME
    }
    if (process.env.DOK_STOCKGEN_PORT) {
        DOK_STOCKGEN_PORT = parseInt(process.env.DOK_STOCKGEN_PORT, 10);
    }
    // read in all stock symbols and create a moving average object for each:
    httpGetJSON(DOK_STOCKGEN_HOSTNAME, DOK_STOCKGEN_PORT, '/stockdata', function (e, stocks) {
        if (stocks == null) {
            return
        }
        for (var i = 0, len = stocks.length; i < len; i++) {
            var stock = stocks[i];
            var ma = avg(period);
            console.info('Creating moving average for symbol ' + stock.symbol);
            sym2Avg.set(stock.symbol, ma);
        }
    });
}

init()
app.listen(SERVICE_PORT);
console.info('DoK stock market consumer service running on port ' + SERVICE_PORT);