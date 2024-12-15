package handler

import (
	"context"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/pismo/console-audit-bff/internal/app/adapter/web/presenter"
	"github.com/pismo/console-audit-bff/internal/app/domain/service"
	"github.com/pismo/console-audit-bff/internal/app/domain/service/mock"
	srvmock "github.com/pismo/psm-sdk/psm/server/mock"
	"github.com/stretchr/testify/assert"
)

func TestSetManager(t *testing.T) {
	expectedManager := service.NewManager(nil)
	SetManager(expectedManager)
	assert.Equal(t, expectedManager, manager)
}

func TestWhoami(t *testing.T) {
	testCases := []struct {
		name string
		mock func(service *mock.MockNetService, req *srvmock.MockRequest, res *srvmock.MockResponse, echoCtx echo.Context)
		err  error
	}{
		{
			name: "ReturnIPAddress",
			mock: func(service *mock.MockNetService, req *srvmock.MockRequest, res *srvmock.MockResponse, ctxEcho echo.Context) {
				req.EXPECT().EchoCtx().Return(ctxEcho)
				service.EXPECT().
					Whoami(context.Background(), ctxEcho).
					Return(&presenter.IP{IP: "123.123.123.123"})
				res.EXPECT().JSON(http.StatusOK, &presenter.IP{
					IP: "123.123.123.123",
				}).Return(nil)
			},
			err: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockedNetService := mock.NewMockNetService(ctrl)
			mockedRequest := srvmock.NewMockRequest(ctrl)
			mockedResponse := srvmock.NewMockResponse(ctrl)
			echoServer := echo.New()
			echoCtx := echoServer.AcquireContext()
			SetManager(&service.Manager{
				NetService: mockedNetService,
			})

			testCase.mock(mockedNetService, mockedRequest, mockedResponse, echoCtx)
			err := Whoami(context.Background(), mockedRequest, mockedResponse)
			assert.Equal(t, testCase.err, err)
		})
	}
}

