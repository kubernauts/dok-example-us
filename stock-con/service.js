'use strict';

const express = require('express');
const avg = require('moving-average');
const http = require("http");
const https = require("https");

const SERVICE_PORT = 9898;
const app = express();

var DOK_STOCKGEN_HOSTNAME = "stock-gen"
var DOK_STOCKGEN_PORT = 80
var period = 60 * 1000; // 60 seconds 
var ma = avg(period);

app.get('/average/:symbol', function (req, res) {
    var symbol = req.params.symbol;
    console.info('\n===\nTargeting symbol ' + symbol + ' over the past ' + period/1000 + ' seconds');
    httpGetJSON(DOK_STOCKGEN_HOSTNAME, DOK_STOCKGEN_PORT, '/stockdata', function (e, stocks) {
        for (var i = 0, len = stocks.length; i < len; i++) {
            var stock = stocks[i]
            if (symbol == stock.symbol) {
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
                res.json(result)
                res.end();
                return
            }
            // res.status(404).end();
        }
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