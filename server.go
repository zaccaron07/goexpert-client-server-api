package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	apiBaseUrl     = "https://economia.awesomeapi.com.br/json"
	requestTimeout = 200 * time.Millisecond
)

type ExchangeRateResponse struct {
	USDBRL ExchangeRateDetails `json:"USDBRL"`
}

type ExchangeRateDetails struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type ExchangeRateClientResponse struct {
	Bid string `json:"bid"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", exchangeRateHandler)
	http.ListenAndServe("127.0.0.1:8080", mux)
}

func exchangeRateHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	resp, err := fetchExchangeRate(ctx)

	if err != nil {
		fmt.Printf("Error fetching exchange rate: %v", err)
		http.Error(w, "Failed to fetch exchange rate", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(&resp)
}

func fetchExchangeRate(ctx context.Context) (*ExchangeRateDetails, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", apiBaseUrl+"/last/USD-BRL", nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("request timed out")
		}
		return nil, fmt.Errorf("performing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	var exchangeRateResponse ExchangeRateResponse
	err = json.Unmarshal(body, &exchangeRateResponse)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response: %w", err)
	}

	return &exchangeRateResponse.USDBRL, nil
}
