package adapters

import "os"

// ExchangeRateFetcherConfig hold exchange rate fetcher config.
type ExchangeRateFetcherConfig struct {
	URL    string
	APIKey string
}

// NewExchangeRateFetcherConfig instantiates fetcher config.
func NewExchangeRateFetcherConfig() ExchangeRateFetcherConfig {
	return ExchangeRateFetcherConfig{
		URL:    os.Getenv("EXCHANGE_RATE_FETCHER_URL"),
		APIKey: os.Getenv("EXCHANGE_RATE_FETCHER_KEY"),
	}
}
