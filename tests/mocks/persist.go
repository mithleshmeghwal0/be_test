// Code generated by mockery v2.23.1. DO NOT EDIT.

package mocks

import (
	context "context"

	models "example.com/be_test/internal/models"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// Persist is an autogenerated mock type for the Persist type
type Persist struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: ctx, resource
func (_m *Persist) CreateUser(ctx context.Context, resource *models.User) (*models.User, error) {
	ret := _m.Called(ctx, resource)

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.User) (*models.User, error)); ok {
		return rf(ctx, resource)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.User) *models.User); ok {
		r0 = rf(ctx, resource)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.User) error); ok {
		r1 = rf(ctx, resource)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteUser provides a mock function with given fields: ctx, id
func (_m *Persist) DeleteUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	ret := _m.Called(ctx, id)

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*models.User, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *models.User); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUser provides a mock function with given fields: ctx, id
func (_m *Persist) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	ret := _m.Called(ctx, id)

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*models.User, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *models.User); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListUser provides a mock function with given fields: ctx, filter, pageSize, nextPageToken
func (_m *Persist) ListUser(ctx context.Context, filter string, pageSize int, nextPageToken string) ([]*models.User, string, error) {
	ret := _m.Called(ctx, filter, pageSize, nextPageToken)

	var r0 []*models.User
	var r1 string
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int, string) ([]*models.User, string, error)); ok {
		return rf(ctx, filter, pageSize, nextPageToken)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int, string) []*models.User); ok {
		r0 = rf(ctx, filter, pageSize, nextPageToken)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int, string) string); ok {
		r1 = rf(ctx, filter, pageSize, nextPageToken)
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string, int, string) error); ok {
		r2 = rf(ctx, filter, pageSize, nextPageToken)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// UpdateUser provides a mock function with given fields: ctx, resource, fields
func (_m *Persist) UpdateUser(ctx context.Context, resource *models.User, fields []string) (*models.User, error) {
	ret := _m.Called(ctx, resource, fields)

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.User, []string) (*models.User, error)); ok {
		return rf(ctx, resource, fields)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.User, []string) *models.User); ok {
		r0 = rf(ctx, resource, fields)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.User, []string) error); ok {
		r1 = rf(ctx, resource, fields)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewPersist interface {
	mock.TestingT
	Cleanup(func())
}

// NewPersist creates a new instance of Persist. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPersist(t mockConstructorTestingTNewPersist) *Persist {
	mock := &Persist{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}