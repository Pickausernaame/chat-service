// Code generated by MockGen. DO NOT EDIT.
// Source: job.go

// Package sendclientmessagejobmocks is a generated GoMock package.
package sendclientmessagejobmocks

import (
	context "context"
	reflect "reflect"

	messagesrepo "github.com/Pickausernaame/chat-service/internal/repositories/messages"
	eventstream "github.com/Pickausernaame/chat-service/internal/services/event-stream"
	types "github.com/Pickausernaame/chat-service/internal/types"
	gomock "github.com/golang/mock/gomock"
)

// MockmessageRepository is a mock of messageRepository interface.
type MockmessageRepository struct {
	ctrl     *gomock.Controller
	recorder *MockmessageRepositoryMockRecorder
}

// MockmessageRepositoryMockRecorder is the mock recorder for MockmessageRepository.
type MockmessageRepositoryMockRecorder struct {
	mock *MockmessageRepository
}

// NewMockmessageRepository creates a new mock instance.
func NewMockmessageRepository(ctrl *gomock.Controller) *MockmessageRepository {
	mock := &MockmessageRepository{ctrl: ctrl}
	mock.recorder = &MockmessageRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockmessageRepository) EXPECT() *MockmessageRepositoryMockRecorder {
	return m.recorder
}

// GetMessageByID mocks base method.
func (m *MockmessageRepository) GetMessageByID(ctx context.Context, msgID types.MessageID) (*messagesrepo.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessageByID", ctx, msgID)
	ret0, _ := ret[0].(*messagesrepo.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMessageByID indicates an expected call of GetMessageByID.
func (mr *MockmessageRepositoryMockRecorder) GetMessageByID(ctx, msgID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessageByID", reflect.TypeOf((*MockmessageRepository)(nil).GetMessageByID), ctx, msgID)
}

// MockeventStream is a mock of eventStream interface.
type MockeventStream struct {
	ctrl     *gomock.Controller
	recorder *MockeventStreamMockRecorder
}

// MockeventStreamMockRecorder is the mock recorder for MockeventStream.
type MockeventStreamMockRecorder struct {
	mock *MockeventStream
}

// NewMockeventStream creates a new mock instance.
func NewMockeventStream(ctrl *gomock.Controller) *MockeventStream {
	mock := &MockeventStream{ctrl: ctrl}
	mock.recorder = &MockeventStreamMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockeventStream) EXPECT() *MockeventStreamMockRecorder {
	return m.recorder
}

// Publish mocks base method.
func (m *MockeventStream) Publish(ctx context.Context, userID types.UserID, event eventstream.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", ctx, userID, event)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish.
func (mr *MockeventStreamMockRecorder) Publish(ctx, userID, event interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockeventStream)(nil).Publish), ctx, userID, event)
}
