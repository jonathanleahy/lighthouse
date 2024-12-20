// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/domain/service/interfaces.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/pismo/backoffice-core-bff/internal/app/domain/entity"
)

// MockAuditService is a mock of AuditService interface.
type MockAuditService struct {
	ctrl     *gomock.Controller
	recorder *MockAuditServiceMockRecorder
}

// MockAuditServiceMockRecorder is the mock recorder for MockAuditService.
type MockAuditServiceMockRecorder struct {
	mock *MockAuditService
}

// NewMockAuditService creates a new mock instance.
func NewMockAuditService(ctrl *gomock.Controller) *MockAuditService {
	mock := &MockAuditService{ctrl: ctrl}
	mock.recorder = &MockAuditServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuditService) EXPECT() *MockAuditServiceMockRecorder {
	return m.recorder
}

// SendAuditEvent mocks base method.
func (m *MockAuditService) SendAuditEvent(ctx context.Context, domain, action, domainID, request, response string, responseCode int) chan error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendAuditEvent", ctx, domain, action, domainID, request, response, responseCode)
	ret0, _ := ret[0].(chan error)
	return ret0
}

// SendAuditEvent indicates an expected call of SendAuditEvent.
func (mr *MockAuditServiceMockRecorder) SendAuditEvent(ctx, domain, action, domainID, request, response, responseCode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendAuditEvent", reflect.TypeOf((*MockAuditService)(nil).SendAuditEvent), ctx, domain, action, domainID, request, response, responseCode)
}

// WaitFinish mocks base method.
func (m *MockAuditService) WaitFinish() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WaitFinish")
}

// WaitFinish indicates an expected call of WaitFinish.
func (mr *MockAuditServiceMockRecorder) WaitFinish() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitFinish", reflect.TypeOf((*MockAuditService)(nil).WaitFinish))
}

// MockDisputeService is a mock of DisputeService interface.
type MockDisputeService struct {
	ctrl     *gomock.Controller
	recorder *MockDisputeServiceMockRecorder
}

// MockDisputeServiceMockRecorder is the mock recorder for MockDisputeService.
type MockDisputeServiceMockRecorder struct {
	mock *MockDisputeService
}

// NewMockDisputeService creates a new mock instance.
func NewMockDisputeService(ctrl *gomock.Controller) *MockDisputeService {
	mock := &MockDisputeService{ctrl: ctrl}
	mock.recorder = &MockDisputeServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDisputeService) EXPECT() *MockDisputeServiceMockRecorder {
	return m.recorder
}

// CreateFraudReport mocks base method.
func (m *MockDisputeService) CreateFraudReport(ctx context.Context, disputeId int, input entity.FraudReportInput) (*entity.FraudReport, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFraudReport", ctx, disputeId, input)
	ret0, _ := ret[0].(*entity.FraudReport)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateFraudReport indicates an expected call of CreateFraudReport.
func (mr *MockDisputeServiceMockRecorder) CreateFraudReport(ctx, disputeId, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFraudReport", reflect.TypeOf((*MockDisputeService)(nil).CreateFraudReport), ctx, disputeId, input)
}

// GetDisputeStatus mocks base method.
func (m *MockDisputeService) GetDisputeStatus(ctx context.Context, disputeId int, disputeInstallmentId *int) ([]*entity.DisputeStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDisputeStatus", ctx, disputeId, disputeInstallmentId)
	ret0, _ := ret[0].([]*entity.DisputeStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDisputeStatus indicates an expected call of GetDisputeStatus.
func (mr *MockDisputeServiceMockRecorder) GetDisputeStatus(ctx, disputeId, disputeInstallmentId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDisputeStatus", reflect.TypeOf((*MockDisputeService)(nil).GetDisputeStatus), ctx, disputeId, disputeInstallmentId)
}

// UpdateDisputeStatusEvent mocks base method.
func (m *MockDisputeService) UpdateDisputeStatusEvent(ctx context.Context, disputeId int, input entity.DisputeEventInput) (*entity.DisputeEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateDisputeStatusEvent", ctx, disputeId, input)
	ret0, _ := ret[0].(*entity.DisputeEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateDisputeStatusEvent indicates an expected call of UpdateDisputeStatusEvent.
func (mr *MockDisputeServiceMockRecorder) UpdateDisputeStatusEvent(ctx, disputeId, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDisputeStatusEvent", reflect.TypeOf((*MockDisputeService)(nil).UpdateDisputeStatusEvent), ctx, disputeId, input)
}

// MockeventSender is a mock of eventSender interface.
type MockeventSender struct {
	ctrl     *gomock.Controller
	recorder *MockeventSenderMockRecorder
}

// MockeventSenderMockRecorder is the mock recorder for MockeventSender.
type MockeventSenderMockRecorder struct {
	mock *MockeventSender
}

// NewMockeventSender creates a new mock instance.
func NewMockeventSender(ctrl *gomock.Controller) *MockeventSender {
	mock := &MockeventSender{ctrl: ctrl}
	mock.recorder = &MockeventSenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockeventSender) EXPECT() *MockeventSenderMockRecorder {
	return m.recorder
}

// SendEvent mocks base method.
func (m *MockeventSender) SendEvent(ctx context.Context, domain, eventType string, body interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendEvent", ctx, domain, eventType, body)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendEvent indicates an expected call of SendEvent.
func (mr *MockeventSenderMockRecorder) SendEvent(ctx, domain, eventType, body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendEvent", reflect.TypeOf((*MockeventSender)(nil).SendEvent), ctx, domain, eventType, body)
}

// MockHealthService is a mock of HealthService interface.
type MockHealthService struct {
	ctrl     *gomock.Controller
	recorder *MockHealthServiceMockRecorder
}

// MockHealthServiceMockRecorder is the mock recorder for MockHealthService.
type MockHealthServiceMockRecorder struct {
	mock *MockHealthService
}

// NewMockHealthService creates a new mock instance.
func NewMockHealthService(ctrl *gomock.Controller) *MockHealthService {
	mock := &MockHealthService{ctrl: ctrl}
	mock.recorder = &MockHealthServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHealthService) EXPECT() *MockHealthServiceMockRecorder {
	return m.recorder
}

// GetMessage mocks base method.
func (m *MockHealthService) GetMessage() (*entity.Health, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessage")
	ret0, _ := ret[0].(*entity.Health)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMessage indicates an expected call of GetMessage.
func (mr *MockHealthServiceMockRecorder) GetMessage() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessage", reflect.TypeOf((*MockHealthService)(nil).GetMessage))
}

