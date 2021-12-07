package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/tracer"
)

type exchRateResponse struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
}

// ExchangeRateFetcher represent exchange rate fetcher.
type ExchangeRateFetcher struct {
	config ExchangeRateFetcherConfig
}

// ExchangeRateFetcherInterface defines a contract to fetch rates.
type ExchangeRateFetcherInterface interface {
	Fetch(ctx context.Context, dates []time.Time) ([]domain.ExchangeRates, error)
}

// NewExchangeRateFetcher returns a ExchangeRateFetcher.
func NewExchangeRateFetcher(config ExchangeRateFetcherConfig) ExchangeRateFetcher {
	return ExchangeRateFetcher{
		config: config,
	}
}

// Fetch fetches exchange rate data.
func (f ExchangeRateFetcher) Fetch(ctx context.Context, dates []time.Time) ([]domain.ExchangeRates, error) {
	_, span := tracer.NewSpan(ctx, "fetch exchange rates from the provider")
	defer span.End()

	exchRates := make([]domain.ExchangeRates, 0)
	for _, date := range dates {
		url := fmt.Sprintf("%s/%s.json?app_id=%s", f.config.URL, date.Format("2006-01-02"), f.config.APIKey)
		resp, respErr := http.Get(url)
		if respErr != nil {
			tracer.AddSpanError(span, respErr)
			return nil, errors.Wrap(respErr, "failed response")
		}

		body, bodyErr := ioutil.ReadAll(resp.Body)
		if bodyErr != nil {
			return nil, errors.Wrap(bodyErr, "response body read")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			if resp.StatusCode == http.StatusUnauthorized {
				return nil, fmt.Errorf("failed to authorize: %s", string(body))
			}

			return nil, fmt.Errorf("unsuccessful reply: %s", string(body))
		}

		var response exchRateResponse
		if jsonErr := json.Unmarshal(body, &response); jsonErr != nil {
			tracer.AddSpanError(span, jsonErr)
			return nil, errors.Wrap(jsonErr, "response decode")
		}

		exchRate, exchRateErr := domain.NewExchageRate(date, response.Base, response.Rates)
		if exchRateErr != nil {
			return nil, errors.Wrap(exchRateErr, "invalid exchange rate")
		}
		exchRates = append(exchRates, *exchRate)
	}
	return exchRates, nil
}
