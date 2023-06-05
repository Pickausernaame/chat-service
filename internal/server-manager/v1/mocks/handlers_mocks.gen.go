// Code generated by MockGen. DO NOT EDIT.
// Source: handlers.go

// Package managerv1mocks is a generated GoMock package.
package managerv1mocks

import (
	context "context"
	reflect "reflect"

	canreceiveproblems "github.com/Pickausernaame/chat-service/internal/usecases/manager/can-receive-problems"
	getassignedproblems "github.com/Pickausernaame/chat-service/internal/usecases/manager/get-assigned-problems"
	setreadyreceiveproblems "github.com/Pickausernaame/chat-service/internal/usecases/manager/set-ready-receive-problems"
	gomock "github.com/golang/mock/gomock"
)

// MockcanReceiveProblemsUseCase is a mock of canReceiveProblemsUseCase interface.
type MockcanReceiveProblemsUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockcanReceiveProblemsUseCaseMockRecorder
}

// MockcanReceiveProblemsUseCaseMockRecorder is the mock recorder for MockcanReceiveProblemsUseCase.
type MockcanReceiveProblemsUseCaseMockRecorder struct {
	mock *MockcanReceiveProblemsUseCase
}

// NewMockcanReceiveProblemsUseCase creates a new mock instance.
func NewMockcanReceiveProblemsUseCase(ctrl *gomock.Controller) *MockcanReceiveProblemsUseCase {
	mock := &MockcanReceiveProblemsUseCase{ctrl: ctrl}
	mock.recorder = &MockcanReceiveProblemsUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockcanReceiveProblemsUseCase) EXPECT() *MockcanReceiveProblemsUseCaseMockRecorder {
	return m.recorder
}

// Handle mocks base method.
func (m *MockcanReceiveProblemsUseCase) Handle(ctx context.Context, req canreceiveproblems.Request) (canreceiveproblems.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handle", ctx, req)
	ret0, _ := ret[0].(canreceiveproblems.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Handle indicates an expected call of Handle.
func (mr *MockcanReceiveProblemsUseCaseMockRecorder) Handle(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockcanReceiveProblemsUseCase)(nil).Handle), ctx, req)
}

// MocksetReadyReceiveProblemsUseCase is a mock of setReadyReceiveProblemsUseCase interface.
type MocksetReadyReceiveProblemsUseCase struct {
	ctrl     *gomock.Controller
	recorder *MocksetReadyReceiveProblemsUseCaseMockRecorder
}

// MocksetReadyReceiveProblemsUseCaseMockRecorder is the mock recorder for MocksetReadyReceiveProblemsUseCase.
type MocksetReadyReceiveProblemsUseCaseMockRecorder struct {
	mock *MocksetReadyReceiveProblemsUseCase
}

// NewMocksetReadyReceiveProblemsUseCase creates a new mock instance.
func NewMocksetReadyReceiveProblemsUseCase(ctrl *gomock.Controller) *MocksetReadyReceiveProblemsUseCase {
	mock := &MocksetReadyReceiveProblemsUseCase{ctrl: ctrl}
	mock.recorder = &MocksetReadyReceiveProblemsUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocksetReadyReceiveProblemsUseCase) EXPECT() *MocksetReadyReceiveProblemsUseCaseMockRecorder {
	return m.recorder
}

// Handle mocks base method.
func (m *MocksetReadyReceiveProblemsUseCase) Handle(ctx context.Context, req setreadyreceiveproblems.Request) (setreadyreceiveproblems.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handle", ctx, req)
	ret0, _ := ret[0].(setreadyreceiveproblems.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Handle indicates an expected call of Handle.
func (mr *MocksetReadyReceiveProblemsUseCaseMockRecorder) Handle(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MocksetReadyReceiveProblemsUseCase)(nil).Handle), ctx, req)
}

// MockgetAssignedProblemsUseCase is a mock of getAssignedProblemsUseCase interface.
type MockgetAssignedProblemsUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockgetAssignedProblemsUseCaseMockRecorder
}

// MockgetAssignedProblemsUseCaseMockRecorder is the mock recorder for MockgetAssignedProblemsUseCase.
type MockgetAssignedProblemsUseCaseMockRecorder struct {
	mock *MockgetAssignedProblemsUseCase
}

// NewMockgetAssignedProblemsUseCase creates a new mock instance.
func NewMockgetAssignedProblemsUseCase(ctrl *gomock.Controller) *MockgetAssignedProblemsUseCase {
	mock := &MockgetAssignedProblemsUseCase{ctrl: ctrl}
	mock.recorder = &MockgetAssignedProblemsUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockgetAssignedProblemsUseCase) EXPECT() *MockgetAssignedProblemsUseCaseMockRecorder {
	return m.recorder
}

// Handle mocks base method.
func (m *MockgetAssignedProblemsUseCase) Handle(ctx context.Context, req getassignedproblems.Request) (getassignedproblems.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handle", ctx, req)
	ret0, _ := ret[0].(getassignedproblems.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Handle indicates an expected call of Handle.
func (mr *MockgetAssignedProblemsUseCaseMockRecorder) Handle(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockgetAssignedProblemsUseCase)(nil).Handle), ctx, req)
}
