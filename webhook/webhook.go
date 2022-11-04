package webhook

import (
	"context"
	"github.com/contbank/bankly-sdk"
	"github.com/contbank/bankly-sdk/client"
	"github.com/sirupsen/logrus"
	"net/http"
)

var RegisterWebhookConflict = &RegisterWebhookResponse{}

type Webhook interface {
	RegisterWebhook(ctx context.Context, data RegisterWebhookRequest) (out *RegisterWebhookResponse, err error)
}

type webhook struct {
	client client.Client
}

// NewWebhook ...
func NewWebhook(client client.Client) Webhook {
	return &webhook{client: client}
}

func (w webhook) RegisterWebhook(ctx context.Context, data RegisterWebhookRequest) (*RegisterWebhookResponse, error) {
	log := logrus.WithFields(logrus.Fields{
		"request_id": bankly.GetRequestID(ctx),
		"object":     data,
	})
	banklyResponse, err := w.client.Post(ctx, "/webhooks/configurations", data)
	if err != nil {
		log.WithError(err).Error("Error registering the webhook")
		return nil, err
	}
	log = log.WithFields(logrus.Fields{"status": banklyResponse.Status,
		"code": banklyResponse.StatusCode})
	if banklyResponse.StatusCode == http.StatusConflict {
		return RegisterWebhookConflict, nil
	}
	response := &RegisterWebhookResponse{}
	err = banklyResponse.Json(response)
	if err != nil {
		log.WithError(err).Error("Error registering the webhook")
		return nil, err
	}
	return response, nil
}
