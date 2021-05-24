// Code generated by mockery 2.7.4. DO NOT EDIT.

package mocks

import (
	context "context"

	command "dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/app/command"

	domain "dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/domain"

	mock "github.com/stretchr/testify/mock"
)

// UpdateCategoryHandlerInterface is an autogenerated mock type for the UpdateCategoryHandlerInterface type
type UpdateCategoryHandlerInterface struct {
	mock.Mock
}

// Handle provides a mock function with given fields: ctx, cmd
func (_m *UpdateCategoryHandlerInterface) Handle(ctx context.Context, cmd command.UpdateCategoryCommand) (*domain.UpdateResult, error) {
	ret := _m.Called(ctx, cmd)

	var r0 *domain.UpdateResult
	if rf, ok := ret.Get(0).(func(context.Context, command.UpdateCategoryCommand) *domain.UpdateResult); ok {
		r0 = rf(ctx, cmd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.UpdateResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, command.UpdateCategoryCommand) error); ok {
		r1 = rf(ctx, cmd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
