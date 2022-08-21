// Code generated by mockery v2.14.0. DO NOT EDIT.

package server

import (
	http "net/http"

	storage "github.com/iamwavecut/ct-mend/internal/storage"
	mock "github.com/stretchr/testify/mock"
)

// MockRESTHandler is an autogenerated mock type for the RESTHandler type
type MockRESTHandler struct {
	mock.Mock
}

// Delete provides a mock function with given fields: _a0, _a1
func (_m *MockRESTHandler) Delete(_a0 http.ResponseWriter, _a1 *http.Request) {
	_m.Called(_a0, _a1)
}

// Get provides a mock function with given fields: _a0, _a1
func (_m *MockRESTHandler) Get(_a0 http.ResponseWriter, _a1 *http.Request) {
	_m.Called(_a0, _a1)
}

// Post provides a mock function with given fields: _a0, _a1
func (_m *MockRESTHandler) Post(_a0 http.ResponseWriter, _a1 *http.Request) {
	_m.Called(_a0, _a1)
}

// Put provides a mock function with given fields: _a0, _a1
func (_m *MockRESTHandler) Put(_a0 http.ResponseWriter, _a1 *http.Request) {
	_m.Called(_a0, _a1)
}

// Select provides a mock function with given fields: _a0, _a1
func (_m *MockRESTHandler) Select(_a0 http.ResponseWriter, _a1 *http.Request) {
	_m.Called(_a0, _a1)
}

// WithStorageAdapter provides a mock function with given fields: db
func (_m *MockRESTHandler) WithStorageAdapter(db storage.Adapter) RESTHandler {
	ret := _m.Called(db)

	var r0 RESTHandler
	if rf, ok := ret.Get(0).(func(storage.Adapter) RESTHandler); ok {
		r0 = rf(db)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(RESTHandler)
		}
	}

	return r0
}

type mockConstructorTestingTNewMockRESTHandler interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockRESTHandler creates a new instance of MockRESTHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockRESTHandler(t mockConstructorTestingTNewMockRESTHandler) *MockRESTHandler {
	mock := &MockRESTHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}