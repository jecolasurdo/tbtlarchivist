// Code generated by MockGen. DO NOT EDIT.
// Source: go/internal/accessors/messagebus/acknowledger/acknack.go

// Package mock_acknowledger is a generated GoMock package.
package mock_acknowledger

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockAckNack is a mock of AckNack interface
type MockAckNack struct {
	ctrl     *gomock.Controller
	recorder *MockAckNackMockRecorder
}

// MockAckNackMockRecorder is the mock recorder for MockAckNack
type MockAckNackMockRecorder struct {
	mock *MockAckNack
}

// NewMockAckNack creates a new mock instance
func NewMockAckNack(ctrl *gomock.Controller) *MockAckNack {
	mock := &MockAckNack{ctrl: ctrl}
	mock.recorder = &MockAckNackMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAckNack) EXPECT() *MockAckNackMockRecorder {
	return m.recorder
}

// Ack mocks base method
func (m *MockAckNack) Ack() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ack")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ack indicates an expected call of Ack
func (mr *MockAckNackMockRecorder) Ack() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ack", reflect.TypeOf((*MockAckNack)(nil).Ack))
}

// Nack mocks base method
func (m *MockAckNack) Nack(requeue bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Nack", requeue)
	ret0, _ := ret[0].(error)
	return ret0
}

// Nack indicates an expected call of Nack
func (mr *MockAckNackMockRecorder) Nack(requeue interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Nack", reflect.TypeOf((*MockAckNack)(nil).Nack), requeue)
}
