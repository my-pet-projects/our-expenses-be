// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	command "dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/app/command"

	domain "dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/domain"

	mock "github.com/stretchr/testify/mock"
)

// FetchExchangeRatesHandlerInterface is an autogenerated mock type for the FetchExchangeRatesHandlerInterface type
type FetchExchangeRatesHandlerInterface struct {
	mock.Mock
}

// Handle provides a mock function with given fields: ctx, cmd
func (_m *FetchExchangeRatesHandlerInterface) Handle(ctx context.Context, cmd command.FetchExchangeRatesCommand) ([]domain.ExchangeRates, error) {
	ret := _m.Called(ctx, cmd)

	var r0 []domain.ExchangeRates
	if rf, ok := ret.Get(0).(func(context.Context, command.FetchExchangeRatesCommand) []domain.ExchangeRates); ok {
		r0 = rf(ctx, cmd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.ExchangeRates)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, command.FetchExchangeRatesCommand) error); ok {
		r1 = rf(ctx, cmd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
