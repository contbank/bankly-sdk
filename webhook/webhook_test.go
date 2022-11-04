package webhook

import (
	"context"
	"encoding/json"
	"github.com/contbank/bankly-sdk/client"
	mocks "github.com/contbank/bankly-sdk/mocks/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func jsonDumps(data interface{}) []byte {
	dataBytes, _ := json.Marshal(data)
	return dataBytes
}

func TestWebhook_RegisterWebhook(t *testing.T) {
	testId := "Id"
	mockClient := mocks.NewClient(t)
	instance := NewWebhook(mockClient)
	requestData := RegisterWebhookRequest{}
	expectedResponse := &RegisterWebhookResponse{Data: WebhookRecord{Id: testId}, Links: []SchemaLink{}}
	banklyResponse := &client.BanklyResponse{
		Response: &http.Response{
			StatusCode: http.StatusOK,
		},
		Body: jsonDumps(expectedResponse),
	}
	mockClient.EXPECT().Post(mock.Anything, "/webhooks/configurations", requestData).
		Return(banklyResponse, nil)
	actualResponse, err := instance.RegisterWebhook(context.Background(), requestData)
	require.Nil(t, err)
	assert.Equal(t, expectedResponse, actualResponse)
}
