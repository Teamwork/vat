// Code generated by MockGen. DO NOT EDIT.
// Source: vies_service.go

// Package vat is a generated GoMock package.
package vat

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockLookupServiceInterface is a mock of LookupServiceInterface interface.
type MockLookupServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockLookupServiceInterfaceMockRecorder
}

// MockLookupServiceInterfaceMockRecorder is the mock recorder for MockLookupServiceInterface.
type MockLookupServiceInterfaceMockRecorder struct {
	mock *MockLookupServiceInterface
}

// NewMockLookupServiceInterface creates a new mock instance.
func NewMockLookupServiceInterface(ctrl *gomock.Controller) *MockLookupServiceInterface {
	mock := &MockLookupServiceInterface{ctrl: ctrl}
	mock.recorder = &MockLookupServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLookupServiceInterface) EXPECT() *MockLookupServiceInterfaceMockRecorder {
	return m.recorder
}

// Validate mocks base method.
func (m *MockLookupServiceInterface) Validate(vatNumber string, opts ValidatorOpts) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate", vatNumber, opts)
	ret0, _ := ret[0].(error)
	return ret0
}

// Validate indicates an expected call of Validate.
func (mr *MockLookupServiceInterfaceMockRecorder) Validate(vatNumber, opts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockLookupServiceInterface)(nil).Validate), vatNumber, opts)
}
