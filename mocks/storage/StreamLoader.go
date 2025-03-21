// Code generated by mockery v2.52.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// StreamLoader is an autogenerated mock type for the StreamLoader type
type StreamLoader struct {
	mock.Mock
}

type StreamLoader_Expecter struct {
	mock *mock.Mock
}

func (_m *StreamLoader) EXPECT() *StreamLoader_Expecter {
	return &StreamLoader_Expecter{mock: &_m.Mock}
}

// SetModel provides a mock function with given fields: model
func (_m *StreamLoader) SetModel(model interface{}) {
	_m.Called(model)
}

// StreamLoader_SetModel_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetModel'
type StreamLoader_SetModel_Call struct {
	*mock.Call
}

// SetModel is a helper method to define mock.On call
//   - model interface{}
func (_e *StreamLoader_Expecter) SetModel(model interface{}) *StreamLoader_SetModel_Call {
	return &StreamLoader_SetModel_Call{Call: _e.mock.On("SetModel", model)}
}

func (_c *StreamLoader_SetModel_Call) Run(run func(model interface{})) *StreamLoader_SetModel_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *StreamLoader_SetModel_Call) Return() *StreamLoader_SetModel_Call {
	_c.Call.Return()
	return _c
}

func (_c *StreamLoader_SetModel_Call) RunAndReturn(run func(interface{})) *StreamLoader_SetModel_Call {
	_c.Run(run)
	return _c
}

// Start provides a mock function with no fields
func (_m *StreamLoader) Start() (<-chan interface{}, <-chan error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Start")
	}

	var r0 <-chan interface{}
	var r1 <-chan error
	if rf, ok := ret.Get(0).(func() (<-chan interface{}, <-chan error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() <-chan interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan interface{})
		}
	}

	if rf, ok := ret.Get(1).(func() <-chan error); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(<-chan error)
		}
	}

	return r0, r1
}

// StreamLoader_Start_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Start'
type StreamLoader_Start_Call struct {
	*mock.Call
}

// Start is a helper method to define mock.On call
func (_e *StreamLoader_Expecter) Start() *StreamLoader_Start_Call {
	return &StreamLoader_Start_Call{Call: _e.mock.On("Start")}
}

func (_c *StreamLoader_Start_Call) Run(run func()) *StreamLoader_Start_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *StreamLoader_Start_Call) Return(_a0 <-chan interface{}, _a1 <-chan error) *StreamLoader_Start_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *StreamLoader_Start_Call) RunAndReturn(run func() (<-chan interface{}, <-chan error)) *StreamLoader_Start_Call {
	_c.Call.Return(run)
	return _c
}

// NewStreamLoader creates a new instance of StreamLoader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStreamLoader(t interface {
	mock.TestingT
	Cleanup(func())
}) *StreamLoader {
	mock := &StreamLoader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
