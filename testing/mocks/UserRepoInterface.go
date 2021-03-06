// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/users/domain"
	mock "github.com/stretchr/testify/mock"
)

// UserRepoInterface is an autogenerated mock type for the UserRepoInterface type
type UserRepoInterface struct {
	mock.Mock
}

// GetOne provides a mock function with given fields: ctx, username
func (_m *UserRepoInterface) GetOne(ctx context.Context, username string) (*domain.User, error) {
	ret := _m.Called(ctx, username)

	var r0 *domain.User
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.User); ok {
		r0 = rf(ctx, username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Insert provides a mock function with given fields: ctx, user
func (_m *UserRepoInterface) Insert(ctx context.Context, user *domain.User) (*domain.InsertResult, error) {
	ret := _m.Called(ctx, user)

	var r0 *domain.InsertResult
	if rf, ok := ret.Get(0).(func(context.Context, *domain.User) *domain.InsertResult); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.InsertResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *domain.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, user
func (_m *UserRepoInterface) Update(ctx context.Context, user *domain.User) (*domain.UpdateResult, error) {
	ret := _m.Called(ctx, user)

	var r0 *domain.UpdateResult
	if rf, ok := ret.Get(0).(func(context.Context, *domain.User) *domain.UpdateResult); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.UpdateResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *domain.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
