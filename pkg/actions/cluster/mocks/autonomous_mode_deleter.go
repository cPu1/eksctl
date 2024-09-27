// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// AutonomousModeDeleter is an autogenerated mock type for the AutonomousModeDeleter type
type AutonomousModeDeleter struct {
	mock.Mock
}

type AutonomousModeDeleter_Expecter struct {
	mock *mock.Mock
}

func (_m *AutonomousModeDeleter) EXPECT() *AutonomousModeDeleter_Expecter {
	return &AutonomousModeDeleter_Expecter{mock: &_m.Mock}
}

// DeleteIfRequired provides a mock function with given fields: ctx
func (_m *AutonomousModeDeleter) DeleteIfRequired(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for DeleteIfRequired")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AutonomousModeDeleter_DeleteIfRequired_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteIfRequired'
type AutonomousModeDeleter_DeleteIfRequired_Call struct {
	*mock.Call
}

// DeleteIfRequired is a helper method to define mock.On call
//   - ctx context.Context
func (_e *AutonomousModeDeleter_Expecter) DeleteIfRequired(ctx interface{}) *AutonomousModeDeleter_DeleteIfRequired_Call {
	return &AutonomousModeDeleter_DeleteIfRequired_Call{Call: _e.mock.On("DeleteIfRequired", ctx)}
}

func (_c *AutonomousModeDeleter_DeleteIfRequired_Call) Run(run func(ctx context.Context)) *AutonomousModeDeleter_DeleteIfRequired_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *AutonomousModeDeleter_DeleteIfRequired_Call) Return(_a0 error) *AutonomousModeDeleter_DeleteIfRequired_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AutonomousModeDeleter_DeleteIfRequired_Call) RunAndReturn(run func(context.Context) error) *AutonomousModeDeleter_DeleteIfRequired_Call {
	_c.Call.Return(run)
	return _c
}

// NewAutonomousModeDeleter creates a new instance of AutonomousModeDeleter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAutonomousModeDeleter(t interface {
	mock.TestingT
	Cleanup(func())
}) *AutonomousModeDeleter {
	mock := &AutonomousModeDeleter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
