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

const (
	minPrice = 2.0
	maxPrice = 1000.0
)

var tick int

func main() {
	// load the symbols and their initial values:
	stocks, err := loadsym()
	if err != nil {
		fmt.Printf("Couldn't load symbols due to %v\n", err)
	}
	// kick of random updates in the background:
	go func() {
		for {
			switch {
			// every 500 ticks we shuffle stock values randomly
			case tick%500 == 0:
				stocks = crash(stocks)
			// normally, we increase or decrease relative to a random peer:
			default:
				stocks = update(stocks)
			}
			tick++
			time.Sleep(1 * time.Second)
		}
	}()
	// HTTP API:
	http.HandleFunc("/stockdata", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(stocks)
		if err != nil {
			fmt.Printf("HTTP %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	// the service:
	port := os.Getenv("DOK_STOCKGEN_PORT")
	if port == "" {
		port = "80"
	}
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Error launching service %v\n", err)
	}
}

func update(in []Stock) (out []Stock) {
	reallyrandom := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, s := range in {
		cs := s
		idx := reallyrandom.Intn(len(in))
		cs.Value = genval(cs.Value, in[idx].Value)
		out = append(out, cs)
	}
	return out
}

func crash(in []Stock) (out []Stock) {
	for _, s := range in {
		cs := s
		reallyrandom := rand.New(rand.NewSource(time.Now().UnixNano()))
		cs.Value = minPrice + reallyrandom.Float64()*maxPrice
		out = append(out, cs)
	}
	return out
}

func genval(current, peer float64) (new float64) {
	switch {
	case current < peer:
		new = current + 1
	case current > peer:
		new = current - 1
	default:
		new = current
	}
	return new
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
