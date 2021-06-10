package bankly_test

import (
	"os"
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
		ClientID:     bankly.String(os.Getenv("BANKLY_CLIENT_ID")),
		ClientSecret: bankly.String(os.Getenv("BANKLY_CLIENT_SECRET")),
	})

	s.assert.NoError(err)

	s.session = session
	s.payment = bankly.NewPayment(*s.session)
}

// func (s *PaymentTestSuite) TestValidatePayment() {
// 	r, err := s.payment.ValidatePayment("cd38008a-dc96-4872-a1a1-96a35a153efd", &bankly.ValidatePaymentRequest{
// 		Code: "836200000021355100403185234319172032100181841691",
// 	})

// 	s.assert.NoError(err)
// 	s.assert.NotNil(r)
// }

// func (s *PaymentTestSuite) TestConfirmPayment() {
// 	r, err := s.payment.ValidatePayment("cd38008a-dc96-4872-a1a1-96a35a153efd", &bankly.ValidatePaymentRequest{
// 		Code: "836200000021355100403185234319172032100181841691",
// 	})

// 	s.assert.NoError(err)
// 	s.assert.NotNil(r)

// 	description := "test payment"
// 	r2, err := s.payment.ConfirmPayment("cd38008a-dc96-4872-a1a1-96a35a153efd", &bankly.ConfirmPaymentRequest{
// 		ID:          r.ID,
// 		Amount:      r.Amount,
// 		Description: &description,
// 		BankBranch:  "0001",
// 		BankAccount: "184152",
// 	})

// 	s.assert.NoError(err)
// 	s.assert.NotNil(r2)
// }

// func (s *PaymentTestSuite) TestFilterPayments() {
// 	r, err := s.payment.FilterPayments("e0f5ff37-a75c-4e57-96e6-b17d7003a0e9", &bankly.FilterPaymentsRequest{
// 		BankBranch:  "0001",
// 		BankAccount: "184152",
// 		PageSize:    10,
// 	})

// 	s.assert.NoError(err)
// 	s.assert.NotNil(r)
// }

// func (s *PaymentTestSuite) TestDetailPayment() {
// 	r, err := s.payment.DetailPayment("e0f5ff37-a75c-4e57-96e6-b17d7003a0e9", &bankly.DetailPaymentRequest{
// 		BankBranch:         "0001",
// 		BankAccount:        "184152",
// 		AuthenticationCode: "710dd95c-7aa6-49f2-94d0-79c93a945528",
// 	})

// 	s.assert.NoError(err)
// 	s.assert.NotNil(r)
// }
