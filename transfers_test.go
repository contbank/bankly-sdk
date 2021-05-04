package bankly_test

import (
	"github.com/contbank/bankly-sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type TransfersTestSuite struct {
	suite.Suite
	assert    	*assert.Assertions
	session   	*bankly.Session
	transfers  	*bankly.Transfers
}

func TestTransfersTestSuite(t *testing.T) {
	suite.Run(t, new(TransfersTestSuite))
}

func (s *TransfersTestSuite) SetupTest() {
	s.assert = assert.New(s.T())

	session, err := bankly.NewSession(bankly.Config{
		ClientID:     bankly.String(os.Getenv("BANKLY_CLIENT_ID")),
		ClientSecret: bankly.String(os.Getenv("BANKLY_CLIENT_SECRET")),
	})

	s.assert.NoError(err)

	s.session = session
	s.transfers = bankly.NewTransfers(*s.session)
}

func (s *TransfersTestSuite) TestCreateInternalTransfer() {

	// TODO Verificar o correlation id correto

	correlationID := "correlation_id_" + randStringBytes(10)
	transferRequest := createInternalTransferRequest()

	resp, err := s.transfers.CreateInternalTransfer(correlationID, transferRequest)

	s.assert.NoError(err)
	s.assert.Nil(err)
	s.assert.NotNil(resp)
}

func (s *TransfersTestSuite) TestCreateInternalTransferInvalidCorrelationId() {
	correlationID := "invalid_correlation_id"
	transferRequest := createInternalTransferRequest()

	resp, err := s.transfers.CreateInternalTransfer(correlationID, transferRequest)

	s.assert.Error(err)
	s.assert.Equal(err, bankly.ErrInvalidCorrelationId)
	s.assert.Nil(resp)
}

func (s *TransfersTestSuite) TestCreateExternalTransfer() {

	// TODO Verificar o correlation id correto

	correlationID := "correlation_id_" + randStringBytes(10)
	transferRequest := createExternalTransferRequest()

	resp, err := s.transfers.CreateExternalTransfer(correlationID, transferRequest)

	s.assert.NoError(err)
	s.assert.Nil(err)
	s.assert.NotNil(resp)
}

func (s *TransfersTestSuite) TestCreateExternalTransferInvalidCorrelationId() {
	correlationID := "invalid_correlation_id"
	transferRequest := createExternalTransferRequest()

	resp, err := s.transfers.CreateExternalTransfer(correlationID, transferRequest)

	s.assert.Error(err)
	s.assert.Equal(err, bankly.ErrInvalidCorrelationId)
	s.assert.Nil(resp)
}

func createInternalTransferRequest() bankly.TransfersRequest {
	senderRequest := createSenderRequest("0001", "189162", "82895341000137", "NOME DA EMPRESA 1245312")
	recipientRequest := createRecipientRequest("332", "0001", "189081", "59619372000143", "Nome da Empresa XVlBzgbaiC")
	return bankly.TransfersRequest {
		Amount : 100.00,
		Sender : *senderRequest,
		Recipient : *recipientRequest,
		Description : "Descrição da Transação para uma Conta Interna",
	}
}

func createExternalTransferRequest() bankly.TransfersRequest {
	senderRequest := createSenderRequest("0001", "189162", "82895341000137", "NOME DA EMPRESA 1245312")
	recipientRequest := createRecipientRequest("301", "1000", "13122-1", "11111111111", "Nome Qualquer")
	return bankly.TransfersRequest {
		Amount : 25.00,
		Sender : *senderRequest,
		Recipient : *recipientRequest,
		Description : "Descrição da Transação para uma Conta Externa",
	}
}

func createSenderRequest(branch string, account string, document string, name string) *bankly.SenderRequest {
	return &bankly.SenderRequest {
		Branch : branch,
		Account : account,
		Document : document,
		Name : name,
	}
}

func createRecipientRequest(bankCode string, branch string, account string, document string, name string) *bankly.RecipientRequest {
	return &bankly.RecipientRequest {
		BankCode : bankCode,
		Branch : branch,
		Account : account,
		Document : document,
		Name : name,
		TransfersAccountType : bankly.CheckingAccount,
	}
}