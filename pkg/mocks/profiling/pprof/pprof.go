// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/v4/pkg/mocks/profiling/models (interfaces: Profiler)

// Package mocked_profiler is a generated GoMock package.
package mocked_profiler

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockProfiler is a mock of Profiler interface.
type MockProfiler struct {
	ctrl     *gomock.Controller
	recorder *MockProfilerMockRecorder
}

// MockProfilerMockRecorder is the mock recorder for MockProfiler.
type MockProfilerMockRecorder struct {
	mock *MockProfiler
}

// NewMockProfiler creates a new mock instance.
func NewMockProfiler(ctrl *gomock.Controller) *MockProfiler {
	mock := &MockProfiler{ctrl: ctrl}
	mock.recorder = &MockProfilerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProfiler) EXPECT() *MockProfilerMockRecorder {
	return m.recorder
}

// Stop mocks base method.
func (m *MockProfiler) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop.
func (mr *MockProfilerMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockProfiler)(nil).Stop))
}
