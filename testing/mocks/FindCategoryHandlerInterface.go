// Code generated by mockery 2.7.4. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/domain"
	mock "github.com/stretchr/testify/mock"

	query "dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/app/query"
)

// FindCategoryHandlerInterface is an autogenerated mock type for the FindCategoryHandlerInterface type
type FindCategoryHandlerInterface struct {
	mock.Mock
}

// Handle provides a mock function with given fields: ctx, _a1
func (_m *FindCategoryHandlerInterface) Handle(ctx context.Context, _a1 query.FindCategoryQuery) (*domain.Category, error) {
	ret := _m.Called(ctx, _a1)

	var r0 *domain.Category
	if rf, ok := ret.Get(0).(func(context.Context, query.FindCategoryQuery) *domain.Category); ok {
		r0 = rf(ctx, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Category)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, query.FindCategoryQuery) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
