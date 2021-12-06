package command_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/app/command"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/testing/mocks"
)

func TestFetchExchangeRatesHandler_ReturnsHandler(t *testing.T) {
	t.Parallel()
	// Arrange
	fetcher := new(mocks.ExchangeRateFetcherInterface)
	repo := new(mocks.ExchangeRateRepoInterface)
	log := new(mocks.LogInterface)

	// Act
	err := command.NewFetchExchangeRatesHandler(fetcher, repo, log)

	// Assert
	assert.NotNil(t, err)
}

func TestFetchExchangeRatesHandler_RepoFails_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	fetcher := new(mocks.ExchangeRateFetcherInterface)
	repo := new(mocks.ExchangeRateRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	timeTo := time.Now()
	timeFrom := timeTo.Add(-1 * 24 * time.Hour)
	dateRange, _ := domain.NewDateRange(timeFrom, timeTo)
	cmd := command.FetchExchangeRatesCommand{
		DateRange: *dateRange,
	}

	matchDateRangeFn := func(dr domain.DateRange) bool {
		return dr == *dateRange
	}
	repo.On("GetAll", mock.Anything,
		mock.MatchedBy(matchDateRangeFn)).Return(nil, errors.New("error"))

	// SUT
	sut := command.NewFetchExchangeRatesHandler(fetcher, repo, log)

	// Act
	res, resErr := sut.Handle(ctx, cmd)

	// Assert
	assert.Nil(t, res)
	assert.NotNil(t, resErr)
}

func TestFetchExchangeRatesHandler_NoMissingRates_ReturnsRates(t *testing.T) {
	t.Parallel()
	// Arrange
	fetcher := new(mocks.ExchangeRateFetcherInterface)
	repo := new(mocks.ExchangeRateRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	time1 := time.Now()
	time2 := time1.Add(-1 * 24 * time.Hour)
	time3 := time1.Add(-2 * 24 * time.Hour)
	dateRange, _ := domain.NewDateRange(time2, time1)
	cmd := command.FetchExchangeRatesCommand{
		DateRange: *dateRange,
	}
	rate1, _ := domain.NewExchageRate(time1, "EUR", map[string]float64{"USD": 1})
	rate2, _ := domain.NewExchageRate(time2, "EUR", map[string]float64{"USD": 1})
	rate3, _ := domain.NewExchageRate(time3, "EUR", map[string]float64{"USD": 1})
	rates := []domain.ExchangeRates{*rate1, *rate2, *rate3}

	matchDateRangeFn := func(dr domain.DateRange) bool {
		return dr == *dateRange
	}
	repo.On("GetAll", mock.Anything,
		mock.MatchedBy(matchDateRangeFn)).Return(rates, nil)

	// SUT
	sut := command.NewFetchExchangeRatesHandler(fetcher, repo, log)

	// Act
	res, resErr := sut.Handle(ctx, cmd)

	// Assert
	assert.NotNil(t, res)
	assert.Equal(t, rates, res)
	assert.Nil(t, resErr)
}

func TestFetchExchangeRatesHandler_MissingRates_FetcherFails_ReturnRatesThatHave(t *testing.T) {
	t.Parallel()
	// Arrange
	fetcher := new(mocks.ExchangeRateFetcherInterface)
	repo := new(mocks.ExchangeRateRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	time1 := time.Now()
	time2 := time1.Add(-1 * 24 * time.Hour)
	time3 := time1.Add(-2 * 24 * time.Hour)
	time4 := time1.Add(-3 * 24 * time.Hour)
	time5 := time1.Add(-4 * 24 * time.Hour)
	dateRange, _ := domain.NewDateRange(time5, time1)
	cmd := command.FetchExchangeRatesCommand{
		DateRange: *dateRange,
	}
	rate1, _ := domain.NewExchageRate(time1, "EUR", map[string]float64{"USD": 1})
	rate2, _ := domain.NewExchageRate(time2, "EUR", map[string]float64{"USD": 1})
	rate3, _ := domain.NewExchageRate(time3, "EUR", map[string]float64{"USD": 1})
	rates := []domain.ExchangeRates{*rate1, *rate2, *rate3}

	matchDateRangeFn := func(dr domain.DateRange) bool {
		return dr == *dateRange
	}
	repo.On("GetAll", mock.Anything,
		mock.MatchedBy(matchDateRangeFn)).Return(rates, nil)

	matchMissingDatesFn := func(d []time.Time) bool {
		return reflect.DeepEqual(d, []time.Time{time5, time4})
	}
	fetcher.On("Fetch", mock.Anything,
		mock.MatchedBy(matchMissingDatesFn)).Return(nil, errors.New("error"))

	log.On("Warnf", mock.Anything, mock.Anything, mock.Anything).Return()

	// SUT
	sut := command.NewFetchExchangeRatesHandler(fetcher, repo, log)

	// Act
	res, resErr := sut.Handle(ctx, cmd)

	// Assert
	assert.NotNil(t, res)
	assert.Equal(t, rates, res)
	assert.Nil(t, resErr)
}

func TestFetchExchangeRatesHandler_MissingRates_RepoFails_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	fetcher := new(mocks.ExchangeRateFetcherInterface)
	repo := new(mocks.ExchangeRateRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	time1 := time.Now()
	time2 := time1.Add(-1 * 24 * time.Hour)
	time3 := time1.Add(-2 * 24 * time.Hour)
	time4 := time1.Add(-3 * 24 * time.Hour)
	time5 := time1.Add(-4 * 24 * time.Hour)
	dateRange, _ := domain.NewDateRange(time5, time1)
	cmd := command.FetchExchangeRatesCommand{
		DateRange: *dateRange,
	}
	rate1, _ := domain.NewExchageRate(time1, "EUR", map[string]float64{"USD": 1})
	rate2, _ := domain.NewExchageRate(time2, "EUR", map[string]float64{"USD": 1})
	rate3, _ := domain.NewExchageRate(time3, "EUR", map[string]float64{"USD": 1})
	rate4, _ := domain.NewExchageRate(time4, "EUR", map[string]float64{"USD": 1})
	rate5, _ := domain.NewExchageRate(time5, "EUR", map[string]float64{"USD": 1})
	rates := []domain.ExchangeRates{*rate1, *rate2, *rate3}
	missingRates := []domain.ExchangeRates{*rate4, *rate5}

	matchDateRangeFn := func(dr domain.DateRange) bool {
		return dr == *dateRange
	}
	repo.On("GetAll", mock.Anything,
		mock.MatchedBy(matchDateRangeFn)).Return(rates, nil)

	matchMissingDatesFn := func(d []time.Time) bool {
		return reflect.DeepEqual(d, []time.Time{time5, time4})
	}
	fetcher.On("Fetch", mock.Anything,
		mock.MatchedBy(matchMissingDatesFn)).Return(missingRates, nil)

	matchMissingRatesFn := func(r []domain.ExchangeRates) bool {
		return reflect.DeepEqual(r, missingRates)
	}
	repo.On("InsertAll", mock.Anything,
		mock.MatchedBy(matchMissingRatesFn)).Return(nil, errors.New("error"))

	// SUT
	sut := command.NewFetchExchangeRatesHandler(fetcher, repo, log)

	// Act
	res, resErr := sut.Handle(ctx, cmd)

	// Assert
	assert.Nil(t, res)
	assert.NotNil(t, resErr)
}

func TestFetchExchangeRatesHandler_MissingRates_ReturnCombinedRates(t *testing.T) {
	t.Parallel()
	// Arrange
	fetcher := new(mocks.ExchangeRateFetcherInterface)
	repo := new(mocks.ExchangeRateRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	time1 := time.Now().Add(1 * 24 * time.Hour)
	time2 := time1.Add(-1 * 24 * time.Hour)
	time3 := time1.Add(-2 * 24 * time.Hour)
	time4 := time1.Add(-3 * 24 * time.Hour)
	time5 := time1.Add(-4 * 24 * time.Hour)
	dateRange, _ := domain.NewDateRange(time5, time1)
	cmd := command.FetchExchangeRatesCommand{
		DateRange: *dateRange,
	}
	rate1, _ := domain.NewExchageRate(time1, "EUR", map[string]float64{"USD": 1})
	rate2, _ := domain.NewExchageRate(time2, "EUR", map[string]float64{"USD": 1})
	rate3, _ := domain.NewExchageRate(time3, "EUR", map[string]float64{"USD": 1})
	rate4, _ := domain.NewExchageRate(time4, "EUR", map[string]float64{"USD": 1})
	rate5, _ := domain.NewExchageRate(time5, "EUR", map[string]float64{"USD": 1})
	rates := []domain.ExchangeRates{*rate1, *rate2, *rate3}
	missingRates := []domain.ExchangeRates{*rate4, *rate5}

	matchDateRangeFn := func(dr domain.DateRange) bool {
		return dr == *dateRange
	}
	repo.On("GetAll", mock.Anything,
		mock.MatchedBy(matchDateRangeFn)).Return(rates, nil)

	matchMissingDatesFn := func(d []time.Time) bool {
		return reflect.DeepEqual(d, []time.Time{time5, time4})
	}
	fetcher.On("Fetch", mock.Anything,
		mock.MatchedBy(matchMissingDatesFn)).Return(missingRates, nil)

	matchMissingRatesFn := func(r []domain.ExchangeRates) bool {
		return reflect.DeepEqual(r, missingRates)
	}
	repo.On("InsertAll", mock.Anything,
		mock.MatchedBy(matchMissingRatesFn)).Return(&domain.InsertResult{}, nil)

	// SUT
	sut := command.NewFetchExchangeRatesHandler(fetcher, repo, log)

	// Act
	res, resErr := sut.Handle(ctx, cmd)

	// Assert
	assert.NotNil(t, res)
	assert.Equal(t, append(rates, missingRates...), res)
	assert.Nil(t, resErr)
}
