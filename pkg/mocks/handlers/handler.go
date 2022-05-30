// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/v4/pkg/models/handlers (interfaces: Server)

// Package mocked_transport_handlers is a generated GoMock package.
package mocked_transport_handlers

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockServer is a mock of Server interface.
type MockServer struct {
	ctrl     *gomock.Controller
	recorder *MockServerMockRecorder
}

// MockServerMockRecorder is the mock recorder for MockServer.
type MockServerMockRecorder struct {
	mock *MockServer
}

// NewMockServer creates a new mock instance.
func NewMockServer(ctrl *gomock.Controller) *MockServer {
	mock := &MockServer{ctrl: ctrl}
	mock.recorder = &MockServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServer) EXPECT() *MockServerMockRecorder {
	return m.recorder
}

// Handle mocks base method.
func (m *MockServer) Handle() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Handle")
}

// Handle indicates an expected call of Handle.
func (mr *MockServerMockRecorder) Handle() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockServer)(nil).Handle))
}

// Start mocks base method.
func (m *MockServer) Start() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockServerMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockServer)(nil).Start))
}
