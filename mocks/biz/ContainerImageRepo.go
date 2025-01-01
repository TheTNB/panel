// Code generated by mockery. DO NOT EDIT.

package biz

import (
	mock "github.com/stretchr/testify/mock"
	request "github.com/tnb-labs/panel/internal/http/request"
	types "github.com/tnb-labs/panel/pkg/types"
)

// ContainerImageRepo is an autogenerated mock type for the ContainerImageRepo type
type ContainerImageRepo struct {
	mock.Mock
}

type ContainerImageRepo_Expecter struct {
	mock *mock.Mock
}

func (_m *ContainerImageRepo) EXPECT() *ContainerImageRepo_Expecter {
	return &ContainerImageRepo_Expecter{mock: &_m.Mock}
}

// List provides a mock function with no fields
func (_m *ContainerImageRepo) List() ([]types.ContainerImage, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 []types.ContainerImage
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]types.ContainerImage, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []types.ContainerImage); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.ContainerImage)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ContainerImageRepo_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type ContainerImageRepo_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
func (_e *ContainerImageRepo_Expecter) List() *ContainerImageRepo_List_Call {
	return &ContainerImageRepo_List_Call{Call: _e.mock.On("List")}
}

func (_c *ContainerImageRepo_List_Call) Run(run func()) *ContainerImageRepo_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ContainerImageRepo_List_Call) Return(_a0 []types.ContainerImage, _a1 error) *ContainerImageRepo_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ContainerImageRepo_List_Call) RunAndReturn(run func() ([]types.ContainerImage, error)) *ContainerImageRepo_List_Call {
	_c.Call.Return(run)
	return _c
}

// Prune provides a mock function with no fields
func (_m *ContainerImageRepo) Prune() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Prune")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ContainerImageRepo_Prune_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Prune'
type ContainerImageRepo_Prune_Call struct {
	*mock.Call
}

// Prune is a helper method to define mock.On call
func (_e *ContainerImageRepo_Expecter) Prune() *ContainerImageRepo_Prune_Call {
	return &ContainerImageRepo_Prune_Call{Call: _e.mock.On("Prune")}
}

func (_c *ContainerImageRepo_Prune_Call) Run(run func()) *ContainerImageRepo_Prune_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ContainerImageRepo_Prune_Call) Return(_a0 error) *ContainerImageRepo_Prune_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ContainerImageRepo_Prune_Call) RunAndReturn(run func() error) *ContainerImageRepo_Prune_Call {
	_c.Call.Return(run)
	return _c
}

// Pull provides a mock function with given fields: req
func (_m *ContainerImageRepo) Pull(req *request.ContainerImagePull) error {
	ret := _m.Called(req)

	if len(ret) == 0 {
		panic("no return value specified for Pull")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*request.ContainerImagePull) error); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ContainerImageRepo_Pull_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Pull'
type ContainerImageRepo_Pull_Call struct {
	*mock.Call
}

// Pull is a helper method to define mock.On call
//   - req *request.ContainerImagePull
func (_e *ContainerImageRepo_Expecter) Pull(req interface{}) *ContainerImageRepo_Pull_Call {
	return &ContainerImageRepo_Pull_Call{Call: _e.mock.On("Pull", req)}
}

func (_c *ContainerImageRepo_Pull_Call) Run(run func(req *request.ContainerImagePull)) *ContainerImageRepo_Pull_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*request.ContainerImagePull))
	})
	return _c
}

func (_c *ContainerImageRepo_Pull_Call) Return(_a0 error) *ContainerImageRepo_Pull_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ContainerImageRepo_Pull_Call) RunAndReturn(run func(*request.ContainerImagePull) error) *ContainerImageRepo_Pull_Call {
	_c.Call.Return(run)
	return _c
}

// Remove provides a mock function with given fields: id
func (_m *ContainerImageRepo) Remove(id string) error {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for Remove")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ContainerImageRepo_Remove_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Remove'
type ContainerImageRepo_Remove_Call struct {
	*mock.Call
}

// Remove is a helper method to define mock.On call
//   - id string
func (_e *ContainerImageRepo_Expecter) Remove(id interface{}) *ContainerImageRepo_Remove_Call {
	return &ContainerImageRepo_Remove_Call{Call: _e.mock.On("Remove", id)}
}

func (_c *ContainerImageRepo_Remove_Call) Run(run func(id string)) *ContainerImageRepo_Remove_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *ContainerImageRepo_Remove_Call) Return(_a0 error) *ContainerImageRepo_Remove_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ContainerImageRepo_Remove_Call) RunAndReturn(run func(string) error) *ContainerImageRepo_Remove_Call {
	_c.Call.Return(run)
	return _c
}

// NewContainerImageRepo creates a new instance of ContainerImageRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewContainerImageRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *ContainerImageRepo {
	mock := &ContainerImageRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
