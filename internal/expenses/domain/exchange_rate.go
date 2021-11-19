package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type ExchangeRate struct {
	date         time.Time
	baseCurrency Currency
	rates        map[string]float32
}

func NewExchageRate(date time.Time, baseCurrency string, rates map[string]float32) ExchangeRate {
	return ExchangeRate{
		date:         date,
		baseCurrency: Currency(baseCurrency),
		rates:        rates,
	}
}

func (er ExchangeRate) Date() time.Time {
	return er.date
}

func (er ExchangeRate) BaseCurrency() Currency {
	return er.baseCurrency
}

func (er ExchangeRate) Rates() map[Currency]decimal.Decimal {
	rates := make(map[Currency]decimal.Decimal, len(er.rates))
	for currency, rate := range er.rates {
		rates[Currency(currency)] = decimal.NewFromFloat32(rate)
	}
	return rates
}
