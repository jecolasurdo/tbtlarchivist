// Code generated by MockGen. DO NOT EDIT.
// Source: go/internal/accessors/analyst/command.go

// Package mock_analyst is a generated GoMock package.
package mock_analyst

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	analyst "github.com/jecolasurdo/tbtlarchivist/go/internal/accessors/analyst"
	reflect "reflect"
)

// MockCommandBuilder is a mock of CommandBuilder interface
type MockCommandBuilder struct {
	ctrl     *gomock.Controller
	recorder *MockCommandBuilderMockRecorder
}

// MockCommandBuilderMockRecorder is the mock recorder for MockCommandBuilder
type MockCommandBuilderMockRecorder struct {
	mock *MockCommandBuilder
}

// NewMockCommandBuilder creates a new mock instance
func NewMockCommandBuilder(ctrl *gomock.Controller) *MockCommandBuilder {
	mock := &MockCommandBuilder{ctrl: ctrl}
	mock.recorder = &MockCommandBuilderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCommandBuilder) EXPECT() *MockCommandBuilderMockRecorder {
	return m.recorder
}

// CommandContext mocks base method
func (m *MockCommandBuilder) CommandContext(arg0 context.Context, arg1 string, arg2 ...string) analyst.Command {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CommandContext", varargs...)
	ret0, _ := ret[0].(analyst.Command)
	return ret0
}

// CommandContext indicates an expected call of CommandContext
func (mr *MockCommandBuilderMockRecorder) CommandContext(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommandContext", reflect.TypeOf((*MockCommandBuilder)(nil).CommandContext), varargs...)
}

// MockCommand is a mock of Command interface
type MockCommand struct {
	ctrl     *gomock.Controller
	recorder *MockCommandMockRecorder
}

// MockCommandMockRecorder is the mock recorder for MockCommand
type MockCommandMockRecorder struct {
	mock *MockCommand
}

// NewMockCommand creates a new mock instance
func NewMockCommand(ctrl *gomock.Controller) *MockCommand {
	mock := &MockCommand{ctrl: ctrl}
	mock.recorder = &MockCommandMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCommand) EXPECT() *MockCommandMockRecorder {
	return m.recorder
}

// StdoutPipe mocks base method
func (m *MockCommand) StdoutPipe() (analyst.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StdoutPipe")
	ret0, _ := ret[0].(analyst.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StdoutPipe indicates an expected call of StdoutPipe
func (mr *MockCommandMockRecorder) StdoutPipe() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StdoutPipe", reflect.TypeOf((*MockCommand)(nil).StdoutPipe))
}

// StdinPipe mocks base method
func (m *MockCommand) StdinPipe() (analyst.WriteCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StdinPipe")
	ret0, _ := ret[0].(analyst.WriteCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StdinPipe indicates an expected call of StdinPipe
func (mr *MockCommandMockRecorder) StdinPipe() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StdinPipe", reflect.TypeOf((*MockCommand)(nil).StdinPipe))
}

// StderrPipe mocks base method
func (m *MockCommand) StderrPipe() (analyst.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StderrPipe")
	ret0, _ := ret[0].(analyst.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StderrPipe indicates an expected call of StderrPipe
func (mr *MockCommandMockRecorder) StderrPipe() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StderrPipe", reflect.TypeOf((*MockCommand)(nil).StderrPipe))
}

// Start mocks base method
func (m *MockCommand) Start() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start
func (mr *MockCommandMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockCommand)(nil).Start))
}

// Wait mocks base method
func (m *MockCommand) Wait() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Wait")
	ret0, _ := ret[0].(error)
	return ret0
}

// Wait indicates an expected call of Wait
func (mr *MockCommandMockRecorder) Wait() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Wait", reflect.TypeOf((*MockCommand)(nil).Wait))
}

// MockReadCloser is a mock of ReadCloser interface
type MockReadCloser struct {
	ctrl     *gomock.Controller
	recorder *MockReadCloserMockRecorder
}

// MockReadCloserMockRecorder is the mock recorder for MockReadCloser
type MockReadCloserMockRecorder struct {
	mock *MockReadCloser
}

// NewMockReadCloser creates a new mock instance
func NewMockReadCloser(ctrl *gomock.Controller) *MockReadCloser {
	mock := &MockReadCloser{ctrl: ctrl}
	mock.recorder = &MockReadCloserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockReadCloser) EXPECT() *MockReadCloserMockRecorder {
	return m.recorder
}

// Read mocks base method
func (m *MockReadCloser) Read(p []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", p)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read
func (mr *MockReadCloserMockRecorder) Read(p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockReadCloser)(nil).Read), p)
}

// Close mocks base method
func (m *MockReadCloser) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockReadCloserMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockReadCloser)(nil).Close))
}

// MockWriteCloser is a mock of WriteCloser interface
type MockWriteCloser struct {
	ctrl     *gomock.Controller
	recorder *MockWriteCloserMockRecorder
}

// MockWriteCloserMockRecorder is the mock recorder for MockWriteCloser
type MockWriteCloserMockRecorder struct {
	mock *MockWriteCloser
}

// NewMockWriteCloser creates a new mock instance
func NewMockWriteCloser(ctrl *gomock.Controller) *MockWriteCloser {
	mock := &MockWriteCloser{ctrl: ctrl}
	mock.recorder = &MockWriteCloserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockWriteCloser) EXPECT() *MockWriteCloserMockRecorder {
	return m.recorder
}

// Write mocks base method
func (m *MockWriteCloser) Write(p []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", p)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Write indicates an expected call of Write
func (mr *MockWriteCloserMockRecorder) Write(p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockWriteCloser)(nil).Write), p)
}

// Close mocks base method
func (m *MockWriteCloser) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockWriteCloserMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockWriteCloser)(nil).Close))
}
