// Code generated by mockery 2.7.4. DO NOT EDIT.

package mocks

import (
	echo "github.com/labstack/echo/v4"
	mock "github.com/stretchr/testify/mock"

	ports "dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/ports"
)

// ServerInterface is an autogenerated mock type for the ServerInterface type
type ServerInterface struct {
	mock.Mock
}

// FindCategories provides a mock function with given fields: ctx, params
func (_m *ServerInterface) FindCategories(ctx echo.Context, params ports.FindCategoriesParams) error {
	ret := _m.Called(ctx, params)

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context, ports.FindCategoriesParams) error); ok {
		r0 = rf(ctx, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindCategoryByID provides a mock function with given fields: ctx, id
func (_m *ServerInterface) FindCategoryByID(ctx echo.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}