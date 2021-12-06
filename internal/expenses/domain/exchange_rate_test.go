package domain_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/domain"
)

func TestExchangeRate_NewExchangeRate_InstantiatesAndReturnsProperties(t *testing.T) {
	t.Parallel()
	// Arrange
	date := time.Now()
	base := "EUR"
	rates := map[string]float64{"USD": 123.45}

	// Act
	res, resErr := domain.NewExchageRate(date, base, rates)

	// Assert
	assert.NotNil(t, res)
	assert.Nil(t, resErr)
	assert.Equal(t, domain.Currency(base), res.BaseCurrency())
	assert.Equal(t, date, res.Date())
	assert.Equal(t, decimal.NewFromFloat32(123.45), res.Rates()[domain.Currency("USD")])
}

func TestExchangeRate_EmptyBaseCurrency_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	date := time.Now()
	base := ""
	rates := map[string]float64{"EUR": 2}

	// Act
	res, resErr := domain.NewExchageRate(date, base, rates)

	// Assert
	assert.Nil(t, res)
	assert.NotNil(t, resErr)
}

func TestExchangeRate_EmptyRates_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	date := time.Now()
	base := "EUR"
	rates := map[string]float64{}

	// Act
	res, resErr := domain.NewExchageRate(date, base, rates)

	// Assert
	assert.Nil(t, res)
	assert.NotNil(t, resErr)
}

func TestChangeBaseCurrency_ToSameCurrency_ReturnsSameRates(t *testing.T) {
	t.Parallel()
	// Arrange
	date := time.Now()
	baseCurrency := "USD"
	rates := map[string]float64{
		"HRK": 6.2632,
		"GBP": 0.720181,
		"EUR": 0.82725,
	}

	// SUT
	sut, _ := domain.NewExchageRate(date, baseCurrency, rates)

	// Act
	result := sut.ChangeBaseCurrency("USD")

	// Assert
	assert.Equal(t, date, result.Date())
	assert.Equal(t, domain.Currency(baseCurrency), result.BaseCurrency())
	assert.Len(t, result.Rates(), len(rates))
	assert.Equal(t, fmt.Sprintf("%v", rates["HRK"]), result.Rates()[domain.Currency("HRK")].String())
	assert.Equal(t, fmt.Sprintf("%v", rates["GBP"]), result.Rates()[domain.Currency("GBP")].String())
	assert.Equal(t, fmt.Sprintf("%v", rates["EUR"]), result.Rates()[domain.Currency("EUR")].String())
}

func TestChangeBaseCurrency_ToDiffCurrency_RecalculatesRates(t *testing.T) {
	t.Parallel()
	// Arrange
	date := time.Now()
	baseCurrency := "USD"
	rates := map[string]float64{
		"HRK": 6.2632,
		"GBP": 0.720181,
		"EUR": 0.82725,
	}

	// SUT
	sut, _ := domain.NewExchageRate(date, baseCurrency, rates)

	// Act
	result := sut.ChangeBaseCurrency("EUR")

	// Assert
	hrkRate := decimal.NewFromFloat(rates["HRK"]).Div(decimal.NewFromFloat(rates["EUR"]))
	gbpRate := decimal.NewFromFloat(rates["GBP"]).Div(decimal.NewFromFloat(rates["EUR"]))
	usdRate := decimal.NewFromFloat(1).Div(decimal.NewFromFloat(rates["EUR"]))

	assert.Equal(t, date, result.Date())
	assert.Equal(t, domain.Currency("EUR"), result.BaseCurrency())
	assert.Len(t, result.Rates(), len(rates))
	assert.Equal(t, hrkRate, result.Rates()[domain.Currency("HRK")])
	assert.Equal(t, gbpRate, result.Rates()[domain.Currency("GBP")])
	assert.Equal(t, usdRate, result.Rates()[domain.Currency("USD")])
	assert.NotContains(t, result.Rates(), domain.Currency("EUR"))
}
