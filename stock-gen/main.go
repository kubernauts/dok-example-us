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
	"time"
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
		reallyrandom := rand.New(rand.NewSource(time.Now().UnixNano()))
		idx := reallyrandom.Intn(len(stocks))
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
	fmt.Printf("using symbols from %v\n", symfile)
	c, err := os.Open(symfile)
	if err != nil {
		return stocks, err
	}
	csvreader := csv.NewReader(bufio.NewReader(c))
	for {
		record, error := csvreader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			return stocks, err
		}
		v, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return stocks, err
		}
		stock := Stock{
			Symbol:   record[0],
			Value:    v,
			Currency: "USD",
		}
		stocks = append(stocks, stock)
		fmt.Printf("adding symbol %v\n", stock)
	}
	return stocks, nil
}
