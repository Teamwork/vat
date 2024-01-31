// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/teamwork/vat (interfaces: ViesServiceInterface)

// Package mocks is a generated GoMock package.
package mocks

import (
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockViesServiceInterface is a mock of ViesServiceInterface interface.
type MockViesServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockViesServiceInterfaceMockRecorder
}

// MockViesServiceInterfaceMockRecorder is the mock recorder for MockViesServiceInterface.
type MockViesServiceInterfaceMockRecorder struct {
	mock *MockViesServiceInterface
}

// NewMockViesServiceInterface creates a new mock instance.
func NewMockViesServiceInterface(ctrl *gomock.Controller) *MockViesServiceInterface {
	mock := &MockViesServiceInterface{ctrl: ctrl}
	mock.recorder = &MockViesServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockViesServiceInterface) EXPECT() *MockViesServiceInterfaceMockRecorder {
	return m.recorder
}

// Lookup mocks base method.
func (m *MockViesServiceInterface) Lookup(arg0 string) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Lookup", arg0)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Lookup indicates an expected call of Lookup.
func (mr *MockViesServiceInterfaceMockRecorder) Lookup(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Lookup", reflect.TypeOf((*MockViesServiceInterface)(nil).Lookup), arg0)
}
