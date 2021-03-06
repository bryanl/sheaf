// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/bryanl/sheaf/pkg/sheaf (interfaces: ImageService)

// Package mocks is a generated GoMock package.
package mocks

import (
	sheaf "github.com/bryanl/sheaf/pkg/sheaf"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockImageService is a mock of ImageService interface
type MockImageService struct {
	ctrl     *gomock.Controller
	recorder *MockImageServiceMockRecorder
}

// MockImageServiceMockRecorder is the mock recorder for MockImageService
type MockImageServiceMockRecorder struct {
	mock *MockImageService
}

// NewMockImageService creates a new mock instance
func NewMockImageService(ctrl *gomock.Controller) *MockImageService {
	mock := &MockImageService{ctrl: ctrl}
	mock.recorder = &MockImageServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockImageService) EXPECT() *MockImageServiceMockRecorder {
	return m.recorder
}

// List mocks base method
func (m *MockImageService) List() ([]sheaf.BundleImage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List")
	ret0, _ := ret[0].([]sheaf.BundleImage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockImageServiceMockRecorder) List() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockImageService)(nil).List))
}
