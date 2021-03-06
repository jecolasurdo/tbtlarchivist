// Code generated by MockGen. DO NOT EDIT.
// Source: go/internal/accessors/messagebus/messagebus.go

// Package mock_messagebus is a generated GoMock package.
package mock_messagebus

import (
	gomock "github.com/golang/mock/gomock"
	messagebustypes "github.com/jecolasurdo/tbtlarchivist/go/internal/accessors/messagebus/messagebustypes"
	reflect "reflect"
)

// MockSender is a mock of Sender interface
type MockSender struct {
	ctrl     *gomock.Controller
	recorder *MockSenderMockRecorder
}

// MockSenderMockRecorder is the mock recorder for MockSender
type MockSenderMockRecorder struct {
	mock *MockSender
}

// NewMockSender creates a new mock instance
func NewMockSender(ctrl *gomock.Controller) *MockSender {
	mock := &MockSender{ctrl: ctrl}
	mock.recorder = &MockSenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSender) EXPECT() *MockSenderMockRecorder {
	return m.recorder
}

// Send mocks base method
func (m *MockSender) Send(arg0 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send
func (mr *MockSenderMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockSender)(nil).Send), arg0)
}

// Inspect mocks base method
func (m *MockSender) Inspect() (*messagebustypes.QueueInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Inspect")
	ret0, _ := ret[0].(*messagebustypes.QueueInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Inspect indicates an expected call of Inspect
func (mr *MockSenderMockRecorder) Inspect() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Inspect", reflect.TypeOf((*MockSender)(nil).Inspect))
}

// MockReceiver is a mock of Receiver interface
type MockReceiver struct {
	ctrl     *gomock.Controller
	recorder *MockReceiverMockRecorder
}

// MockReceiverMockRecorder is the mock recorder for MockReceiver
type MockReceiverMockRecorder struct {
	mock *MockReceiver
}

// NewMockReceiver creates a new mock instance
func NewMockReceiver(ctrl *gomock.Controller) *MockReceiver {
	mock := &MockReceiver{ctrl: ctrl}
	mock.recorder = &MockReceiverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockReceiver) EXPECT() *MockReceiverMockRecorder {
	return m.recorder
}

// Receive mocks base method
func (m *MockReceiver) Receive() (*messagebustypes.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Receive")
	ret0, _ := ret[0].(*messagebustypes.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Receive indicates an expected call of Receive
func (mr *MockReceiverMockRecorder) Receive() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Receive", reflect.TypeOf((*MockReceiver)(nil).Receive))
}

// MockSenderReceiver is a mock of SenderReceiver interface
type MockSenderReceiver struct {
	ctrl     *gomock.Controller
	recorder *MockSenderReceiverMockRecorder
}

// MockSenderReceiverMockRecorder is the mock recorder for MockSenderReceiver
type MockSenderReceiverMockRecorder struct {
	mock *MockSenderReceiver
}

// NewMockSenderReceiver creates a new mock instance
func NewMockSenderReceiver(ctrl *gomock.Controller) *MockSenderReceiver {
	mock := &MockSenderReceiver{ctrl: ctrl}
	mock.recorder = &MockSenderReceiverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSenderReceiver) EXPECT() *MockSenderReceiverMockRecorder {
	return m.recorder
}

// Send mocks base method
func (m *MockSenderReceiver) Send(arg0 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send
func (mr *MockSenderReceiverMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockSenderReceiver)(nil).Send), arg0)
}

// Inspect mocks base method
func (m *MockSenderReceiver) Inspect() (*messagebustypes.QueueInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Inspect")
	ret0, _ := ret[0].(*messagebustypes.QueueInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Inspect indicates an expected call of Inspect
func (mr *MockSenderReceiverMockRecorder) Inspect() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Inspect", reflect.TypeOf((*MockSenderReceiver)(nil).Inspect))
}

// Receive mocks base method
func (m *MockSenderReceiver) Receive() (*messagebustypes.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Receive")
	ret0, _ := ret[0].(*messagebustypes.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Receive indicates an expected call of Receive
func (mr *MockSenderReceiverMockRecorder) Receive() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Receive", reflect.TypeOf((*MockSenderReceiver)(nil).Receive))
}
