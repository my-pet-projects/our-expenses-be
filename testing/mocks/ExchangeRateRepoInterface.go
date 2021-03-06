// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/domain"
	mock "github.com/stretchr/testify/mock"
)

// ExchangeRateRepoInterface is an autogenerated mock type for the ExchangeRateRepoInterface type
type ExchangeRateRepoInterface struct {
	mock.Mock
}

// GetAll provides a mock function with given fields: ctx, dateRange
func (_m *ExchangeRateRepoInterface) GetAll(ctx context.Context, dateRange domain.DateRange) ([]domain.ExchangeRates, error) {
	ret := _m.Called(ctx, dateRange)

	var r0 []domain.ExchangeRates
	if rf, ok := ret.Get(0).(func(context.Context, domain.DateRange) []domain.ExchangeRates); ok {
		r0 = rf(ctx, dateRange)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.ExchangeRates)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.DateRange) error); ok {
		r1 = rf(ctx, dateRange)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertAll provides a mock function with given fields: ctx, rates
func (_m *ExchangeRateRepoInterface) InsertAll(ctx context.Context, rates []domain.ExchangeRates) (*domain.InsertResult, error) {
	ret := _m.Called(ctx, rates)

	var r0 *domain.InsertResult
	if rf, ok := ret.Get(0).(func(context.Context, []domain.ExchangeRates) *domain.InsertResult); ok {
		r0 = rf(ctx, rates)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.InsertResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []domain.ExchangeRates) error); ok {
		r1 = rf(ctx, rates)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
