// Code generated by MockGen. DO NOT EDIT.
// Source: handler.go

// Package mock_grpc is a generated GoMock package.
package mock_grpc

import (
	pb "anti_bruteforce/internal/delivery/grpc/pb"
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockUseCaseI is a mock of UseCaseI interface.
type MockUseCaseI struct {
	ctrl     *gomock.Controller
	recorder *MockUseCaseIMockRecorder
}

// MockUseCaseIMockRecorder is the mock recorder for MockUseCaseI.
type MockUseCaseIMockRecorder struct {
	mock *MockUseCaseI
}

// NewMockUseCaseI creates a new mock instance.
func NewMockUseCaseI(ctrl *gomock.Controller) *MockUseCaseI {
	mock := &MockUseCaseI{ctrl: ctrl}
	mock.recorder = &MockUseCaseIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUseCaseI) EXPECT() *MockUseCaseIMockRecorder {
	return m.recorder
}

// AddToBlackList mocks base method.
func (m *MockUseCaseI) AddToBlackList(ctx context.Context, subnet string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToBlackList", ctx, subnet)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddToBlackList indicates an expected call of AddToBlackList.
func (mr *MockUseCaseIMockRecorder) AddToBlackList(ctx, subnet interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToBlackList", reflect.TypeOf((*MockUseCaseI)(nil).AddToBlackList), ctx, subnet)
}

// AddToWhiteList mocks base method.
func (m *MockUseCaseI) AddToWhiteList(ctx context.Context, subnet string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToWhiteList", ctx, subnet)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddToWhiteList indicates an expected call of AddToWhiteList.
func (mr *MockUseCaseIMockRecorder) AddToWhiteList(ctx, subnet interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToWhiteList", reflect.TypeOf((*MockUseCaseI)(nil).AddToWhiteList), ctx, subnet)
}

// CheckAuth mocks base method.
func (m *MockUseCaseI) CheckAuth(ctx context.Context, in *pb.AuthCheckRequest) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckAuth", ctx, in)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckAuth indicates an expected call of CheckAuth.
func (mr *MockUseCaseIMockRecorder) CheckAuth(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAuth", reflect.TypeOf((*MockUseCaseI)(nil).CheckAuth), ctx, in)
}

// ClearLists mocks base method.
func (m *MockUseCaseI) ClearLists(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClearLists", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// ClearLists indicates an expected call of ClearLists.
func (mr *MockUseCaseIMockRecorder) ClearLists(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearLists", reflect.TypeOf((*MockUseCaseI)(nil).ClearLists), ctx)
}

// RemoveFromBlackList mocks base method.
func (m *MockUseCaseI) RemoveFromBlackList(ctx context.Context, subnet string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveFromBlackList", ctx, subnet)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveFromBlackList indicates an expected call of RemoveFromBlackList.
func (mr *MockUseCaseIMockRecorder) RemoveFromBlackList(ctx, subnet interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveFromBlackList", reflect.TypeOf((*MockUseCaseI)(nil).RemoveFromBlackList), ctx, subnet)
}

// RemoveFromWhiteList mocks base method.
func (m *MockUseCaseI) RemoveFromWhiteList(ctx context.Context, subnet string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveFromWhiteList", ctx, subnet)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveFromWhiteList indicates an expected call of RemoveFromWhiteList.
func (mr *MockUseCaseIMockRecorder) RemoveFromWhiteList(ctx, subnet interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveFromWhiteList", reflect.TypeOf((*MockUseCaseI)(nil).RemoveFromWhiteList), ctx, subnet)
}

// Reset mocks base method.
func (m *MockUseCaseI) Reset(ctx context.Context, in *pb.AuthCheckRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reset", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// Reset indicates an expected call of Reset.
func (mr *MockUseCaseIMockRecorder) Reset(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reset", reflect.TypeOf((*MockUseCaseI)(nil).Reset), ctx, in)
}