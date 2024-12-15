package service

import (
	"context"
	"errors"
	"testing"

	"net/http"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/pismo/console-audit-bff/internal/app/adapter/web/presenter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoleService_Whoami(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
		res        presenter.IP
		err        error
	}{
		{
			name:       "Success",
			statusCode: http.StatusOK,
			res:        presenter.IP{IP: ""},
		},
		{
			name:       "Error",
			statusCode: http.StatusInternalServerError,
			err:        errors.New("connection_err"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			srv := NewNetService()
			require.NotNil(t, srv)

			ctx := context.TODO()
			echoServer := echo.New()

			echoCtx := echoServer.NewContext(&http.Request{RemoteAddr: "192.168.0.1"}, nil)

			// executes the operati on to be tested
			response := srv.Whoami(ctx, echoCtx)
			if testCase.err == nil {
				assert.Equal(t, response, &testCase.res)
			} else {
				assert.Equal(t, testCase.err, testCase.err)
				assert.Equal(t, response, &testCase.res)
			}
		})
	}
}

