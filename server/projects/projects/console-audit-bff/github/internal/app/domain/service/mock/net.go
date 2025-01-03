// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/domain/service/net.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	echo "github.com/labstack/echo/v4"
	presenter "github.com/pismo/console-audit-bff/internal/app/adapter/web/presenter"
)

// MockNetService is a mock of NetService interface.
type MockNetService struct {
	ctrl     *gomock.Controller
	recorder *MockNetServiceMockRecorder
}

// MockNetServiceMockRecorder is the mock recorder for MockNetService.
type MockNetServiceMockRecorder struct {
	mock *MockNetService
}

// NewMockNetService creates a new mock instance.
func NewMockNetService(ctrl *gomock.Controller) *MockNetService {
	mock := &MockNetService{ctrl: ctrl}
	mock.recorder = &MockNetServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNetService) EXPECT() *MockNetServiceMockRecorder {
	return m.recorder
}

// Whoami mocks base method.
func (m *MockNetService) Whoami(ctx context.Context, netCtx echo.Context) *presenter.IP {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Whoami", ctx, netCtx)
	ret0, _ := ret[0].(*presenter.IP)
	return ret0
}

// Whoami indicates an expected call of Whoami.
func (mr *MockNetServiceMockRecorder) Whoami(ctx, netCtx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Whoami", reflect.TypeOf((*MockNetService)(nil).Whoami), ctx, netCtx)
}

