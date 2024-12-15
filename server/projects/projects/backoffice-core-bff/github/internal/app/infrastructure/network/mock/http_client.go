package mock

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/network"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func NewMockHttpClientCustom(t *testing.T, ctrl *gomock.Controller, statusCode int, response interface{}, err error) network.HttpClient {
	var bytesReader *bytes.Reader

	if response == nil {
		bytesReader = bytes.NewReader(make([]byte, 0))
	} else {
		jsonBytes, err := json.Marshal(response)
		if err != nil {
			assert.FailNow(t, err.Error())
		}

		bytesReader = bytes.NewReader(jsonBytes)
	}

	mockedHttpResponse := &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytesReader),
	}

	client := NewMockHttpClient(ctrl)
	client.EXPECT().
		Do(gomock.Any()).
		Return(mockedHttpResponse, err)

	return client
}

