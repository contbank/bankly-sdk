package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/contbank/bankly-sdk/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

var nilHeader *http.Header

func jsonDumps(data interface{}) io.ReadCloser {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return io.NopCloser(bytes.NewReader(dataBytes))
}

func TestWebhook_RegisterWebhook(t *testing.T) {
	testId := "Id"
	banklyHttpClient := mocks.NewBanklyHttpClient(t)
	instance := NewWebhook(banklyHttpClient)
	testWebhook := ConfigItem{
		Name:      "TED_TED_CASH_IN_WAS_RECEIVED",
		Context:   "Ted",
		EventName: "TED_CASH_IN_WAS_RECEIVED",
		Uri:       "http://test/bankly/event",
		PublicKey: "public-key",
	}
	requestData := RegisterWebhookRequest{
		ConfigItem: testWebhook,
		PrivateKey: "private-key",
	}
	expectedResponse := &RegisterWebhookResponse{
		Data:  ConfigEntity{Id: testId, ConfigItem: testWebhook},
		Links: []SchemaLink{},
	}
	banklyHttpClient.EXPECT().Post(mock.Anything, "/webhooks/configurations", requestData, nilHeader).
		Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       jsonDumps(expectedResponse),
		}, nil)
	actualResponse, err := instance.RegisterWebhook(context.Background(), requestData)
	require.Nil(t, err)
	assert.Equal(t, expectedResponse, actualResponse)
}
