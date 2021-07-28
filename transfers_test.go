package bankly_test

import (
	"github.com/contbank/bankly-sdk"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"math"
	"net/http"
	"os"
	"testing"
	"time"
)

type TransfersTestSuite struct {
	suite.Suite
	assert    *assert.Assertions
	session   *bankly.Session
	transfers *bankly.Transfers
	balance   *bankly.Balance
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

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	s.session = session
	s.transfers = bankly.NewTransfers(httpClient, *s.session)
	s.balance = bankly.NewBalance(httpClient, *s.session)
}

var internalTransferAmount = []int64{
	075,    // R$ 0,75
	18375,  // R$ 183,75
	1001,   // R$ 10,01
	121374, // R$ 1.213,74
	300,    // R$ 3,00
	1741,   // R$ 17,41
}

func (s *TransfersTestSuite) AfterTest() {
	s.reversalBalanceAmount()
}

func (s *TransfersTestSuite) TestCreateInternalTransfer0() {
	s.createInternalTransferTestLogic(internalTransferAmount[0], *accountA(), *accountB())
}

/*
func (s *TransfersTestSuite) TestCreateInternalTransfer1() {
	s.createInternalTransferTestLogic(internalTransferAmount[1], *accountA(), *accountB())
}

func (s *TransfersTestSuite) TestCreateInternalTransfer2() {
	s.createInternalTransferTestLogic(internalTransferAmount[2], *accountA(), *accountB())
}

func (s *TransfersTestSuite) TestCreateInternalTransfer3() {
	s.createInternalTransferTestLogic(internalTransferAmount[3], *accountA(), *accountB())
}

func (s *TransfersTestSuite) TestCreateInternalTransfer4() {
	s.createInternalTransferTestLogic(internalTransferAmount[4], *accountA(), *accountB())
}

func (s *TransfersTestSuite) TestCreateInternalTransfer5() {
	s.createInternalTransferTestLogic(internalTransferAmount[5], *accountA(), *accountB())
}

func (s *TransfersTestSuite) TestCreateInternalTransferAllAvailableAmountBalance1() {
	s.transferAllAvailableAmountBalance(*accountC(), *accountD())
}

func (s *TransfersTestSuite) TestCreateInternalTransferAllAvailableAmountBalance2() {
	s.transferAllAvailableAmountBalance(*accountD(), *accountC())
}

func (s *TransfersTestSuite) TestCreateExternalTransfer() {
	correlationID := uuid.New().String()

	sender := *accountA()
	transferRequest := createExternalTransferRequest(50.00, sender)

	resp, err := s.transfers.CreateExternalTransfer(correlationID, transferRequest)

	if isOutOfServicePeriod() {
		s.assert.Nil(resp)
		s.assert.NotNil(err)
		s.assert.Error(err, bankly.ErrOutOfServicePeriod)
	} else {
		s.assert.NotNil(resp)
		s.assert.NoError(err)
		s.assert.Nil(err)
	}
}

func (s *TransfersTestSuite) TestCreateExternalTransferInvalidAccountNumberFormat() {
	correlationID := uuid.New().String()

	invalidAccountNumber := "13122-1" // hifen

	sender := *accountA()
	transferRequest := bankly.TransfersRequest{
		Amount:      int64(500),
		Sender:      *createSenderRequest(sender.Branch, sender.Account, sender.Document, sender.Name),
		Recipient:   *createRecipientRequest("301", "1000", invalidAccountNumber, "11111111111", "Nome Qualquer"),
		Description: "Descrição da Transação para uma Conta Externa",
	}

	resp, err := s.transfers.CreateExternalTransfer(correlationID, transferRequest)

	s.assert.Error(err)
	s.assert.Contains(err.Error(), "Invalid account number format")
	s.assert.Nil(resp)
}

func (s *TransfersTestSuite) TestCreateExternalTransferInvalidCorrelationId() {
	correlationID := "invalid_correlation_id"
	transferRequest := createExternalTransferRequest(50.00, *accountA())

	resp, err := s.transfers.CreateExternalTransfer(correlationID, transferRequest)

	s.assert.Error(err)
	s.assert.Equal(err, bankly.ErrInvalidCorrelationId)
	s.assert.Nil(resp)
}
*/

func (s *TransfersTestSuite) TestFindTransferByCode1() {
	correlationID := uuid.New().String()
	authenticationCode := "4ebdaae8-663f-4c78-835d-112177f139e8"
	branch := "0001"
	account := "189081"
	receipt, err := s.transfers.FindTransfersByCode(&correlationID, &authenticationCode, &branch, &account)

	s.assert.NoError(err)
	s.assert.NotNil(receipt)
	s.assert.Equal(branch, receipt.Sender.Account.Branch)
	s.assert.Equal(account, receipt.Sender.Account.Number)
	s.assert.Equal(authenticationCode, receipt.AuthenticationCode)
}

func (s *TransfersTestSuite) TestFindTransferByCode2() {
	correlationID := uuid.New().String()
	authenticationCode := "8330d81e-0501-4c7d-8b31-5961fa6f8550"
	branch := "0001"
	account := "189162"
	receipt, err := s.transfers.FindTransfersByCode(&correlationID, &authenticationCode, &branch, &account)

	s.assert.NoError(err)
	s.assert.NotNil(receipt)
	s.assert.Equal(branch, receipt.Sender.Account.Branch)
	s.assert.Equal(account, receipt.Sender.Account.Number)
	s.assert.Equal(authenticationCode, receipt.AuthenticationCode)
}

func (s *TransfersTestSuite) TestFindTransfers() {
	correlationID := uuid.New().String()
	branch := "0001"
	account := "189081"
	pageSize := 10
	transfers, err := s.transfers.FindTransfers(&correlationID, &branch, &account, &pageSize, nil)

	s.assert.NoError(err)
	s.assert.NotNil(transfers)
}

func createInternalTransferRequest(amount int64, from AccountToTest, to AccountToTest) bankly.TransfersRequest {
	senderRequest := createSenderRequest(from.Branch, from.Account, from.Document, from.Name)
	recipientRequest := createRecipientRequest(to.BankCode, to.Branch, to.Account, to.Document, to.Name)
	return bankly.TransfersRequest{
		Amount:      amount,
		Sender:      *senderRequest,
		Recipient:   *recipientRequest,
		Description: "Descrição da Transação para uma Conta Interna",
	}
}

/*
func createExternalTransferRequest(amount int64, from AccountToTest) bankly.TransfersRequest {
	senderRequest := createSenderRequest(from.Branch, from.Account, from.Document, from.Name)
	recipientRequest := createRecipientRequest("301", "1000", "131221", "11111111111", "Nome Qualquer")
	return bankly.TransfersRequest{
		Amount:      amount,
		Sender:      *senderRequest,
		Recipient:   *recipientRequest,
		Description: "Descrição da Transação para uma Conta Externa",
	}
}
*/

func (s *TransfersTestSuite) createInternalTransferTestLogic(amount int64, sender AccountToTest, recipient AccountToTest) {
	correlationID := uuid.New().String()

	senderBalance, _ := s.balance.Balance(sender.Account)
	expectedSenderAvailableAmount := senderBalance.Balance.Available.Amount - (float64(amount) / 100)

	recipientBalance, _ := s.balance.Balance(recipient.Account)
	expectedRecipientAvailableAmount := recipientBalance.Balance.Available.Amount + (float64(amount) / 100)

	transferRequest := createInternalTransferRequest(amount, sender, recipient)

	resp, err := s.transfers.CreateInternalTransfer(correlationID, transferRequest)

	time.Sleep(time.Second * 3)

	senderBalance, _ = s.balance.Balance(sender.Account)
	recipientBalance, _ = s.balance.Balance(recipient.Account)

	s.assert.NoError(err)
	s.assert.NotNil(resp)
	s.assert.NotNil(resp.AuthenticationCode)
	s.assert.Equal(toDecimal(expectedSenderAvailableAmount), toDecimal(senderBalance.Balance.Available.Amount))
	s.assert.Equal(toDecimal(expectedRecipientAvailableAmount), toDecimal(recipientBalance.Balance.Available.Amount))
}

func (s *TransfersTestSuite) transferAllAvailableAmountBalance(from AccountToTest, to AccountToTest) {
	correlationID := uuid.New().String()

	senderBalance, _ := s.balance.Balance(from.Account)
	recipientBalance, _ := s.balance.Balance(to.Account)
	expectedRecipientAvailableAmount :=
		float64(recipientBalance.Balance.Available.Amount) + float64(senderBalance.Balance.Available.Amount)

	amount := int64(senderBalance.Balance.Available.Amount * 100)
	transferRequest := createInternalTransferRequest(amount, from, to)

	resp, err := s.transfers.CreateInternalTransfer(correlationID, transferRequest)

	time.Sleep(time.Second)

	senderBalance, _ = s.balance.Balance(from.Account)
	afterSenderAvailableAmount := float64(senderBalance.Balance.Available.Amount)

	recipientBalance, _ = s.balance.Balance(to.Account)
	afterRecipientAvailableAmount := float64(recipientBalance.Balance.Available.Amount)

	if amount != 0 {
		s.assert.NoError(err)
		s.assert.NotNil(resp)
		s.assert.NotNil(resp.AuthenticationCode)
	}

	s.assert.Equal(float64(0), afterSenderAvailableAmount)
	s.assert.Equal(expectedRecipientAvailableAmount, afterRecipientAvailableAmount)
}

func createSenderRequest(branch string, account string, document string, name string) *bankly.SenderRequest {
	return &bankly.SenderRequest{
		Branch:   branch,
		Account:  account,
		Document: document,
		Name:     name,
	}
}

func createRecipientRequest(bankCode string, branch string, account string, document string, name string) *bankly.RecipientRequest {
	return &bankly.RecipientRequest{
		BankCode:             bankCode,
		Branch:               branch,
		Account:              account,
		Document:             document,
		Name:                 name,
		TransfersAccountType: bankly.CheckingAccount,
	}
}

func (s *TransfersTestSuite) reversalBalanceAmount() {
	// sum of values
	total := int64(0)
	for _, value := range internalTransferAmount {
		total += value
	}
	// revert values
	s.createInternalTransferTestLogic(total, *accountB(), *accountA())
}

func toDecimal(value float64) float64 {
	return math.Round(value*100) / 100
}

/*
// TODO Não está levando em consideração dia útil ou feriados. Verificar melhor alternativa para implementar isto.
func isOutOfServicePeriod() bool {
	currentTime := time.Now()
	initialTime := time.Duration(time.Hour * 7).Minutes()
	limitTime := time.Duration(time.Hour * 17).Minutes()
	now := float64((currentTime.Hour() * 60) + currentTime.Minute())
	if now < initialTime || now > limitTime {
		return true
	}
	return false
}
*/
