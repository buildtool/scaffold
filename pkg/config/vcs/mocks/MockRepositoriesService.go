// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/buildtool/scaffold/pkg/config (interfaces: RepositoriesService)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	github "github.com/google/go-github/v28/github"
	reflect "reflect"
)

// MockRepositoriesService is a mock of RepositoriesService interface
type MockRepositoriesService struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoriesServiceMockRecorder
}

// MockRepositoriesServiceMockRecorder is the mock recorder for MockRepositoriesService
type MockRepositoriesServiceMockRecorder struct {
	mock *MockRepositoriesService
}

// NewMockRepositoriesService creates a new mock instance
func NewMockRepositoriesService(ctrl *gomock.Controller) *MockRepositoriesService {
	mock := &MockRepositoriesService{ctrl: ctrl}
	mock.recorder = &MockRepositoriesServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepositoriesService) EXPECT() *MockRepositoriesServiceMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockRepositoriesService) Create(arg0 context.Context, arg1 string, arg2 *github.Repository) (*github.Repository, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1, arg2)
	ret0, _ := ret[0].(*github.Repository)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Create indicates an expected call of Create
func (mr *MockRepositoriesServiceMockRecorder) Create(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepositoriesService)(nil).Create), arg0, arg1, arg2)
}

// CreateHook mocks base method
func (m *MockRepositoriesService) CreateHook(arg0 context.Context, arg1, arg2 string, arg3 *github.Hook) (*github.Hook, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateHook", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*github.Hook)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateHook indicates an expected call of CreateHook
func (mr *MockRepositoriesServiceMockRecorder) CreateHook(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateHook", reflect.TypeOf((*MockRepositoriesService)(nil).CreateHook), arg0, arg1, arg2, arg3)
}

// UpdateBranchProtection mocks base method
func (m *MockRepositoriesService) UpdateBranchProtection(arg0 context.Context, arg1, arg2, arg3 string, arg4 *github.ProtectionRequest) (*github.Protection, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBranchProtection", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(*github.Protection)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// UpdateBranchProtection indicates an expected call of UpdateBranchProtection
func (mr *MockRepositoriesServiceMockRecorder) UpdateBranchProtection(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBranchProtection", reflect.TypeOf((*MockRepositoriesService)(nil).UpdateBranchProtection), arg0, arg1, arg2, arg3, arg4)
}
