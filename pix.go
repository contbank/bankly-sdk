package bankly

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Pix struct {
	httpClient BanklyHttpClient
}

//NewPix ...
func NewPix(newHttpClient BanklyHttpClient) *Pix {
	return &Pix{newHttpClient}
}

//Cards ...
func (p *Pix) GetBankByKey(ctx context.Context, key string, currentIdentity string) (*PixBanksByKeyResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
		"key":        key,
	}

	url := "pix/entries/" + key

	header := http.Header{}
	header.Add("x-bkly-pix-user-id", currentIdentity)

	resp, err := p.httpClient.Get(ctx, url, nil, &header)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding body response")
		return nil, err
	}

	response := new(PixBanksByKeyResponse)
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, ErrDefaultCard
	}

	defer resp.Body.Close()
	return response, nil
}
