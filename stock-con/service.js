'use strict';

const express = require('express');
const http = require("http");
const https = require("https");

const SERVICE_PORT = 9898;
const app = express();

var DOK_STOCKGEN_HOSTNAME = "stock-gen"
var DOK_STOCKGEN_PORT = 80

app.get('/average/:symbol/:period', function (req, res) {
    var symbol = req.params.symbol,
        period = req.params.period;
    var result = [];
    console.info('Calculating average stock price of symbol ' + symbol + ' over the past ' + period + ' ticks');
    httpGetJSON(DOK_STOCKGEN_HOSTNAME, DOK_STOCKGEN_PORT, '/stockdata', function (e, stocks) {
        for (var i = 0, len = stocks.length; i < len; i++) {
            var stock = stocks[i]
            console.info('Stock: ' + stock.symbol + ' @ ' + stock.value);
            if (symbol == stock.symbol) {
                result.push(stock)
            }
        }
        res.json(result)
        res.end();
        // res.status(404).end();
    });
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
        console.error('Got error on request: ', err.message);
        callback(err);
    });
}

if(process.env.DOK_STOCKGEN_HOSTNAME) {
    DOK_STOCKGEN_HOSTNAME = process.env.DOK_STOCKGEN_HOSTNAME
}

if (process.env.DOK_STOCKGEN_PORT) {
    DOK_STOCKGEN_PORT = parseInt(process.env.DOK_STOCKGEN_PORT, 10);
}

app.listen(SERVICE_PORT);
console.info('DoK stock market consumer service running on port ' + SERVICE_PORT);