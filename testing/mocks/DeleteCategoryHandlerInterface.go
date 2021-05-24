// Code generated by mockery 2.7.4. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/domain"
	mock "github.com/stretchr/testify/mock"
)

// DeleteCategoryHandlerInterface is an autogenerated mock type for the DeleteCategoryHandlerInterface type
type DeleteCategoryHandlerInterface struct {
	mock.Mock
}

// Handle provides a mock function with given fields: ctx, id
func (_m *DeleteCategoryHandlerInterface) Handle(ctx context.Context, id string) (*domain.DeleteResult, error) {
	ret := _m.Called(ctx, id)

	var r0 *domain.DeleteResult
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.DeleteResult); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.DeleteResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
