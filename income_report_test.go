package bankly_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/contbank/bankly-sdk"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type IncomeReportTestSuite struct {
	suite.Suite
	assert  		 *assert.Assertions
	ctx     		 context.Context
	session 		 *bankly.Session
	bankIncomeReport *bankly.IncomeReport
}

func TestIncomeReportTestSuite(t *testing.T) {
	suite.Run(t, new(IncomeReportTestSuite))
}

func (s *IncomeReportTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.ctx = context.Background()

	session, err := bankly.NewSession(bankly.Config{
		ClientID : bankly.String(*bankly.GetEnvBanklyClientID()),
		ClientSecret : bankly.String(*bankly.GetEnvBanklyClientSecret()),
	})

	s.assert.NoError(err)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	s.session = session

	newHttpClient := bankly.BanklyHttpClient{
		Session:        *session,
		HttpClient:     httpClient,
		Authentication: bankly.NewAuthentication(httpClient, *session),
	}

	s.bankIncomeReport = bankly.NewIncomeReport(newHttpClient)
}

// TestIncomeReport_SUCCESS ...
func (s *IncomeReportTestSuite) TestIncomeReport_SUCCESS() {
	report, err := s.bankIncomeReport.GetIncomeReport(s.ctx,
		&bankly.IncomeReportRequest{
			Account: "184152",
			Year: "2021",
		})

	s.assert.NoError(err)
	s.assert.NotNil(report)
	s.assert.NotEmpty(report.FileName)
	s.assert.NotEmpty(report.IncomeFile)
}

// TestIncomeReport_INVALID_ACCOUNT_NUMBER ...
func (s *IncomeReportTestSuite) TestIncomeReport_INVALID_ACCOUNT_NUMBER() {
	report, err := s.bankIncomeReport.GetIncomeReport(s.ctx,
		&bankly.IncomeReportRequest{
			Account: "100000",
			Year: "2021",
		})

	s.assert.Error(err)
	s.assert.Nil(report)
}