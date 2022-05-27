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

type TransactionalHashTestSuite struct {
	suite.Suite
	assert *assert.Assertions
	pix    *bankly.Pix
	ctx    context.Context
}

func TestTransactionalHashTestSuite(t *testing.T) {
	suite.Run(t, new(PixTestSuite))
}

func (s *TransactionalHashTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.ctx = context.Background()

	session, err := bankly.NewSession(bankly.Config{
		ClientID:     bankly.String(*bankly.GetEnvBanklyClientID()),
		ClientSecret: bankly.String(*bankly.GetEnvBanklyClientSecret()),
		Scopes:       bankly.String("pix.account.read pix.entries.create pix.entries.delete pix.entries.read pix.qrcode.create pix.qrcode.read pix.cashout.create pix.cashout.read"),
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

func (c *PixTestSuite) TestTransactionalHash_OK() {
	identifier := "07485815024"
	transactional := buildTransactionalHashRquest(bankly.PixPHONE, "+5511946559874", identifier)
	response, err := c.pix.TransactionalHash(c.ctx, transactional, identifier)
	if err != nil {
		c.assert.Error(err)
	}

	c.assert.NotNil(response)
}

func (c *PixTestSuite) TestTransactionalHash_IdentifierNil_NOK() {
	identifier := ""
	transactional := buildTransactionalHashRquest(bankly.PixPHONE, "+0000000000000", identifier)
	response, err := c.pix.TransactionalHash(c.ctx, transactional, identifier)

	c.assert.Error(err)
	c.assert.Nil(response)
}
func (c *PixTestSuite) TestTransactionalHashValidate_OK() {
	identifier := "00000000000"
	buildTransactional := buildTransactionalHashRquest(bankly.PixPHONE, "+0000000000000", identifier)
	transactional, _ := c.pix.TransactionalHash(c.ctx, buildTransactional, identifier)

	response, _ := c.pix.TransactionalHashValidate(c.ctx, *transactional, identifier)
	c.assert.NotNil(response)
}

func buildTransactionalHashRquest(typePix bankly.PixType, valuePix, accountNumber string) bankly.TransactionalHashRequest {
	return bankly.TransactionalHashRequest{
		Context:   "Pix",
		Operation: "RegisterEntry",
		Data: bankly.TransactionalHashData{
			AddressingKey: bankly.TransactionalHashAddressingKey{
				Type:  typePix,
				Value: valuePix,
			},
		},
	}
}