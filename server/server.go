package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

const (
	apiBaseUrl            = "https://economia.awesomeapi.com.br/json"
	requestTimeout        = 200 * time.Millisecond
	insertDatabaseTimeout = 10 * time.Millisecond
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
	db := initializeDatabase()
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		exchangeRateHandler(w, r, db)
	})
	http.ListenAndServe(":8080", mux)
}

func initializeDatabase() *sql.DB {
	db, err := sql.Open("sqlite", "./exchange_rate.db")
	if err != nil {
		log.Fatal(err)
	}

	schema, err := os.ReadFile("./schema.sql")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func exchangeRateHandler(w http.ResponseWriter, _ *http.Request, db *sql.DB) {
	resp, err := fetchExchangeRate()

	if err != nil {
		fmt.Printf("Error fetching exchange rate: %v\n", err)
		http.Error(w, "Failed to fetch exchange rate", http.StatusInternalServerError)
		return
	}

	err = insertExchangeRate(db, resp)
	if err != nil {
		fmt.Printf("Error inserting exchange rate: %v\n", err)
		http.Error(w, "Failed to insert exchange rate", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(&resp)
}

func fetchExchangeRate() (*ExchangeRateDetails, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
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

func insertExchangeRate(db *sql.DB, exchangeRate *ExchangeRateDetails) error {
	ctx, cancel := context.WithTimeout(context.Background(), insertDatabaseTimeout)
	defer cancel()

	_, err := db.ExecContext(ctx, "INSERT INTO exchange_rate (code, codein, name, high, low, var_bid, pct_change, bid, ask, timestamp, create_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		exchangeRate.Code, exchangeRate.Codein, exchangeRate.Name, exchangeRate.High, exchangeRate.Low, exchangeRate.VarBid, exchangeRate.PctChange, exchangeRate.Bid, exchangeRate.Ask, exchangeRate.Timestamp, exchangeRate.CreateDate)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("request timed out")
		}
		return fmt.Errorf("inserting exchange rate: %w", err)
	}

	return nil
}
