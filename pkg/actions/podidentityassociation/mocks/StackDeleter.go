// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	types "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"

	v1alpha5 "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// StackDeleter is an autogenerated mock type for the StackDeleter type
type StackDeleter struct {
	mock.Mock
}

type StackDeleter_Expecter struct {
	mock *mock.Mock
}

func (_m *StackDeleter) EXPECT() *StackDeleter_Expecter {
	return &StackDeleter_Expecter{mock: &_m.Mock}
}

// DeleteStackBySpecSync provides a mock function with given fields: ctx, stack, errCh
func (_m *StackDeleter) DeleteStackBySpecSync(ctx context.Context, stack *types.Stack, errCh chan error) error {
	ret := _m.Called(ctx, stack, errCh)

	if len(ret) == 0 {
		panic("no return value specified for DeleteStackBySpecSync")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.Stack, chan error) error); ok {
		r0 = rf(ctx, stack, errCh)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StackDeleter_DeleteStackBySpecSync_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteStackBySpecSync'
type StackDeleter_DeleteStackBySpecSync_Call struct {
	*mock.Call
}

// DeleteStackBySpecSync is a helper method to define mock.On call
//   - ctx context.Context
//   - stack *types.Stack
//   - errCh chan error
func (_e *StackDeleter_Expecter) DeleteStackBySpecSync(ctx interface{}, stack interface{}, errCh interface{}) *StackDeleter_DeleteStackBySpecSync_Call {
	return &StackDeleter_DeleteStackBySpecSync_Call{Call: _e.mock.On("DeleteStackBySpecSync", ctx, stack, errCh)}
}

func (_c *StackDeleter_DeleteStackBySpecSync_Call) Run(run func(ctx context.Context, stack *types.Stack, errCh chan error)) *StackDeleter_DeleteStackBySpecSync_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*types.Stack), args[2].(chan error))
	})
	return _c
}

func (_c *StackDeleter_DeleteStackBySpecSync_Call) Return(_a0 error) *StackDeleter_DeleteStackBySpecSync_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *StackDeleter_DeleteStackBySpecSync_Call) RunAndReturn(run func(context.Context, *types.Stack, chan error) error) *StackDeleter_DeleteStackBySpecSync_Call {
	_c.Call.Return(run)
	return _c
}

// DescribeStack provides a mock function with given fields: ctx, stack
func (_m *StackDeleter) DescribeStack(ctx context.Context, stack *types.Stack) (*types.Stack, error) {
	ret := _m.Called(ctx, stack)

	if len(ret) == 0 {
		panic("no return value specified for DescribeStack")
	}

	var r0 *types.Stack
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.Stack) (*types.Stack, error)); ok {
		return rf(ctx, stack)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.Stack) *types.Stack); ok {
		r0 = rf(ctx, stack)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Stack)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.Stack) error); ok {
		r1 = rf(ctx, stack)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StackDeleter_DescribeStack_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DescribeStack'
type StackDeleter_DescribeStack_Call struct {
	*mock.Call
}

// DescribeStack is a helper method to define mock.On call
//   - ctx context.Context
//   - stack *types.Stack
func (_e *StackDeleter_Expecter) DescribeStack(ctx interface{}, stack interface{}) *StackDeleter_DescribeStack_Call {
	return &StackDeleter_DescribeStack_Call{Call: _e.mock.On("DescribeStack", ctx, stack)}
}

func (_c *StackDeleter_DescribeStack_Call) Run(run func(ctx context.Context, stack *types.Stack)) *StackDeleter_DescribeStack_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*types.Stack))
	})
	return _c
}

func (_c *StackDeleter_DescribeStack_Call) Return(_a0 *types.Stack, _a1 error) *StackDeleter_DescribeStack_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *StackDeleter_DescribeStack_Call) RunAndReturn(run func(context.Context, *types.Stack) (*types.Stack, error)) *StackDeleter_DescribeStack_Call {
	_c.Call.Return(run)
	return _c
}

// GetIAMServiceAccounts provides a mock function with given fields: ctx
func (_m *StackDeleter) GetIAMServiceAccounts(ctx context.Context) ([]*v1alpha5.ClusterIAMServiceAccount, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetIAMServiceAccounts")
	}

	var r0 []*v1alpha5.ClusterIAMServiceAccount
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*v1alpha5.ClusterIAMServiceAccount, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*v1alpha5.ClusterIAMServiceAccount); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*v1alpha5.ClusterIAMServiceAccount)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StackDeleter_GetIAMServiceAccounts_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetIAMServiceAccounts'
type StackDeleter_GetIAMServiceAccounts_Call struct {
	*mock.Call
}

// GetIAMServiceAccounts is a helper method to define mock.On call
//   - ctx context.Context
func (_e *StackDeleter_Expecter) GetIAMServiceAccounts(ctx interface{}) *StackDeleter_GetIAMServiceAccounts_Call {
	return &StackDeleter_GetIAMServiceAccounts_Call{Call: _e.mock.On("GetIAMServiceAccounts", ctx)}
}

func (_c *StackDeleter_GetIAMServiceAccounts_Call) Run(run func(ctx context.Context)) *StackDeleter_GetIAMServiceAccounts_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *StackDeleter_GetIAMServiceAccounts_Call) Return(_a0 []*v1alpha5.ClusterIAMServiceAccount, _a1 error) *StackDeleter_GetIAMServiceAccounts_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *StackDeleter_GetIAMServiceAccounts_Call) RunAndReturn(run func(context.Context) ([]*v1alpha5.ClusterIAMServiceAccount, error)) *StackDeleter_GetIAMServiceAccounts_Call {
	_c.Call.Return(run)
	return _c
}

// GetStackTemplate provides a mock function with given fields: ctx, stackName
func (_m *StackDeleter) GetStackTemplate(ctx context.Context, stackName string) (string, error) {
	ret := _m.Called(ctx, stackName)

	if len(ret) == 0 {
		panic("no return value specified for GetStackTemplate")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, stackName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, stackName)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, stackName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StackDeleter_GetStackTemplate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetStackTemplate'
type StackDeleter_GetStackTemplate_Call struct {
	*mock.Call
}

// GetStackTemplate is a helper method to define mock.On call
//   - ctx context.Context
//   - stackName string
func (_e *StackDeleter_Expecter) GetStackTemplate(ctx interface{}, stackName interface{}) *StackDeleter_GetStackTemplate_Call {
	return &StackDeleter_GetStackTemplate_Call{Call: _e.mock.On("GetStackTemplate", ctx, stackName)}
}

func (_c *StackDeleter_GetStackTemplate_Call) Run(run func(ctx context.Context, stackName string)) *StackDeleter_GetStackTemplate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *StackDeleter_GetStackTemplate_Call) Return(_a0 string, _a1 error) *StackDeleter_GetStackTemplate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *StackDeleter_GetStackTemplate_Call) RunAndReturn(run func(context.Context, string) (string, error)) *StackDeleter_GetStackTemplate_Call {
	_c.Call.Return(run)
	return _c
}

// ListPodIdentityStackNames provides a mock function with given fields: ctx
func (_m *StackDeleter) ListPodIdentityStackNames(ctx context.Context) ([]string, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for ListPodIdentityStackNames")
	}

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]string, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []string); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StackDeleter_ListPodIdentityStackNames_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListPodIdentityStackNames'
type StackDeleter_ListPodIdentityStackNames_Call struct {
	*mock.Call
}

// ListPodIdentityStackNames is a helper method to define mock.On call
//   - ctx context.Context
func (_e *StackDeleter_Expecter) ListPodIdentityStackNames(ctx interface{}) *StackDeleter_ListPodIdentityStackNames_Call {
	return &StackDeleter_ListPodIdentityStackNames_Call{Call: _e.mock.On("ListPodIdentityStackNames", ctx)}
}

func (_c *StackDeleter_ListPodIdentityStackNames_Call) Run(run func(ctx context.Context)) *StackDeleter_ListPodIdentityStackNames_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *StackDeleter_ListPodIdentityStackNames_Call) Return(_a0 []string, _a1 error) *StackDeleter_ListPodIdentityStackNames_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *StackDeleter_ListPodIdentityStackNames_Call) RunAndReturn(run func(context.Context) ([]string, error)) *StackDeleter_ListPodIdentityStackNames_Call {
	_c.Call.Return(run)
	return _c
}

// NewStackDeleter creates a new instance of StackDeleter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStackDeleter(t interface {
	mock.TestingT
	Cleanup(func())
}) *StackDeleter {
	mock := &StackDeleter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}