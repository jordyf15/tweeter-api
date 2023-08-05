// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	models "github.com/jordyf15/tweeter-api/models"
	mock "github.com/stretchr/testify/mock"

	user "github.com/jordyf15/tweeter-api/user"
)

// Usecase is an autogenerated mock type for the Usecase type
type Usecase struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0
func (_m *Usecase) Create(_a0 *models.User) (map[string]interface{}, error) {
	ret := _m.Called(_a0)

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(*models.User) map[string]interface{}); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.User) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// For provides a mock function with given fields: _a0
func (_m *Usecase) For(_a0 *models.User) user.InstanceUsecase {
	ret := _m.Called(_a0)

	var r0 user.InstanceUsecase
	if rf, ok := ret.Get(0).(func(*models.User) user.InstanceUsecase); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(user.InstanceUsecase)
		}
	}

	return r0
}

// Login provides a mock function with given fields: login, password
func (_m *Usecase) Login(login string, password string) (map[string]interface{}, error) {
	ret := _m.Called(login, password)

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(string, string) map[string]interface{}); ok {
		r0 = rf(login, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(login, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewUsecase interface {
	mock.TestingT
	Cleanup(func())
}

// NewUsecase creates a new instance of Usecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUsecase(t mockConstructorTestingTNewUsecase) *Usecase {
	mock := &Usecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
