// Code generated by mockery v2.52.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Runner is an autogenerated mock type for the Runner type
type Runner struct {
	mock.Mock
}

type Runner_Expecter struct {
	mock *mock.Mock
}

func (_m *Runner) EXPECT() *Runner_Expecter {
	return &Runner_Expecter{mock: &_m.Mock}
}

// Run provides a mock function with given fields: ctx
func (_m *Runner) Run(ctx context.Context) (int, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Run")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (int, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) int); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Runner_Run_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Run'
type Runner_Run_Call struct {
	*mock.Call
}

// Run is a helper method to define mock.On call
//   - ctx context.Context
func (_e *Runner_Expecter) Run(ctx interface{}) *Runner_Run_Call {
	return &Runner_Run_Call{Call: _e.mock.On("Run", ctx)}
}

func (_c *Runner_Run_Call) Run(run func(ctx context.Context)) *Runner_Run_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *Runner_Run_Call) Return(finalTaskUID int, err error) *Runner_Run_Call {
	_c.Call.Return(finalTaskUID, err)
	return _c
}

func (_c *Runner_Run_Call) RunAndReturn(run func(context.Context) (int, error)) *Runner_Run_Call {
	_c.Call.Return(run)
	return _c
}

// NewRunner creates a new instance of Runner. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRunner(t interface {
	mock.TestingT
	Cleanup(func())
}) *Runner {
	mock := &Runner{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
