package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Stock represents a company listed at a stock market.
type Stock struct {
	Symbol   string  `json:"symbol"`
	Value    float32 `json:"current"`
	Currency string  `json:"currency"`
}

func main() {
	dummy := Stock{
		Symbol:   "GOOG",
		Value:    100,
		Currency: "USD",
	}
	http.HandleFunc("/stockdata", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(dummy)
		if err != nil {
			fmt.Printf("HTTP %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	port := os.Getenv("DOK_STOCKGEN_PORT")
	if port == "" {
		port = "80"
	}
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Error launching service %v\n", err)
	}
}
