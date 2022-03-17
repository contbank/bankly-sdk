package bankly

import (
	"context"
	"encoding/json"
	errors "github.com/contbank/bankly-sdk/pkg/errors"
	models "github.com/contbank/bankly-sdk/pkg/models"
	"github.com/contbank/bankly-sdk/pkg/services/authentication"
	utils "github.com/contbank/bankly-sdk/pkg/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/sirupsen/logrus"
)

// Bank ...
type Bank struct {
	session    bankly.Session
	httpClient *http.Client
	authentication *bankly.Authentication
}

//NewBank ...
func NewBank(httpClient *http.Client, session bankly.Session) *Bank {
	return &Bank{
		session:        session,
		httpClient:     httpClient,
		authentication: bankly.NewAuthentication(httpClient, session),
	}
}

//GetByID returns a list with all available financial instituitions
func (c *Bank) GetByID(ctx context.Context, id string) (*models.BankDataResponse, error) {

	fields := logrus.Fields{
		"bank_id" : id,
		"request_id": utils.GetRequestID(ctx),
	}

	u, err := url.Parse(c.session.APIEndpoint)

	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error parsing api endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, models.BanksPath)
	u.Path = path.Join(u.Path, id)
	endpoint := u.String()

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)

	if err != nil {
		logrus.
			WithError(err).
			WithFields(fields).
			Error("error creating request")
		return nil, err
	}

	token, err := c.authentication.Token(ctx)

	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error in authentication request")
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("api-version", c.session.APIVersion)

	resp, err := c.httpClient.Do(req)

	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error performing the request")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response *models.BankDataResponse

		err = json.Unmarshal(respBody, &response)

		if err != nil {
			logrus.WithFields(fields).
				WithError(err).Error("error decoding json response")
			return nil, err
		}

		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.ErrEntryNotFound
	}

	var bodyErr *errors.ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)

	if err != nil {
		logrus.
			WithError(err).
			WithFields(fields).
			Error("error decoding json response")

		return nil, errors.ErrDefaultBank
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		err = errors.FindError(errModel.Code, errModel.Messages...)
		logrus.
			WithError(err).
			WithFields(fields).
			WithField("bankly_error", bodyErr).
			Error("bankly get bank by id error")
		return nil, err
	}

	return nil, errors.ErrDefaultBank
}

//List returns a list with all available financial instituitions
func (c *Bank) List(ctx context.Context, filter *models.FilterBankListRequest) ([]*models.BankDataResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
	}

	u, err := url.Parse(c.session.APIEndpoint)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error parsing api endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, models.BanksPath)

	q := u.Query()

	for _, id := range filter.IDs {
		q.Add("id", id)
	}

	if filter.Name != nil {
		q.Set("name", *filter.Name)
	}

	if filter.Product != nil {
		q.Set("product", *filter.Product)
	}

	if filter.Page != nil {
		q.Set("page", strconv.Itoa(*filter.Page))
	}

	if filter.Name != nil {
		q.Set("pageSize", strconv.Itoa(*filter.PageSize))
	}

	u.RawQuery = q.Encode()
	endpoint := u.String()

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error creating request")
		return nil, err
	}

	token, err := c.authentication.Token(ctx)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error in authentication request")
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("api-version", c.session.APIVersion)

	resp, err := c.httpClient.Do(req)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error performing the request")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response []*models.BankDataResponse

		err = json.Unmarshal(respBody, &response)

		if err != nil {
			logrus.
				WithFields(fields).
				WithError(err).
				Error("error decoding json response")
			return nil, errors.ErrDefaultBank
		}

		return response, nil
	}

	var bodyErr *errors.ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error decoding json response")
		return nil, errors.ErrDefaultBank
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		err = errors.FindError(errModel.Code, errModel.Messages...)

		logrus.
			WithField("bankly_error", bodyErr).
			WithFields(fields).
			WithError(err).
			Error("bankly filter banks")

		return nil, err
	}

	return nil, errors.ErrDefaultBank
}