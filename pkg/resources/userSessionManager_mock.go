// Code generated by MockGen. DO NOT EDIT.
// Source: pkg\managers\userSessionManager.go

// Package resources is a generated GoMock package.
package resources

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	managers "github.com/kwhite17/Neighbors/pkg/managers"
	reflect "reflect"
)

// MockSessionManger is a mock of SessionManger interface
type MockSessionManger struct {
	ctrl     *gomock.Controller
	recorder *MockSessionMangerMockRecorder
}

// MockSessionMangerMockRecorder is the mock recorder for MockSessionManger
type MockSessionMangerMockRecorder struct {
	mock *MockSessionManger
}

// NewMockSessionManger creates a new mock instance
func NewMockSessionManger(ctrl *gomock.Controller) *MockSessionManger {
	mock := &MockSessionManger{ctrl: ctrl}
	mock.recorder = &MockSessionMangerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSessionManger) EXPECT() *MockSessionMangerMockRecorder {
	return m.recorder
}

// GetUserSession mocks base method
func (m *MockSessionManger) GetUserSession(ctx context.Context, sessionKey interface{}) (*managers.UserSession, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserSession", ctx, sessionKey)
	ret0, _ := ret[0].(*managers.UserSession)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserSession indicates an expected call of GetUserSession
func (mr *MockSessionMangerMockRecorder) GetUserSession(ctx, sessionKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserSession", reflect.TypeOf((*MockSessionManger)(nil).GetUserSession), ctx, sessionKey)
}

// WriteUserSession mocks base method
func (m *MockSessionManger) WriteUserSession(ctx context.Context, userID int64, userType managers.UserType) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteUserSession", ctx, userID, userType)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WriteUserSession indicates an expected call of WriteUserSession
func (mr *MockSessionMangerMockRecorder) WriteUserSession(ctx, userID, userType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteUserSession", reflect.TypeOf((*MockSessionManger)(nil).WriteUserSession), ctx, userID, userType)
}

// UpdateUserSession mocks base method
func (m *MockSessionManger) UpdateUserSession(ctx context.Context, userID, loginTime, lastSeenTime int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserSession", ctx, userID, loginTime, lastSeenTime)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserSession indicates an expected call of UpdateUserSession
func (mr *MockSessionMangerMockRecorder) UpdateUserSession(ctx, userID, loginTime, lastSeenTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserSession", reflect.TypeOf((*MockSessionManger)(nil).UpdateUserSession), ctx, userID, loginTime, lastSeenTime)
}

// DeleteUserSession mocks base method
func (m *MockSessionManger) DeleteUserSession(ctx context.Context, sessionKey interface{}) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserSession", ctx, sessionKey)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteUserSession indicates an expected call of DeleteUserSession
func (mr *MockSessionMangerMockRecorder) DeleteUserSession(ctx, sessionKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserSession", reflect.TypeOf((*MockSessionManger)(nil).DeleteUserSession), ctx, sessionKey)
}
