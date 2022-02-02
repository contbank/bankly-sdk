package bankly_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	bankly "github.com/contbank/bankly-sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PixTestSuite struct {
	suite.Suite
	assert *assert.Assertions
	pix    *bankly.Pix
	ctx    context.Context
}

func TestPixTestSuite(t *testing.T) {
	suite.Run(t, new(PixTestSuite))
}

func (s *PixTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.ctx = context.Background()

	session, err := bankly.NewSession(bankly.Config{
		ClientID:     bankly.String(*bankly.GetEnvBanklyClientID()),
		ClientSecret: bankly.String(*bankly.GetEnvBanklyClientSecret()),
	})
	s.assert.NoError(err)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	newHttpClient := bankly.BanklyHttpClient{
		Session:        *session,
		HttpClient:     httpClient,
		Authentication: bankly.NewAuthentication(httpClient, *session),
	}

	s.pix = bankly.NewPix(newHttpClient)
}

func (c *PixTestSuite) TestGetAddresskey_OK() {
	key := "16246241620"
	currentIdentity := "36183588814"
	response, err := c.pix.GetAddresskey(context.Background(), key, currentIdentity)
	c.assert.NoError(err)
	c.assert.NotNil(response)
}

func (c *PixTestSuite) TestGetAddresskey_NOT_FOUND() {
	key := "96337603052"
	currentIdentity := "36183588814"
	response, err := c.pix.GetAddresskey(context.Background(), key, currentIdentity)
	c.assert.Error(err)
	c.assert.Nil(response)
	c.assert.EqualError(bankly.ErrEntryNotFound, err.Error())
}

/*
func (c *PixTestSuite) TestGetAddresskey_INVALID_PARAMETER_KEY_TYPE() {
	key := "96337603052"
	currentIdentity := "36183588814"
	response, err := c.pix.GetAddresskey(context.Background(), key, currentIdentity)
	c.assert.Error(err)
	c.assert.Nil(response)
	c.assert.EqualError(bankly.ErrInvalidKeyType, err.Error())
}
*/

func (c *PixTestSuite) TestQrCodeDecode_OK() {
	currentIdentity := "36183588814"

	response, err := c.pix.QrCodeDecode(context.Background(), &bankly.PixQrCodeDecodeRequest{
		EncodedValue: "MDAwMjAxMjYzMzAwMTRici5nb3YuYmNiLnBpeDAxMTEzNjE4MzU4ODgxNDUyMDQwMDAwNTMwMzk4NjU0MDUxMC4wMDU4MDJCUjU5MTlHdWlsaGVybWUgR29uY2FsdmVzNjAwOVNhbyBQYXVsbzYxMDgwMzUwMzAzMDYyMTQwNTEwMjMxMzEyMzEyMzYzMDQyRjU1",
	}, currentIdentity)

	c.assert.NoError(err)
	c.assert.NotNil(response)
}

/*
func (c *PixTestSuite) TestQrCodeDecode_INVALID_QRCODE_PAYLOAD() {
	currentIdentity := "36183588814"

	response, err := c.pix.QrCodeDecode(context.Background(), &bankly.PixQrCodeDecodeRequest{
		EncodedValue: "INVALID",
	}, currentIdentity)

	c.assert.Error(err)
	c.assert.Nil(response)
	c.assert.EqualError(bankly.ErrInvalidQrCodePayload, err.Error())
}

func (c *PixTestSuite) TestCashOut_OK() {
	key := "36183588814"
	currentIdentity := "16246241620"
	keyResponse, err := c.pix.GetAddresskey(context.Background(), key, currentIdentity)

	c.assert.NoError(err)
	c.assert.NotNil(keyResponse)

	response, err := c.pix.CashOut(context.Background(), &bankly.PixCashOutRequest{
		Amount: 0.01,
		Sender: bankly.PixCashOutSenderRequest{
			Account: bankly.PixCashOutAccountRequest{
				Branch: "0001",
				Number: "207802",
			},
			DocumentNumber: currentIdentity,
			Name:           "NOME DA PESSOA OASKDO",
		},
		Description:        "PIX_CASH_OUT_TEST",
		EndToEndID:         keyResponse.EndToEndID,
		InitializationType: bankly.Key,
	})

	c.assert.NoError(err)
	c.assert.NotNil(response)

	key = "16246241620"
	currentIdentity = "36183588814"
	keyResponse, err = c.pix.GetAddresskey(context.Background(), key, currentIdentity)

	c.assert.NoError(err)
	c.assert.NotNil(keyResponse)

	response, err = c.pix.CashOut(context.Background(), &bankly.PixCashOutRequest{
		Amount: 0.01,
		Sender: bankly.PixCashOutSenderRequest{
			Account: bankly.PixCashOutAccountRequest{
				Branch: "0001",
				Number: "184152",
			},
			DocumentNumber: currentIdentity,
			Name:           "NOME DA PESSOA OASKDO",
		},
		Description:        "PIX_CASH_OUT_TEST",
		EndToEndID:         keyResponse.EndToEndID,
		InitializationType: bankly.Key,
	})

	c.assert.NoError(err)
	c.assert.NotNil(response)
}

func (c *PixTestSuite) TestCashOut_INSUFFICIENTBALANCE() {
	key := "36183588814"
	currentIdentity := "16246241620"
	keyResponse, err := c.pix.GetAddresskey(context.Background(), key, currentIdentity)

	c.assert.NoError(err)
	c.assert.NotNil(keyResponse)

	response, err := c.pix.CashOut(context.Background(), &bankly.PixCashOutRequest{
		Amount: 1000000000,
		Sender: bankly.PixCashOutSenderRequest{
			Account: bankly.PixCashOutAccountRequest{
				Branch: "0001",
				Number: "207802",
			},
			DocumentNumber: currentIdentity,
			Name:           "NOME DA PESSOA OASKDO",
		},
		Description:        "PIX_CASH_OUT_TEST",
		EndToEndID:         keyResponse.EndToEndID,
		InitializationType: bankly.Key,
	})

	c.assert.Error(err)
	c.assert.Nil(response)
	c.assert.EqualError(bankly.ErrInsufficientBalance, err.Error())

}
*/
