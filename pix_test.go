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
		Scopes:       bankly.String("pix.account.read pix.entries.create pix.entries.delete pix.entries.read pix.qrcode.create pix.qrcode.read pix.cashout.create pix.cashout.read"),
	})
	s.assert.NoError(err)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	newHttpClient := bankly.NewBanklyHttpClient(*session, httpClient, bankly.NewAuthentication(httpClient, *session))

	s.pix = bankly.NewPix(newHttpClient)
}

// GetAddressKeysByAccount_OK ...
func (c *PixTestSuite) GetAddressKeysByAccount_OK() {
	ctx := context.Background()
	accountNumber := "207802"
	currentIdentity := "35872232000114"
	response, err := c.pix.GetAddressKeysByAccount(ctx, accountNumber, currentIdentity)
	c.assert.NoError(err)
	c.assert.NotNil(response)
}

// GetAddressKeysByAccount_NotFound ...
func (c *PixTestSuite) GetAddressKeysByAccount_NotFound() {
	ctx := context.Background()
	accountNumber := "101010101"
	currentIdentity := "35872232000114"
	response, err := c.pix.GetAddressKeysByAccount(ctx, accountNumber, currentIdentity)
	c.assert.NoError(err)
	c.assert.NotNil(response)
}

// TestGetAddresskey_OK ...
func (c *PixTestSuite) TestGetAddresskey_OK() {
	// TODO Mockar teste
	c.T().Skip("Bankly está retornando erro 504 gateway timeout. Mockar este teste.")

	key := "16246241620"
	currentIdentity := "36183588814"
	response, err := c.pix.GetAddressKey(context.Background(), key, currentIdentity)
	c.assert.NoError(err)
	c.assert.NotNil(response)
}

// TestQrCodeDecode_OK ...
func (c *PixTestSuite) TestQrCodeDecode_OK() {
	// TODO Mockar teste
	c.T().Skip("Bankly está retornando erro 504 gateway timeout. Mockar este teste.")

	currentIdentity := "36183588814"
	response, err := c.pix.QrCodeDecode(context.Background(), &bankly.PixQrCodeDecodeRequest{
		EncodedValue: "MDAwMjAxMjYzMzAwMTRici5nb3YuYmNiLnBpeDAxMTEzNjE4MzU4ODgxNDUyMDQwMDAwNTMwMzk4NjU0MDUxMC4wMDU4MDJCUjU5MTlHdWlsaGVybWUgR29uY2FsdmVzNjAwOVNhbyBQYXVsbzYxMDgwMzUwMzAzMDYyMTQwNTEwMjMxMzEyMzEyMzYzMDQyRjU1",
	}, currentIdentity)

	c.assert.NoError(err)
	c.assert.NotNil(response)
}

func (c *PixTestSuite) TestCreatePixByCPF_OK() {
	c.T().Skip("Está retornando erro 504 Gateway Timeout")

	addressingKeyValue := "41345365373"

	c.pix.DeleteAddressKey(context.Background(), addressingKeyValue, addressingKeyValue)

	pix := builderCreateAddressKeyRequest(bankly.PixCPF, addressingKeyValue, "201928")
	response, err := c.pix.CreateAddressKey(context.Background(), pix)

	c.assert.NoError(err)
	c.assert.NotNil(response)
}

func (c *PixTestSuite) TestCreatePixByCNPJ_OK() {
	// TODO Mockar teste
	c.T().Skip("Corrigir este teste devido ao limite de EVPs que podem ser gerados por account")

	addressKey := "58285483000106"

	c.pix.DeleteAddressKey(context.Background(), addressKey, addressKey)

	pix := builderCreateAddressKeyRequest(bankly.PixCNPJ, addressKey, "201952")

	response, err := c.pix.CreateAddressKey(context.Background(), pix)

	c.assert.NoError(err)
	c.assert.NotNil(response)
}

func (c *PixTestSuite) TestCreatePixByEVP_OK() {
	// TODO Mockar teste
	c.T().Skip("Corrigir este teste devido ao limite de EVPs que podem ser gerados por account")

	pix := builderCreateAddressKeyRequest(bankly.PixEVP, "", "200883")

	response, err := c.pix.CreateAddressKey(context.Background(), pix)

	c.assert.NoError(err)
	c.assert.NotNil(response)
}

func (c *PixTestSuite) TestDeletePixByAddressKey_OK() {
	
	c.T().Skip("Bankly falhando")

	addressingKeyValue := "41345365373"
	indentifier := "41345365373"
	addressKeyRequest := builderCreateAddressKeyRequest(bankly.PixCPF, addressingKeyValue, "201928")
	c.pix.CreateAddressKey(context.Background(), addressKeyRequest)

	time.Sleep(2000)

	err := c.pix.DeleteAddressKey(context.Background(), addressingKeyValue, indentifier)

	c.assert.NoError(err)
}

func builderCreateAddressKeyRequest(typePix bankly.PixType, valuePix, accountNumber string) *bankly.PixAddressKeyCreateRequest {
	return &bankly.PixAddressKeyCreateRequest{
		AddressingKey: bankly.PixTypeValue{
			Type:  typePix,
			Value: valuePix,
		},
		Account: bankly.Account{
			Number:      accountNumber,
			Branch:      "0001",
			AccountType: bankly.CheckingAccount,
		},
	}
}

func (c *PixTestSuite) TestClaimPixByCPF() {
	addressingKeyValue := "16246241620"

	pix := builderClaimRequest(bankly.PixCPF, addressingKeyValue, "207802", bankly.Portability)
	currentIdentity := "16246241620"
	response, err := c.pix.CreatePixClaim(context.Background(), pix, currentIdentity)

	c.assert.NoError(err)
	c.assert.NotNil(response)
}

func builderClaimRequest(typePix bankly.PixType, valuePix string, accountNumber string, claimType bankly.PixClaimType) *bankly.PixClaimRequest {
	return &bankly.PixClaimRequest{
		Type: claimType,
		AddressingKey: bankly.PixTypeValue{
			Type:  typePix,
			Value: valuePix,
		},
		Claimer: bankly.Claimer{
			Branch: "0001",
			Number: accountNumber,
			Bank: bankly.BankClaimer{
				Name: "Acesso Soluções de Pagamento S.A",
				Ispb: "13140088",
			},
		},
		
	}
}