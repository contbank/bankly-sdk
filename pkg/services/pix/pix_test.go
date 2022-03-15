package pix_test

import (
	"context"
	"github.com/contbank/bankly-sdk/pkg/models"
	"github.com/contbank/bankly-sdk/pkg/services/authentication"
	pix2 "github.com/contbank/bankly-sdk/pkg/services/pix"
	"github.com/contbank/bankly-sdk/pkg/utils"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PixTestSuite struct {
	suite.Suite
	assert *assert.Assertions
	pix    *pix2.Pix
	ctx    context.Context
}

func TestPixTestSuite(t *testing.T) {
	suite.Run(t, new(PixTestSuite))
}

func (s *PixTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.ctx = context.Background()

	session, err := authentication.NewSession(authentication.Config{
		ClientID:     utils.String(*utils.GetEnvBanklyClientID()),
		ClientSecret: utils.String(*utils.GetEnvBanklyClientSecret()),
	})
	s.assert.NoError(err)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	newHttpClient := utils.BanklyHttpClient{
		Session:        *session,
		HttpClient:     httpClient,
		Authentication: authentication.NewAuthentication(httpClient, *session),
	}

	s.pix = pix2.NewPix(newHttpClient)
}

func (c *PixTestSuite) TestGetAddresskey_OK() {
	key := "16246241620"
	currentIdentity := "36183588814"
	response, err := c.pix.GetAddresskey(context.Background(), key, currentIdentity)
	c.assert.NoError(err)
	c.assert.NotNil(response)
}

func (c *PixTestSuite) TestQrCodeDecode_OK() {
	currentIdentity := "36183588814"

	response, err := c.pix.QrCodeDecode(context.Background(), &models.PixQrCodeDecodeRequest{
		EncodedValue: "MDAwMjAxMjYzMzAwMTRici5nb3YuYmNiLnBpeDAxMTEzNjE4MzU4ODgxNDUyMDQwMDAwNTMwMzk4NjU0MDUxMC4wMDU4MDJCUjU5MTlHdWlsaGVybWUgR29uY2FsdmVzNjAwOVNhbyBQYXVsbzYxMDgwMzUwMzAzMDYyMTQwNTEwMjMxMzEyMzEyMzYzMDQyRjU1",
	}, currentIdentity)

	c.assert.NoError(err)
	c.assert.NotNil(response)
}
