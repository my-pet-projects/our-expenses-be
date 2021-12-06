package command

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/adapters"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/tracer"
)

// FetchExchangeRatesCommand defines a fetch command.
type FetchExchangeRatesCommand struct {
	DateRange domain.DateRange
}

// FetchExchangeRatesHandler defines a handler to fetch exchange rates.
type FetchExchangeRatesHandler struct {
	fetcher adapters.ExchangeRateFetcherInterface
	repo    adapters.ExchangeRateRepoInterface
	logger  logger.LogInterface
}

// FetchExchangeRatesHandlerInterface defines a contract to handle command.
type FetchExchangeRatesHandlerInterface interface {
	Handle(ctx context.Context, cmd FetchExchangeRatesCommand) ([]domain.ExchangeRates, error)
}

// NewFetchExchangeRatesHandler returns command handler.
func NewFetchExchangeRatesHandler(
	fetcher adapters.ExchangeRateFetcherInterface,
	repo adapters.ExchangeRateRepoInterface,
	logger logger.LogInterface,
) FetchExchangeRatesHandler {
	return FetchExchangeRatesHandler{
		fetcher: fetcher,
		repo:    repo,
		logger:  logger,
	}
}

// Handle handles fetch exchage rates command.
func (h FetchExchangeRatesHandler) Handle(
	ctx context.Context,
	cmd FetchExchangeRatesCommand,
) ([]domain.ExchangeRates, error) {
	ctx, span := tracer.NewSpan(ctx, "execute fetch exchange rates command")
	defer span.End()

	rates, ratesErr := h.repo.GetAll(ctx, cmd.DateRange)
	if ratesErr != nil {
		tracer.AddSpanError(span, ratesErr)
		return nil, errors.Wrap(ratesErr, "get existing exchange rates")
	}

	dates := cmd.DateRange.DatesInBetween()
	missingRateDates := getMissingRateDates(dates, rates)
	if len(missingRateDates) == 0 {
		return rates, nil
	}

	missingRates, missingRatesErr := h.fetcher.Fetch(ctx, missingRateDates)
	if missingRatesErr != nil {
		tracer.AddSpanError(span, missingRatesErr)
		h.logger.Warnf(ctx, "Failed to fetch exchange rates: %v. Skipping.", missingRatesErr)
		return rates, nil
	}

	_, insErr := h.repo.InsertAll(ctx, missingRates)
	if insErr != nil {
		tracer.AddSpanError(span, insErr)
		return nil, errors.Wrap(insErr, "insert exchange rate into database")
	}

	rates = append(rates, missingRates...)

	return rates, nil
}

func getMissingRateDates(dates []time.Time, rates []domain.ExchangeRates) []time.Time {
	missingDates := make([]time.Time, 0)
	rateDatesSet := make(map[time.Time]struct{})
	for _, rate := range rates {
		rateDatesSet[rate.Date()] = struct{}{}
	}

	for _, date := range dates {
		if date.After(time.Now()) {
			continue
		}
		_, ok := rateDatesSet[date]
		if !ok {
			missingDates = append(missingDates, date)
		}
	}

	return missingDates
}
