// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// MockRoundTripper is an autogenerated mock type for the RoundTripper type
type MockRoundTripper struct {
	mock.Mock
}

type MockRoundTripper_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRoundTripper) EXPECT() *MockRoundTripper_Expecter {
	return &MockRoundTripper_Expecter{mock: &_m.Mock}
}

// RoundTrip provides a mock function with given fields: _a0
func (_m *MockRoundTripper) RoundTrip(_a0 *http.Request) (*http.Response, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for RoundTrip")
	}

	var r0 *http.Response
	var r1 error
	if rf, ok := ret.Get(0).(func(*http.Request) (*http.Response, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(*http.Request) *http.Response); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Response)
		}
	}

	if rf, ok := ret.Get(1).(func(*http.Request) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRoundTripper_RoundTrip_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RoundTrip'
type MockRoundTripper_RoundTrip_Call struct {
	*mock.Call
}

// RoundTrip is a helper method to define mock.On call
//   - _a0 *http.Request
func (_e *MockRoundTripper_Expecter) RoundTrip(_a0 interface{}) *MockRoundTripper_RoundTrip_Call {
	return &MockRoundTripper_RoundTrip_Call{Call: _e.mock.On("RoundTrip", _a0)}
}

func (_c *MockRoundTripper_RoundTrip_Call) Run(run func(_a0 *http.Request)) *MockRoundTripper_RoundTrip_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*http.Request))
	})
	return _c
}

func (_c *MockRoundTripper_RoundTrip_Call) Return(_a0 *http.Response, _a1 error) *MockRoundTripper_RoundTrip_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRoundTripper_RoundTrip_Call) RunAndReturn(run func(*http.Request) (*http.Response, error)) *MockRoundTripper_RoundTrip_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockRoundTripper creates a new instance of MockRoundTripper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRoundTripper(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRoundTripper {
	mock := &MockRoundTripper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}