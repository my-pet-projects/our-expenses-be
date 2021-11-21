package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/tracer"
)

type exchangeRateResponse struct {
	Base  string             `json:"base"`
	Rates map[string]float32 `json:"rates"`
}

// ExchangeRateFetcher represent exchange rate fetcher.
type ExchangeRateFetcher struct {
	config ExchangeRateFetcherConfig
}

// ExchangeRateFetcherInterface defines a contract to fetch rates.
type ExchangeRateFetcherInterface interface {
	Fetch(ctx context.Context, dates []time.Time) ([]domain.ExchangeRate, error)
}

// NewExchangeRateFetcher returns a ExchangeRateFetcher.
func NewExchangeRateFetcher(config ExchangeRateFetcherConfig) ExchangeRateFetcherInterface {
	return ExchangeRateFetcher{
		config: config,
	}
}

// Fetch fetches exchange rate data.
func (f ExchangeRateFetcher) Fetch(ctx context.Context, dates []time.Time) ([]domain.ExchangeRate, error) {
	_, span := tracer.NewSpan(ctx, "fetch exchange rates from the provider")
	defer span.End()

	exchangeRates := make([]domain.ExchangeRate, 0)
	for _, date := range dates {
		url := fmt.Sprintf("%s/%s.json?app_id=%s", f.config.Url, date.Format("2006-01-02"), f.config.ApiKey)
		resp, respErr := http.Get(url)
		if respErr != nil {
			tracer.AddSpanError(span, respErr)
			return nil, errors.Wrap(respErr, "failed response")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			buf := new(bytes.Buffer)
			_, copyErr := io.Copy(buf, resp.Body)
			if copyErr != nil {
				return nil, errors.Wrap(copyErr, "buffer copy")
			}
			buf.ReadFrom(resp.Body)
			respString := buf.String()

			if resp.StatusCode == http.StatusUnauthorized {
				return nil, errors.New(fmt.Sprintf("failed to authorize: %s", respString))
			}

			return nil, errors.New(fmt.Sprintf("unsuccessful reply: %s", respString))
		}

		var response exchangeRateResponse
		if jsonErr := json.NewDecoder(resp.Body).Decode(&response); jsonErr != nil {
			tracer.AddSpanError(span, jsonErr)
			return nil, errors.Wrap(jsonErr, "response decode")
		}

		exchangeRate := domain.NewExchageRate(date, response.Base, response.Rates)
		exchangeRates = append(exchangeRates, exchangeRate)
	}
	return exchangeRates, nil
}
