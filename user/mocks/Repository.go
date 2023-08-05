// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	models "github.com/jordyf15/tweeter-api/models"
	mock "github.com/stretchr/testify/mock"

	user "github.com/jordyf15/tweeter-api/user"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0
func (_m *Repository) Create(_a0 *models.User) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.User) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateTransaction provides a mock function with given fields: fn
func (_m *Repository) CreateTransaction(fn func(user.Repository) error) error {
	ret := _m.Called(fn)

	var r0 error
	if rf, ok := ret.Get(0).(func(func(user.Repository) error) error); ok {
		r0 = rf(fn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRepository(t mockConstructorTestingTNewRepository) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}