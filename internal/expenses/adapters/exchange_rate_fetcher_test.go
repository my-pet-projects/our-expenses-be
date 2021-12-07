package adapters_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/adapters"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/domain"
)

func TestExchangeRateFetcher_NewExchangeRateFetcher_ReturnsInstance(t *testing.T) {
	t.Parallel()
	// Arrange
	config := adapters.ExchangeRateFetcherConfig{}

	// Act
	result := adapters.NewExchangeRateFetcher(config)

	// Assert
	assert.NotNil(t, result)
}

func TestFetch_FailedRequestWith401_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	datesRange := []time.Time{
		time.Now(),
		time.Now().Add(-10 * 24 * time.Hour),
		time.Now().Add(-2 * 24 * time.Hour),
	}
	expected := "unauthorized"
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		w.Write([]byte(expected))
	}))
	defer svr.Close()
	config := adapters.ExchangeRateFetcherConfig{
		URL:    svr.URL,
		APIKey: "key",
	}

	// SUT
	sut := adapters.NewExchangeRateFetcher(config)

	// Act
	res, resErr := sut.Fetch(context.Background(), datesRange)

	// Assert
	assert.Nil(t, res)
	assert.NotNil(t, resErr)
}

func TestFetch_FailedRequestWith500_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	datesRange := []time.Time{
		time.Now(),
		time.Now().Add(-10 * 24 * time.Hour),
		time.Now().Add(-2 * 24 * time.Hour),
	}
	expected := "server error"
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte(expected))
	}))
	defer svr.Close()
	config := adapters.ExchangeRateFetcherConfig{
		URL:    svr.URL,
		APIKey: "key",
	}

	// SUT
	sut := adapters.NewExchangeRateFetcher(config)

	// Act
	res, resErr := sut.Fetch(context.Background(), datesRange)

	// Assert
	assert.Nil(t, res)
	assert.NotNil(t, resErr)
}

func TestFetch_ResponseFails_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	datesRange := []time.Time{
		time.Now(),
		time.Now().Add(-10 * 24 * time.Hour),
		time.Now().Add(-2 * 24 * time.Hour),
	}
	config := adapters.ExchangeRateFetcherConfig{
		URL:    "url",
		APIKey: "key",
	}

	// SUT
	sut := adapters.NewExchangeRateFetcher(config)

	// Act
	res, resErr := sut.Fetch(context.Background(), datesRange)

	// Assert
	assert.Nil(t, res)
	assert.NotNil(t, resErr)
}

func TestFetch_FailedToReadRequestBody_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	datesRange := []time.Time{
		time.Now(),
		time.Now().Add(-10 * 24 * time.Hour),
		time.Now().Add(-2 * 24 * time.Hour),
	}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// invalid content length will fail body read
		w.Header().Set("Content-Length", "1")
	}))
	defer svr.Close()
	config := adapters.ExchangeRateFetcherConfig{
		URL:    svr.URL,
		APIKey: "key",
	}

	// SUT
	sut := adapters.NewExchangeRateFetcher(config)

	// Act
	res, resErr := sut.Fetch(context.Background(), datesRange)

	// Assert
	assert.Nil(t, res)
	assert.NotNil(t, resErr)
}

func TestFetch_ResponseDecodeFails_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	datesRange := []time.Time{
		time.Now(),
		time.Now().Add(-10 * 24 * time.Hour),
		time.Now().Add(-2 * 24 * time.Hour),
	}
	expected := "unexpected response"
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(expected))
		w.Header()
	}))
	defer svr.Close()
	config := adapters.ExchangeRateFetcherConfig{
		URL:    svr.URL,
		APIKey: "key",
	}

	// SUT
	sut := adapters.NewExchangeRateFetcher(config)

	// Act
	res, resErr := sut.Fetch(context.Background(), datesRange)

	// Assert
	assert.Nil(t, res)
	assert.NotNil(t, resErr)
}

func TestFetch_Response_ValidRate_ReturnsExchangeRates(t *testing.T) {
	t.Parallel()
	// Arrange
	datesRange := []time.Time{
		time.Now(),
	}
	expected := "{ \"base\":\"USD\", \"rates\": { \"EUR\": 1.123 } }"
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(expected))
	}))
	defer svr.Close()
	config := adapters.ExchangeRateFetcherConfig{
		URL:    svr.URL,
		APIKey: "key",
	}

	// SUT
	sut := adapters.NewExchangeRateFetcher(config)

	// Act
	res, resErr := sut.Fetch(context.Background(), datesRange)

	// Assert
	assert.NotNil(t, res)
	assert.Equal(t, domain.Currency("USD"), res[0].BaseCurrency())
	assert.Equal(t, map[domain.Currency]decimal.Decimal{domain.Currency("EUR"): decimal.NewFromFloat(1.123)},
		res[0].Rates())
	assert.Nil(t, resErr)
}

func TestFetch_Response_InvalidRate_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	datesRange := []time.Time{
		time.Now(),
	}
	expected := "{ \"base\":\"USD\", \"rates\": { } }"
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(expected))
	}))
	defer svr.Close()
	config := adapters.ExchangeRateFetcherConfig{
		URL:    svr.URL,
		APIKey: "key",
	}

	// SUT
	sut := adapters.NewExchangeRateFetcher(config)

	// Act
	res, resErr := sut.Fetch(context.Background(), datesRange)

	// Assert
	assert.Nil(t, res)
	assert.NotNil(t, resErr)
}
