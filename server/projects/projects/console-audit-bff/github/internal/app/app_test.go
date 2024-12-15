//package app
//
//import (
//	"context"
//	"errors"
//	"github.com/pismo/psm-sdk/psm"
//	"github.com/pismo/psm-sdk/psm/factory"
//	"github.com/pismo/psm-sdk/psm/network/http"
//	"github.com/pismo/psm-sdk/psm/network/http/decoder"
//	"github.com/pismo/psm-sdk/psm/server/presenter"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//	"testing"
//	"time"
//)
//
//func TestConfigureServer(t *testing.T) {
//	psm.Configure("console-audit-bff", factory.Default())
//	app, err := NewApplication(psm.HttpClient())
//	require.NoError(t, err)
//	require.NotNil(t, app)
//
//	errChan := make(chan error, 1)
//	go func() {
//		errChan <- app.Start()
//	}()
//
//	select {
//	case err = <-errChan:
//		require.NoError(t, err)
//	case <-time.After(100 * time.Millisecond):
//		// application started
//	}
//
//	defer func() {
//		app.Stop()
//		err = <-errChan
//		assert.Equal(t, errors.New("http: Server closed"), err)
//	}()
//
//	response := psm.HttpClient().Execute(context.Background(), &http.Request{
//		Method: "GET",
//		Host:   "http://localhost:8080",
//		Route:  http.Route("/health"),
//	}, decoder.Json(new(presenter.HealthResponse)))
//
//	assert.NoError(t, response.Error())
//	assert.Equal(t, 200, response.StatusCode())
//	assert.Equal(t, &presenter.HealthResponse{
//		Status: "UP",
//		Indicators: []*presenter.HealthIndicator{
//			{
//				Name:    "http_server",
//				Status:  "UP",
//				Details: make(map[string]interface{}),
//			},
//		},
//	}, response.DecodedBody())
//}
//
