package domain

import (
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type ExchangeRates struct {
	date         time.Time
	baseCurrency Currency
	rates        map[Currency]decimal.Decimal
}

func NewExchageRate(date time.Time, baseCurrency string, rawRates map[string]float64) (*ExchangeRates, error) {
	if baseCurrency == "" {
		return nil, errors.New("base currency should not be empty")
	}
	if len(rawRates) == 0 {
		return nil, errors.New("rates should not be empty")
	}

	rates := make(map[Currency]decimal.Decimal, len(rawRates))
	for currency, rate := range rawRates {
		rates[Currency(currency)] = decimal.NewFromFloat(rate)
	}
	er := ExchangeRates{
		date:         date,
		baseCurrency: Currency(baseCurrency),
		rates:        rates,
	}

	return &er, nil
}

func (er ExchangeRates) Date() time.Time {
	return er.date
}

func (er ExchangeRates) BaseCurrency() Currency {
	return er.baseCurrency
}

func (er ExchangeRates) Rates() map[Currency]decimal.Decimal {
	return er.rates
}

func (er ExchangeRates) ChangeBaseCurrency(targetCurrency Currency) ExchangeRates {
	if er.baseCurrency == targetCurrency {
		return er
	}

	baseRate := er.rates[targetCurrency]
	newRate := ExchangeRates{
		date:         er.date,
		baseCurrency: targetCurrency,
		rates:        make(map[Currency]decimal.Decimal, len(er.rates)),
	}
	for currency, rate := range er.rates {
		if currency == targetCurrency {
			continue
		}
		rate := rate.Div(baseRate)
		newRate.rates[currency] = rate
	}

	if currencyRate, ok := er.rates[targetCurrency]; ok {
		newRate.rates[er.baseCurrency] = decimal.NewFromInt(1).Div(currencyRate)
	}

	return newRate
}
