// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	command "dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/app/command"

	mock "github.com/stretchr/testify/mock"
)

// AddExpenseHandlerInterface is an autogenerated mock type for the AddExpenseHandlerInterface type
type AddExpenseHandlerInterface struct {
	mock.Mock
}

// Handle provides a mock function with given fields: ctx, cmd
func (_m *AddExpenseHandlerInterface) Handle(ctx context.Context, cmd command.AddExpenseCommand) (*string, error) {
	ret := _m.Called(ctx, cmd)

	var r0 *string
	if rf, ok := ret.Get(0).(func(context.Context, command.AddExpenseCommand) *string); ok {
		r0 = rf(ctx, cmd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, command.AddExpenseCommand) error); ok {
		r1 = rf(ctx, cmd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}