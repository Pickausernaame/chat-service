// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go

// Package resolveproblemmocks is a generated GoMock package.
package resolveproblemmocks

import (
	context "context"
	reflect "reflect"
	time "time"

	messagesrepo "github.com/Pickausernaame/chat-service/internal/repositories/messages"
	problemsrepo "github.com/Pickausernaame/chat-service/internal/repositories/problems"
	types "github.com/Pickausernaame/chat-service/internal/types"
	gomock "github.com/golang/mock/gomock"
)

// MockmessagesRepository is a mock of messagesRepository interface.
type MockmessagesRepository struct {
	ctrl     *gomock.Controller
	recorder *MockmessagesRepositoryMockRecorder
}

// MockmessagesRepositoryMockRecorder is the mock recorder for MockmessagesRepository.
type MockmessagesRepositoryMockRecorder struct {
	mock *MockmessagesRepository
}

// NewMockmessagesRepository creates a new mock instance.
func NewMockmessagesRepository(ctrl *gomock.Controller) *MockmessagesRepository {
	mock := &MockmessagesRepository{ctrl: ctrl}
	mock.recorder = &MockmessagesRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockmessagesRepository) EXPECT() *MockmessagesRepositoryMockRecorder {
	return m.recorder
}

// CreateProblemResolvedMessage mocks base method.
func (m *MockmessagesRepository) CreateProblemResolvedMessage(ctx context.Context, chatID types.ChatID, problemID types.ProblemID, reqID types.RequestID) (*messagesrepo.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProblemResolvedMessage", ctx, chatID, problemID, reqID)
	ret0, _ := ret[0].(*messagesrepo.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProblemResolvedMessage indicates an expected call of CreateProblemResolvedMessage.
func (mr *MockmessagesRepositoryMockRecorder) CreateProblemResolvedMessage(ctx, chatID, problemID, reqID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProblemResolvedMessage", reflect.TypeOf((*MockmessagesRepository)(nil).CreateProblemResolvedMessage), ctx, chatID, problemID, reqID)
}

// MockproblemsRepository is a mock of problemsRepository interface.
type MockproblemsRepository struct {
	ctrl     *gomock.Controller
	recorder *MockproblemsRepositoryMockRecorder
}

// MockproblemsRepositoryMockRecorder is the mock recorder for MockproblemsRepository.
type MockproblemsRepositoryMockRecorder struct {
	mock *MockproblemsRepository
}

// NewMockproblemsRepository creates a new mock instance.
func NewMockproblemsRepository(ctrl *gomock.Controller) *MockproblemsRepository {
	mock := &MockproblemsRepository{ctrl: ctrl}
	mock.recorder = &MockproblemsRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockproblemsRepository) EXPECT() *MockproblemsRepositoryMockRecorder {
	return m.recorder
}

// GetProblemByChatAndManagerIDs mocks base method.
func (m *MockproblemsRepository) GetProblemByChatAndManagerIDs(ctx context.Context, chatID types.ChatID, managerID types.UserID) (*problemsrepo.Problem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProblemByChatAndManagerIDs", ctx, chatID, managerID)
	ret0, _ := ret[0].(*problemsrepo.Problem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProblemByChatAndManagerIDs indicates an expected call of GetProblemByChatAndManagerIDs.
func (mr *MockproblemsRepositoryMockRecorder) GetProblemByChatAndManagerIDs(ctx, chatID, managerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProblemByChatAndManagerIDs", reflect.TypeOf((*MockproblemsRepository)(nil).GetProblemByChatAndManagerIDs), ctx, chatID, managerID)
}

// ResolveProblem mocks base method.
func (m *MockproblemsRepository) ResolveProblem(ctx context.Context, problemID types.ProblemID, managerID types.UserID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResolveProblem", ctx, problemID, managerID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResolveProblem indicates an expected call of ResolveProblem.
func (mr *MockproblemsRepositoryMockRecorder) ResolveProblem(ctx, problemID, managerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResolveProblem", reflect.TypeOf((*MockproblemsRepository)(nil).ResolveProblem), ctx, problemID, managerID)
}

// Mocktransactor is a mock of transactor interface.
type Mocktransactor struct {
	ctrl     *gomock.Controller
	recorder *MocktransactorMockRecorder
}

// MocktransactorMockRecorder is the mock recorder for Mocktransactor.
type MocktransactorMockRecorder struct {
	mock *Mocktransactor
}

// NewMocktransactor creates a new mock instance.
func NewMocktransactor(ctrl *gomock.Controller) *Mocktransactor {
	mock := &Mocktransactor{ctrl: ctrl}
	mock.recorder = &MocktransactorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mocktransactor) EXPECT() *MocktransactorMockRecorder {
	return m.recorder
}

// RunInTx mocks base method.
func (m *Mocktransactor) RunInTx(ctx context.Context, f func(context.Context) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunInTx", ctx, f)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunInTx indicates an expected call of RunInTx.
func (mr *MocktransactorMockRecorder) RunInTx(ctx, f interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunInTx", reflect.TypeOf((*Mocktransactor)(nil).RunInTx), ctx, f)
}

// MockoutboxService is a mock of outboxService interface.
type MockoutboxService struct {
	ctrl     *gomock.Controller
	recorder *MockoutboxServiceMockRecorder
}

// MockoutboxServiceMockRecorder is the mock recorder for MockoutboxService.
type MockoutboxServiceMockRecorder struct {
	mock *MockoutboxService
}

// NewMockoutboxService creates a new mock instance.
func NewMockoutboxService(ctrl *gomock.Controller) *MockoutboxService {
	mock := &MockoutboxService{ctrl: ctrl}
	mock.recorder = &MockoutboxServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockoutboxService) EXPECT() *MockoutboxServiceMockRecorder {
	return m.recorder
}

// Put mocks base method.
func (m *MockoutboxService) Put(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", ctx, name, payload, availableAt)
	ret0, _ := ret[0].(types.JobID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Put indicates an expected call of Put.
func (mr *MockoutboxServiceMockRecorder) Put(ctx, name, payload, availableAt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockoutboxService)(nil).Put), ctx, name, payload, availableAt)
}