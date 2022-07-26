// Code generated by MockGen. DO NOT EDIT.
// Source: .\internal\app\interface\grpc_interfaces.go

// Package auth_mocks is a generated GoMock package.
package auth_mocks

import (
	reflect "reflect"
	entity "team3-task/internal/entity"

	gomock "github.com/golang/mock/gomock"
)

// MockAuthAccessChecker is a mock of AuthAccessChecker interface.
type MockAuthAccessChecker struct {
	ctrl     *gomock.Controller
	recorder *MockAuthAccessCheckerMockRecorder
}

// MockAuthAccessCheckerMockRecorder is the mock recorder for MockAuthAccessChecker.
type MockAuthAccessCheckerMockRecorder struct {
	mock *MockAuthAccessChecker
}

// NewMockAuthAccessChecker creates a new mock instance.
func NewMockAuthAccessChecker(ctrl *gomock.Controller) *MockAuthAccessChecker {
	mock := &MockAuthAccessChecker{ctrl: ctrl}
	mock.recorder = &MockAuthAccessCheckerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthAccessChecker) EXPECT() *MockAuthAccessCheckerMockRecorder {
	return m.recorder
}

// CheckAccess mocks base method.
func (m *MockAuthAccessChecker) CheckAccess(arg0 *entity.AuthRequest) (entity.AuthResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckAccess", arg0)
	ret0, _ := ret[0].(entity.AuthResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckAccess indicates an expected call of CheckAccess.
func (mr *MockAuthAccessCheckerMockRecorder) CheckAccess(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAccess", reflect.TypeOf((*MockAuthAccessChecker)(nil).CheckAccess), arg0)
}
