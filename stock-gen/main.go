package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

// Stock represents a company listed at a stock market.
type Stock struct {
	Symbol   string  `json:"symbol"`
	Value    float64 `json:"value"`
	Currency string  `json:"currency"`
}

func main() {
	stocks, err := loadsym()
	if err != nil {
		fmt.Printf("Couldn't load symbols due to %v\n", err)
	}
	http.HandleFunc("/stockdata", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		idx := rand.Intn(len(stocks))
		err = json.NewEncoder(w).Encode(stocks[idx])
		if err != nil {
			fmt.Printf("HTTP %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	port := os.Getenv("DOK_STOCKGEN_PORT")
	if port == "" {
		port = "80"
	}
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Error launching service %v\n", err)
	}
}

func loadsym() (stocks []Stock, err error) {
	symfile := os.Getenv("DOK_STOCKGEN_SYMBOLS")
	if symfile == "" {
		symfile = "./symbols.csv"
	}
	c, err := os.Open(symfile)
	if err != nil {
		return stocks, err
	}
	for {
		record, error := csv.NewReader(bufio.NewReader(c)).Read()
		if error == io.EOF {
			break
		} else if error != nil {
			return stocks, err
		}
		v, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return stocks, err
		}
		stocks = append(stocks, Stock{
			Symbol:   record[0],
			Value:    v,
			Currency: "USD",
		})
	}
	return stocks, nil
}
