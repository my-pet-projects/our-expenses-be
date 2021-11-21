package adapters

import "os"

// ExchangeRateFetcherConfig hold exchange rate fetcher config.
type ExchangeRateFetcherConfig struct {
	Url    string
	ApiKey string
}

// NewExchangeRateFetcherConfig instantiates fetcher config.
func NewExchangeRateFetcherConfig() ExchangeRateFetcherConfig {
	return ExchangeRateFetcherConfig{
		Url:    os.Getenv("EXCHANGE_RATE_FETCHER_URL"),
		ApiKey: os.Getenv("EXCHANGE_RATE_FETCHER_KEY"),
	}
}
