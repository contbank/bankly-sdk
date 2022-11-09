package webhook

import (
	"context"
	"encoding/json"
	"github.com/contbank/bankly-sdk"
	"github.com/sirupsen/logrus"
	"net/http"
)

var RegisterWebhookConflict = &RegisterWebhookResponse{}

type Webhook interface {
	RegisterWebhook(ctx context.Context, data RegisterWebhookRequest) (out *RegisterWebhookResponse, err error)
}

type webhook struct {
	client bankly.BanklyHttpClient
}

// NewWebhook ...
func NewWebhook(client bankly.BanklyHttpClient) Webhook {
	return &webhook{client: client}
}

func (w webhook) RegisterWebhook(ctx context.Context, data RegisterWebhookRequest) (*RegisterWebhookResponse, error) {
	log := logrus.WithFields(logrus.Fields{
		"request_id": bankly.GetRequestID(ctx),
		"object":     data,
	})
	response, err := w.client.Post(ctx, "/webhooks/configurations", data, nil)
	if err != nil {
		log.WithError(err).Error("Error registering the webhook")
		return nil, err
	}
	log = log.WithFields(logrus.Fields{"status": response.Status,
		"code": response.StatusCode})
	if response.StatusCode == http.StatusConflict {
		return RegisterWebhookConflict, nil
	}
	result := &RegisterWebhookResponse{}
	err = json.NewDecoder(response.Body).Decode(result)
	if err != nil {
		log.WithError(err).Error("Error registering the webhook")
		return nil, err
	}
	return result, nil
}
