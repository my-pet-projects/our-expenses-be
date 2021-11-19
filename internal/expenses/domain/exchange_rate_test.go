package domain_test

import (
	"testing"
	"time"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/domain"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestExchangeRate_NewExchangeRate_InstantiatesAndReturnsProperties(t *testing.T) {
	t.Parallel()
	// Arrange
	date := time.Now()
	base := "EUR"
	rates := map[string]float32{"USD": 123.45}

	// Act
	res := domain.NewExchageRate(date, base, rates)

	// Assert
	assert.NotNil(t, res)
	assert.Equal(t, domain.Currency(base), res.BaseCurrency())
	assert.Equal(t, date, res.Date())
	assert.Equal(t, decimal.NewFromFloat32(123.45), res.Rates()[domain.Currency("USD")])
}
