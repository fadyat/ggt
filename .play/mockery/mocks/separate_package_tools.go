// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// SeparatePackageTools is an autogenerated mock type for the SeparatePackageTools type
type SeparatePackageTools struct {
	mock.Mock
}

type SeparatePackageTools_Expecter struct {
	mock *mock.Mock
}

func (_m *SeparatePackageTools) EXPECT() *SeparatePackageTools_Expecter {
	return &SeparatePackageTools_Expecter{mock: &_m.Mock}
}

// String provides a mock function with given fields:
func (_m *SeparatePackageTools) String() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for String")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// SeparatePackageTools_String_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'String'
type SeparatePackageTools_String_Call struct {
	*mock.Call
}

// String is a helper method to define mock.On call
func (_e *SeparatePackageTools_Expecter) String() *SeparatePackageTools_String_Call {
	return &SeparatePackageTools_String_Call{Call: _e.mock.On("String")}
}

func (_c *SeparatePackageTools_String_Call) Run(run func()) *SeparatePackageTools_String_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *SeparatePackageTools_String_Call) Return(_a0 string) *SeparatePackageTools_String_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SeparatePackageTools_String_Call) RunAndReturn(run func() string) *SeparatePackageTools_String_Call {
	_c.Call.Return(run)
	return _c
}

// NewSeparatePackageTools creates a new instance of SeparatePackageTools. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSeparatePackageTools(t interface {
	mock.TestingT
	Cleanup(func())
}) *SeparatePackageTools {
	mock := &SeparatePackageTools{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
