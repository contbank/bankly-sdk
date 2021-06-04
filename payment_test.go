package bankly_test

import (
	"testing"

	"github.com/contbank/bankly-sdk"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PaymentTestSuite struct {
	suite.Suite
	assert  *assert.Assertions
	session *bankly.Session
	payment *bankly.Payment
}

func TestPaymentTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentTestSuite))
}

func (s *PaymentTestSuite) SetupTest() {
	s.assert = assert.New(s.T())

	session, err := bankly.NewSession(bankly.Config{
		ClientID:     bankly.String("0a9f5c95-4b73-44b5-b3dd-03569c570630"),
		ClientSecret: bankly.String("V$^YzR$sI#Qhh4b!e0cHu6B1*r#*vkVj"),
	})

	s.assert.NoError(err)

	s.session = session
	s.payment = bankly.NewPayment(*s.session)
}

func (s *PaymentTestSuite) TestValidatePayment() {
	r, err := s.payment.ValidatePayment(&bankly.ValidatePaymentRequest{
		Code: "23793380296096975582095006333306584750000002400",
	})

	s.assert.NoError(err)
	s.assert.NotNil(r)
}

func (s *PaymentTestSuite) TestConfirmPayment() {
	r, err := s.payment.ValidatePayment(&bankly.ValidatePaymentRequest{
		Code: "23793380296096975582095006333306584750000002400",
	})

	s.assert.NoError(err)
	s.assert.NotNil(r)

	description := "test payment"
	r2, err := s.payment.ConfirmPayment(&bankly.ConfirmPaymentRequest{
		ID:          r.ID,
		Amount:      r.Amount,
		Description: &description,
		BankBranch:  "0001",
		BankAccount: "184152",
	})

	s.assert.NoError(err)
	s.assert.NotNil(r2)
}
