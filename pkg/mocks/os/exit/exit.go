// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/mocks/os/models (interfaces: Exiter)

// Package mocked_exiter is a generated GoMock package.
package mocked_exiter

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockExiter is a mock of Exiter interface.
type MockExiter struct {
	ctrl     *gomock.Controller
	recorder *MockExiterMockRecorder
}

// MockExiterMockRecorder is the mock recorder for MockExiter.
type MockExiterMockRecorder struct {
	mock *MockExiter
}

// NewMockExiter creates a new mock instance.
func NewMockExiter(ctrl *gomock.Controller) *MockExiter {
	mock := &MockExiter{ctrl: ctrl}
	mock.recorder = &MockExiterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExiter) EXPECT() *MockExiterMockRecorder {
	return m.recorder
}

// Exit mocks base method.
func (m *MockExiter) Exit(arg0 int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Exit", arg0)
}

// Exit indicates an expected call of Exit.
func (mr *MockExiterMockRecorder) Exit(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exit", reflect.TypeOf((*MockExiter)(nil).Exit), arg0)
}
