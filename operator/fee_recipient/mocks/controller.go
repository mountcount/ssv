// Code generated by MockGen. DO NOT EDIT.
// Source: ./controller.go
//
// Generated by this command:
//
//	mockgen -package=mocks -destination=./mocks/controller.go -source=./controller.go
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	zap "go.uber.org/zap"
)

// MockRecipientController is a mock of RecipientController interface.
type MockRecipientController struct {
	ctrl     *gomock.Controller
	recorder *MockRecipientControllerMockRecorder
}

// MockRecipientControllerMockRecorder is the mock recorder for MockRecipientController.
type MockRecipientControllerMockRecorder struct {
	mock *MockRecipientController
}

// NewMockRecipientController creates a new mock instance.
func NewMockRecipientController(ctrl *gomock.Controller) *MockRecipientController {
	mock := &MockRecipientController{ctrl: ctrl}
	mock.recorder = &MockRecipientControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRecipientController) EXPECT() *MockRecipientControllerMockRecorder {
	return m.recorder
}

// Start mocks base method.
func (m *MockRecipientController) Start(logger *zap.Logger) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Start", logger)
}

// Start indicates an expected call of Start.
func (mr *MockRecipientControllerMockRecorder) Start(logger any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockRecipientController)(nil).Start), logger)
}
