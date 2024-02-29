package bankly_test

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"testing"
	"time"

	"github.com/contbank/bankly-sdk"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type IncomeReportTestSuite struct {
	suite.Suite
	assert           *assert.Assertions
	ctx              context.Context
	session          *bankly.Session
	bankIncomeReport *bankly.IncomeReport
}

func TestIncomeReportTestSuite(t *testing.T) {
	suite.Run(t, new(IncomeReportTestSuite))
}

func (s *IncomeReportTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.ctx = context.Background()

	session, err := bankly.NewSession(bankly.Config{
		ClientID:     bankly.String(*bankly.GetEnvBanklyClientID()),
		ClientSecret: bankly.String(*bankly.GetEnvBanklyClientSecret()),
		Scopes:       bankly.String("income.report.read"),
	})

	s.assert.NoError(err)

	httpClient := &http.Client{
		Timeout:   30 * time.Second,
		Transport: bankly.LoggingRoundTripper{Proxied: http.DefaultTransport},
	}

	s.session = session
	s.bankIncomeReport = bankly.NewIncomeReport(httpClient, *s.session)
}

// TestIncomeReport_SUCCESS ...
func (s *IncomeReportTestSuite) TestIncomeReport_SUCCESS() {
	s.T().Skip("token")

	s.ctx = context.WithValue(s.ctx, "Request-Id", primitive.NewObjectID().Hex())
	s.ctx = context.WithValue(s.ctx, "access_token", primitive.NewObjectID().Hex())

	report, err := s.bankIncomeReport.GetIncomeReport(s.ctx,
		&bankly.IncomeReportRequest{
			Account: "184152",
			Year:    "2023",
		})

	s.assert.NoError(err)
	s.assert.NotNil(report)
	s.assert.NotEmpty(report.Data)
	s.assert.NotEmpty(report.Links)
}

// TestIncomeReport_INVALID_ACCOUNT_NUMBER ...
func (s *IncomeReportTestSuite) TestIncomeReport_INVALID_ACCOUNT_NUMBER() {
	s.T().Skip("token")

	report, err := s.bankIncomeReport.GetIncomeReport(s.ctx,
		&bankly.IncomeReportRequest{
			Account: "100000",
			Year:    "2023",
		})

	s.assert.Error(err)
	s.assert.Nil(report)
}
