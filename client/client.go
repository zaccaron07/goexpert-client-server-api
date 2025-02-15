package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	apiBaseUrl     = "http://localhost:8080"
	requestTimeout = 300 * time.Millisecond
)

type ExchangeRateDetails struct {
	Bid string `json:"bid"`
}

func main() {
	exchangeRateDetails, err := fetchExchangeRate()
	if err != nil {
		fmt.Printf("error fetching exchange rate: %v\n", err)
		return
	}

	err = writeExchangeRateToFile(exchangeRateDetails.Bid)
	if err != nil {
		fmt.Printf("error writing exchange rate to file: %v\n", err)
		return
	}
}

func fetchExchangeRate() (*ExchangeRateDetails, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", apiBaseUrl+"/cotacao", nil)
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

	var exchangeRateDetails ExchangeRateDetails
	err = json.Unmarshal(body, &exchangeRateDetails)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response: %w", err)
	}
	return &exchangeRateDetails, nil
}

func writeExchangeRateToFile(bid string) error {
	f, err := os.Create("cotacao.txt")
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("Dólar: %s", bid))

	if err != nil {
		return fmt.Errorf("writing to file: %w", err)
	}
	return nil
}
