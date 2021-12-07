package domain

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNewExpense_ValidArgs_InstantiatesExpense(t *testing.T) {
	t.Parallel()
	// Arrange
	id := "id"
	price := 20.0
	quantity := 2.0
	currency := "RUB"
	category := Category{id: "catID"}
	trip := "trip"
	comment := "comment"
	date := time.Date(2021, 8, 2, 0, 0, 0, 0, time.UTC)
	createdBy := "createdBy"
	created := time.Date(2020, 8, 2, 0, 0, 0, 0, time.UTC)
	updatedBy := "updatedBy"
	updated := time.Date(2019, 8, 2, 0, 0, 0, 0, time.UTC)

	// Act
	res, resErr := NewExpense(id, category, price, currency, quantity, &comment, &trip, date,
		SetCreateMetadata(createdBy, created), SetUpdateMetadata(updatedBy, updated))

	// Assert
	assert.NotNil(t, res)
	assert.Nil(t, resErr)
	assert.Equal(t, id, res.ID())
	assert.Equal(t, category, res.Category())
	assert.Equal(t, price, res.Price())
	assert.Equal(t, currency, res.Currency())
	assert.Equal(t, quantity, res.Quantity())
	assert.Equal(t, &comment, res.Comment())
	assert.Equal(t, &trip, res.Trip())
	assert.Equal(t, date, res.Date())
	assert.Equal(t, createdBy, res.CreatedBy())
	assert.Equal(t, created, res.CreatedAt())
	assert.Equal(t, &updatedBy, res.UpdatedBy())
	assert.Equal(t, &updated, res.UpdatedAt())
}

func TestNewExpense_InvalidArgs_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	type test struct {
		id       string
		price    float64
		quantity float64
		category Category
		date     time.Time
		currency string
		comment  *string
		trip     *string
	}
	tests := []test{
		{
			id:       "id",
			price:    -20,
			quantity: 2,
			currency: "USD",
			category: Category{id: "cat"},
			date:     time.Now(),
			comment:  nil,
			trip:     nil,
		},
		{
			id:       "id",
			price:    0,
			quantity: 2,
			currency: "USD",
			category: Category{id: "cat"},
			date:     time.Now(),
			comment:  nil,
			trip:     nil,
		},
		{
			id:       "id",
			price:    10,
			quantity: -2,
			currency: "USD",
			category: Category{id: "cat"},
			date:     time.Now(),
			comment:  nil,
			trip:     nil,
		},
		{
			id:       "id",
			price:    10,
			quantity: 0,
			currency: "USD",
			category: Category{id: "cat"},
			date:     time.Now(),
			comment:  nil,
			trip:     nil,
		},
		{
			id:       "id",
			price:    10,
			quantity: 2,
			currency: "",
			category: Category{id: "cat"},
			date:     time.Now(),
			comment:  nil,
			trip:     nil,
		},
	}

	for _, tc := range tests {
		// Act
		res, resErr := NewExpense(tc.id, tc.category, tc.price, tc.currency, tc.quantity, tc.comment, tc.trip, tc.date)

		// Assert
		assert.Nil(t, res)
		assert.NotNil(t, resErr)
	}
}

func TestCalculateTotal_NoExchangeRate_ReturnOriginalTotal(t *testing.T) {
	t.Parallel()
	// Arrange
	price := 20.0
	quantity := 2.0
	currency := "EUR"

	// SUT
	sut, _ := NewExpense("id", Category{}, price, currency, quantity, nil, nil, time.Now())

	// Act
	res := sut.CalculateTotal(nil)

	// Assert
	assert.NotNil(t, res)
	assert.Equal(t, sut.totalInfo, res)
	assert.Equal(t, decimal.NewFromFloat(price*quantity), res.OriginalTotal.Sum)
	assert.Equal(t, Currency(currency), res.OriginalTotal.Currency)
	assert.Nil(t, res.ExchangeRate)
	assert.Nil(t, res.ConvertedTotal)
}

func TestCalculateTotal_NoSuitableExchangeRate_ReturnOriginalTotal(t *testing.T) {
	t.Parallel()
	// Arrange
	price := 20.0
	quantity := 2.0
	currency := "EUR"
	rate := 0.2
	exchangeRates := &ExchangeRates{
		baseCurrency: "USD",
		rates: map[Currency]decimal.Decimal{
			"HRK": decimal.NewFromFloat(rate),
		},
	}

	// SUT
	sut, _ := NewExpense("id", Category{}, price, currency, quantity, nil, nil, time.Now())

	// Act
	res := sut.CalculateTotal(exchangeRates)

	// Assert
	assert.NotNil(t, res)
	assert.Equal(t, sut.totalInfo, res)
	assert.Equal(t, decimal.NewFromFloat(price*quantity), res.OriginalTotal.Sum)
	assert.Equal(t, Currency(currency), res.OriginalTotal.Currency)
	assert.Nil(t, res.ExchangeRate)
	assert.Nil(t, res.ConvertedTotal)
}

func TestCalculateTotal_WithExchangeRate_ReturnOriginalAndConvertedTotal(t *testing.T) {
	t.Parallel()
	// Arrange
	price := 20.0
	quantity := 2.1
	currency := "USD"
	rate := 0.23
	exchangeRates := &ExchangeRates{
		baseCurrency: "EUR",
		rates: map[Currency]decimal.Decimal{
			"USD": decimal.NewFromFloat(rate),
		},
	}

	// SUT
	sut, _ := NewExpense("id", Category{}, price, currency, quantity, nil, nil, time.Now())

	// Act
	res := sut.CalculateTotal(exchangeRates)

	// Assert
	assert.NotNil(t, res)
	assert.Equal(t, sut.TotalInfo(), res)
	assert.Equal(t, decimal.NewFromFloat(price*quantity), res.OriginalTotal.Sum)
	assert.Equal(t, Currency(currency), res.OriginalTotal.Currency)
	assert.Equal(t, exchangeRates, res.ExchangeRate)
	assert.Equal(t, decimal.NewFromFloat((price * quantity)).Div(decimal.NewFromFloat(rate)),
		res.ConvertedTotal.Sum)
	assert.Equal(t, exchangeRates.baseCurrency, res.ConvertedTotal.Currency)
}
