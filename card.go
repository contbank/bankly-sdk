package bankly

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

type Card struct {
	httpClient NewHttpClient
}

//NewCard ...
func NewCard(newHttpClient NewHttpClient) *Card {
	return &Card{newHttpClient}
}

//Cards ...
func (c *Card) GetCardsByIdentifier(ctx context.Context, identifier string) ([]CardResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
		"identifier": identifier,
	}

	url := "cards/document/" + identifier

	resp, err := c.httpClient.Get(ctx, GET, url)
	if err != nil {
		logErrorWithFields(fields, err, err.Error(), nil)
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logErrorWithFields(fields, err, "error decoding body response", nil)
		return nil, err
	}

	var response []CardResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		logErrorWithFields(fields, err, "error decoding json response", nil)
		return nil, ErrDefaultCard
	}

	defer resp.Body.Close()
	return response, nil
}

func (c *Card) CreateCardVirtual(ctx context.Context, body CardCreateRequest) (*CardCreateResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
	}

	url := "cards/virtual"

	resp, err := c.httpClient.Post(ctx, POST, url, body)
	if err != nil {
		logErrorWithFields(fields, err, err.Error(), nil)
		return nil, err
	}

	respBody, _ := ioutil.ReadAll(resp.Body)

	var response *CardCreateResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		logErrorWithFields(fields, err, "error decoding json response", nil)
		return nil, ErrDefaultCard
	}

	defer resp.Body.Close()
	return response, nil
}

func (c *Card) CreateCardPhysical(ctx context.Context, body CardCreateRequest) (*CardCreateResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
	}

	url := "cards/virtual"

	resp, err := c.httpClient.Post(ctx, POST, url, body)
	if err != nil {
		logErrorWithFields(fields, err, err.Error(), nil)
		return nil, err
	}

	respBody, _ := ioutil.ReadAll(resp.Body)

	var response *CardCreateResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		logErrorWithFields(fields, err, "error decoding json response", nil)
		return nil, ErrDefaultCard
	}

	defer resp.Body.Close()
	return response, nil
}

func logErrorWithFields(fields logrus.Fields, err error, msg string, hasField map[string]interface{}) {
	if hasField != nil {
		for prop, value := range hasField {
			logrus.
				WithField(prop, value).
				WithFields(fields).
				WithError(err).
				Error(msg)
		}
	} else {
		logrus.
			WithFields(fields).
			WithError(err).
			Error(msg)
	}
}

func logInfoWithFields(fields logrus.Fields, msg string) {
	logrus.
		WithFields(fields).
		Info(msg)
}
