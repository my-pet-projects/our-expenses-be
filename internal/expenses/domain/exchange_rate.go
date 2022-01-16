package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// ExchangeRate represts currency exchange rates.
type ExchangeRate struct {
	date           time.Time
	baseCurrency   Currency
	targetCurrency Currency
	rate           decimal.Decimal
}

// NewExchangeRate instantiates exchange rate.
func NewExchangeRate(
	date time.Time,
	baseCurrency Currency,
	targetCurrency Currency,
	rate decimal.Decimal,
) ExchangeRate {
	return ExchangeRate{
		date:           date,
		baseCurrency:   baseCurrency,
		targetCurrency: targetCurrency,
		rate:           rate,
	}
}

// Date returns exchange rates date.
func (er ExchangeRate) Date() time.Time {
	return er.date
}

// BaseCurrency returns exchange rates base currency.
func (er ExchangeRate) BaseCurrency() Currency {
	return er.baseCurrency
}

// TargetCurrency returns exchange rates target currency.
func (er ExchangeRate) TargetCurrency() Currency {
	return er.targetCurrency
}

// Rates returns exchange rates.
func (er ExchangeRate) Rate() decimal.Decimal {
	return er.rate
}
