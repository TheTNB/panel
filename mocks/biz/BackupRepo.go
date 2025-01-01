// Code generated by mockery. DO NOT EDIT.

package biz

import (
	mock "github.com/stretchr/testify/mock"
	biz "github.com/tnb-labs/panel/internal/biz"

	types "github.com/tnb-labs/panel/pkg/types"
)

// BackupRepo is an autogenerated mock type for the BackupRepo type
type BackupRepo struct {
	mock.Mock
}

type BackupRepo_Expecter struct {
	mock *mock.Mock
}

func (_m *BackupRepo) EXPECT() *BackupRepo_Expecter {
	return &BackupRepo_Expecter{mock: &_m.Mock}
}

// ClearExpired provides a mock function with given fields: path, prefix, save
func (_m *BackupRepo) ClearExpired(path string, prefix string, save int) error {
	ret := _m.Called(path, prefix, save)

	if len(ret) == 0 {
		panic("no return value specified for ClearExpired")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, int) error); ok {
		r0 = rf(path, prefix, save)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BackupRepo_ClearExpired_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ClearExpired'
type BackupRepo_ClearExpired_Call struct {
	*mock.Call
}

// ClearExpired is a helper method to define mock.On call
//   - path string
//   - prefix string
//   - save int
func (_e *BackupRepo_Expecter) ClearExpired(path interface{}, prefix interface{}, save interface{}) *BackupRepo_ClearExpired_Call {
	return &BackupRepo_ClearExpired_Call{Call: _e.mock.On("ClearExpired", path, prefix, save)}
}

func (_c *BackupRepo_ClearExpired_Call) Run(run func(path string, prefix string, save int)) *BackupRepo_ClearExpired_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(int))
	})
	return _c
}

func (_c *BackupRepo_ClearExpired_Call) Return(_a0 error) *BackupRepo_ClearExpired_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *BackupRepo_ClearExpired_Call) RunAndReturn(run func(string, string, int) error) *BackupRepo_ClearExpired_Call {
	_c.Call.Return(run)
	return _c
}

// Create provides a mock function with given fields: typ, target, path
func (_m *BackupRepo) Create(typ biz.BackupType, target string, path ...string) error {
	_va := make([]interface{}, len(path))
	for _i := range path {
		_va[_i] = path[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, typ, target)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(biz.BackupType, string, ...string) error); ok {
		r0 = rf(typ, target, path...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BackupRepo_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type BackupRepo_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - typ biz.BackupType
//   - target string
//   - path ...string
func (_e *BackupRepo_Expecter) Create(typ interface{}, target interface{}, path ...interface{}) *BackupRepo_Create_Call {
	return &BackupRepo_Create_Call{Call: _e.mock.On("Create",
		append([]interface{}{typ, target}, path...)...)}
}

func (_c *BackupRepo_Create_Call) Run(run func(typ biz.BackupType, target string, path ...string)) *BackupRepo_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]string, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(string)
			}
		}
		run(args[0].(biz.BackupType), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *BackupRepo_Create_Call) Return(_a0 error) *BackupRepo_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *BackupRepo_Create_Call) RunAndReturn(run func(biz.BackupType, string, ...string) error) *BackupRepo_Create_Call {
	_c.Call.Return(run)
	return _c
}

// CutoffLog provides a mock function with given fields: path, target
func (_m *BackupRepo) CutoffLog(path string, target string) error {
	ret := _m.Called(path, target)

	if len(ret) == 0 {
		panic("no return value specified for CutoffLog")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(path, target)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BackupRepo_CutoffLog_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CutoffLog'
type BackupRepo_CutoffLog_Call struct {
	*mock.Call
}

// CutoffLog is a helper method to define mock.On call
//   - path string
//   - target string
func (_e *BackupRepo_Expecter) CutoffLog(path interface{}, target interface{}) *BackupRepo_CutoffLog_Call {
	return &BackupRepo_CutoffLog_Call{Call: _e.mock.On("CutoffLog", path, target)}
}

func (_c *BackupRepo_CutoffLog_Call) Run(run func(path string, target string)) *BackupRepo_CutoffLog_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *BackupRepo_CutoffLog_Call) Return(_a0 error) *BackupRepo_CutoffLog_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *BackupRepo_CutoffLog_Call) RunAndReturn(run func(string, string) error) *BackupRepo_CutoffLog_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: typ, name
func (_m *BackupRepo) Delete(typ biz.BackupType, name string) error {
	ret := _m.Called(typ, name)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(biz.BackupType, string) error); ok {
		r0 = rf(typ, name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BackupRepo_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type BackupRepo_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - typ biz.BackupType
//   - name string
func (_e *BackupRepo_Expecter) Delete(typ interface{}, name interface{}) *BackupRepo_Delete_Call {
	return &BackupRepo_Delete_Call{Call: _e.mock.On("Delete", typ, name)}
}

func (_c *BackupRepo_Delete_Call) Run(run func(typ biz.BackupType, name string)) *BackupRepo_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(biz.BackupType), args[1].(string))
	})
	return _c
}

func (_c *BackupRepo_Delete_Call) Return(_a0 error) *BackupRepo_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *BackupRepo_Delete_Call) RunAndReturn(run func(biz.BackupType, string) error) *BackupRepo_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// FixPanel provides a mock function with no fields
func (_m *BackupRepo) FixPanel() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for FixPanel")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BackupRepo_FixPanel_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FixPanel'
type BackupRepo_FixPanel_Call struct {
	*mock.Call
}

// FixPanel is a helper method to define mock.On call
func (_e *BackupRepo_Expecter) FixPanel() *BackupRepo_FixPanel_Call {
	return &BackupRepo_FixPanel_Call{Call: _e.mock.On("FixPanel")}
}

func (_c *BackupRepo_FixPanel_Call) Run(run func()) *BackupRepo_FixPanel_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *BackupRepo_FixPanel_Call) Return(_a0 error) *BackupRepo_FixPanel_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *BackupRepo_FixPanel_Call) RunAndReturn(run func() error) *BackupRepo_FixPanel_Call {
	_c.Call.Return(run)
	return _c
}

// GetPath provides a mock function with given fields: typ
func (_m *BackupRepo) GetPath(typ biz.BackupType) (string, error) {
	ret := _m.Called(typ)

	if len(ret) == 0 {
		panic("no return value specified for GetPath")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(biz.BackupType) (string, error)); ok {
		return rf(typ)
	}
	if rf, ok := ret.Get(0).(func(biz.BackupType) string); ok {
		r0 = rf(typ)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(biz.BackupType) error); ok {
		r1 = rf(typ)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BackupRepo_GetPath_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPath'
type BackupRepo_GetPath_Call struct {
	*mock.Call
}

// GetPath is a helper method to define mock.On call
//   - typ biz.BackupType
func (_e *BackupRepo_Expecter) GetPath(typ interface{}) *BackupRepo_GetPath_Call {
	return &BackupRepo_GetPath_Call{Call: _e.mock.On("GetPath", typ)}
}

func (_c *BackupRepo_GetPath_Call) Run(run func(typ biz.BackupType)) *BackupRepo_GetPath_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(biz.BackupType))
	})
	return _c
}

func (_c *BackupRepo_GetPath_Call) Return(_a0 string, _a1 error) *BackupRepo_GetPath_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *BackupRepo_GetPath_Call) RunAndReturn(run func(biz.BackupType) (string, error)) *BackupRepo_GetPath_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: typ
func (_m *BackupRepo) List(typ biz.BackupType) ([]*types.BackupFile, error) {
	ret := _m.Called(typ)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 []*types.BackupFile
	var r1 error
	if rf, ok := ret.Get(0).(func(biz.BackupType) ([]*types.BackupFile, error)); ok {
		return rf(typ)
	}
	if rf, ok := ret.Get(0).(func(biz.BackupType) []*types.BackupFile); ok {
		r0 = rf(typ)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.BackupFile)
		}
	}

	if rf, ok := ret.Get(1).(func(biz.BackupType) error); ok {
		r1 = rf(typ)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BackupRepo_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type BackupRepo_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - typ biz.BackupType
func (_e *BackupRepo_Expecter) List(typ interface{}) *BackupRepo_List_Call {
	return &BackupRepo_List_Call{Call: _e.mock.On("List", typ)}
}

func (_c *BackupRepo_List_Call) Run(run func(typ biz.BackupType)) *BackupRepo_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(biz.BackupType))
	})
	return _c
}

func (_c *BackupRepo_List_Call) Return(_a0 []*types.BackupFile, _a1 error) *BackupRepo_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *BackupRepo_List_Call) RunAndReturn(run func(biz.BackupType) ([]*types.BackupFile, error)) *BackupRepo_List_Call {
	_c.Call.Return(run)
	return _c
}

// Restore provides a mock function with given fields: typ, backup, target
func (_m *BackupRepo) Restore(typ biz.BackupType, backup string, target string) error {
	ret := _m.Called(typ, backup, target)

	if len(ret) == 0 {
		panic("no return value specified for Restore")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(biz.BackupType, string, string) error); ok {
		r0 = rf(typ, backup, target)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BackupRepo_Restore_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Restore'
type BackupRepo_Restore_Call struct {
	*mock.Call
}

// Restore is a helper method to define mock.On call
//   - typ biz.BackupType
//   - backup string
//   - target string
func (_e *BackupRepo_Expecter) Restore(typ interface{}, backup interface{}, target interface{}) *BackupRepo_Restore_Call {
	return &BackupRepo_Restore_Call{Call: _e.mock.On("Restore", typ, backup, target)}
}

func (_c *BackupRepo_Restore_Call) Run(run func(typ biz.BackupType, backup string, target string)) *BackupRepo_Restore_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(biz.BackupType), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *BackupRepo_Restore_Call) Return(_a0 error) *BackupRepo_Restore_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *BackupRepo_Restore_Call) RunAndReturn(run func(biz.BackupType, string, string) error) *BackupRepo_Restore_Call {
	_c.Call.Return(run)
	return _c
}

// UpdatePanel provides a mock function with given fields: version, url, checksum
func (_m *BackupRepo) UpdatePanel(version string, url string, checksum string) error {
	ret := _m.Called(version, url, checksum)

	if len(ret) == 0 {
		panic("no return value specified for UpdatePanel")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string) error); ok {
		r0 = rf(version, url, checksum)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BackupRepo_UpdatePanel_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdatePanel'
type BackupRepo_UpdatePanel_Call struct {
	*mock.Call
}

// UpdatePanel is a helper method to define mock.On call
//   - version string
//   - url string
//   - checksum string
func (_e *BackupRepo_Expecter) UpdatePanel(version interface{}, url interface{}, checksum interface{}) *BackupRepo_UpdatePanel_Call {
	return &BackupRepo_UpdatePanel_Call{Call: _e.mock.On("UpdatePanel", version, url, checksum)}
}

func (_c *BackupRepo_UpdatePanel_Call) Run(run func(version string, url string, checksum string)) *BackupRepo_UpdatePanel_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *BackupRepo_UpdatePanel_Call) Return(_a0 error) *BackupRepo_UpdatePanel_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *BackupRepo_UpdatePanel_Call) RunAndReturn(run func(string, string, string) error) *BackupRepo_UpdatePanel_Call {
	_c.Call.Return(run)
	return _c
}

// NewBackupRepo creates a new instance of BackupRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBackupRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *BackupRepo {
	mock := &BackupRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
