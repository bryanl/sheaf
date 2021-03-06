// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/bryanl/sheaf/pkg/sheaf (interfaces: ArtifactsService)

// Package mocks is a generated GoMock package.
package mocks

import (
	sheaf "github.com/bryanl/sheaf/pkg/sheaf"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockArtifactsService is a mock of ArtifactsService interface
type MockArtifactsService struct {
	ctrl     *gomock.Controller
	recorder *MockArtifactsServiceMockRecorder
}

// MockArtifactsServiceMockRecorder is the mock recorder for MockArtifactsService
type MockArtifactsServiceMockRecorder struct {
	mock *MockArtifactsService
}

// NewMockArtifactsService creates a new mock instance
func NewMockArtifactsService(ctrl *gomock.Controller) *MockArtifactsService {
	mock := &MockArtifactsService{ctrl: ctrl}
	mock.recorder = &MockArtifactsServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockArtifactsService) EXPECT() *MockArtifactsServiceMockRecorder {
	return m.recorder
}

// Image mocks base method
func (m *MockArtifactsService) Image() sheaf.ImageService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Image")
	ret0, _ := ret[0].(sheaf.ImageService)
	return ret0
}

// Image indicates an expected call of Image
func (mr *MockArtifactsServiceMockRecorder) Image() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Image", reflect.TypeOf((*MockArtifactsService)(nil).Image))
}

// Index mocks base method
func (m *MockArtifactsService) Index() ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Index")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Index indicates an expected call of Index
func (mr *MockArtifactsServiceMockRecorder) Index() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Index", reflect.TypeOf((*MockArtifactsService)(nil).Index))
}
