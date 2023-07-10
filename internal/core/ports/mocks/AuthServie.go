// Code generated by mockery v2.28.2. DO NOT EDIT.

package mocks

import (
	context "context"
	dto "robinhood-assignment/internal/dto"

	mock "github.com/stretchr/testify/mock"
)

// AuthServie is an autogenerated mock type for the AuthServie type
type AuthServie struct {
	mock.Mock
}

// CreateStaff provides a mock function with given fields: ctx, params
func (_m *AuthServie) CreateStaff(ctx context.Context, params *dto.CreateStaffRequest) error {
	ret := _m.Called(ctx, params)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *dto.CreateStaffRequest) error); ok {
		r0 = rf(ctx, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Login provides a mock function with given fields: ctx, params
func (_m *AuthServie) Login(ctx context.Context, params *dto.LoginRequest) (string, error) {
	ret := _m.Called(ctx, params)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *dto.LoginRequest) (string, error)); ok {
		return rf(ctx, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *dto.LoginRequest) string); ok {
		r0 = rf(ctx, params)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *dto.LoginRequest) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewAuthServie interface {
	mock.TestingT
	Cleanup(func())
}

// NewAuthServie creates a new instance of AuthServie. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAuthServie(t mockConstructorTestingTNewAuthServie) *AuthServie {
	mock := &AuthServie{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}