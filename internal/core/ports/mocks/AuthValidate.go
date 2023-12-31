// Code generated by mockery v2.28.2. DO NOT EDIT.

package mocks

import (
	dto "robinhood-assignment/internal/dto"

	gin "github.com/gin-gonic/gin"

	mock "github.com/stretchr/testify/mock"
)

// AuthValidate is an autogenerated mock type for the AuthValidate type
type AuthValidate struct {
	mock.Mock
}

// ValidateCreateStaff provides a mock function with given fields: ctx
func (_m *AuthValidate) ValidateCreateStaff(ctx *gin.Context) (*dto.CreateStaffRequest, error) {
	ret := _m.Called(ctx)

	var r0 *dto.CreateStaffRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(*gin.Context) (*dto.CreateStaffRequest, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(*gin.Context) *dto.CreateStaffRequest); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.CreateStaffRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(*gin.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ValidateLogin provides a mock function with given fields: ctx
func (_m *AuthValidate) ValidateLogin(ctx *gin.Context) (*dto.LoginRequest, error) {
	ret := _m.Called(ctx)

	var r0 *dto.LoginRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(*gin.Context) (*dto.LoginRequest, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(*gin.Context) *dto.LoginRequest); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.LoginRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(*gin.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewAuthValidate interface {
	mock.TestingT
	Cleanup(func())
}

// NewAuthValidate creates a new instance of AuthValidate. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAuthValidate(t mockConstructorTestingTNewAuthValidate) *AuthValidate {
	mock := &AuthValidate{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
