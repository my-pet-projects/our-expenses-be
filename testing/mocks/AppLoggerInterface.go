// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// AppLoggerInterface is an autogenerated mock type for the AppLoggerInterface type
type AppLoggerInterface struct {
	mock.Mock
}

// Error provides a mock function with given fields: errMessage, err, fields
func (_m *AppLoggerInterface) Error(errMessage string, err error, fields map[string]interface{}) {
	_m.Called(errMessage, err, fields)
}

// Fatal provides a mock function with given fields: errMessage, err, fields
func (_m *AppLoggerInterface) Fatal(errMessage string, err error, fields map[string]interface{}) {
	_m.Called(errMessage, err, fields)
}

// Info provides a mock function with given fields: msg, fields
func (_m *AppLoggerInterface) Info(msg string, fields map[string]interface{}) {
	_m.Called(msg, fields)
}