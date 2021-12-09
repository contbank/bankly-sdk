package bankly

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/sirupsen/logrus"
)

//Card ...
type Card struct {
	session        Session
	authentication *Authentication
	httpClient     *http.Client
}

//NewCard ...
func NewCard(httpClient *http.Client, session Session) *Card {
	return &Card{
		session:        session,
		httpClient:     httpClient,
		authentication: NewAuthentication(httpClient, session),
	}
}

//Cards ...
func (c *Card) GetCardsByIdentifier(ctx context.Context, identifier string) ([]CardResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
		"identifier": identifier,
	}

	endpoint, _ := c.GetEndpointAPI(fields, CardDocumentPath, identifier)

	resp, _ := c.NewRequest(fields, ctx, "GET", endpoint, nil)
	defer resp.Body.Close()

	return handleResponse(resp, fields)
}

//Utils
func (c *Card) NewRequest(fields logrus.Fields, ctx context.Context, method string, url string, body io.Reader) (*http.Response, error) {

	data, err := json.Marshal(body)
	if err != nil {
		logErrorWithFields(fields, err, "error marshal body request", nil)
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(data))
	if err != nil {
		logErrorWithFields(fields, err, "error new request", nil)
		return nil, err
	}

	token, err := c.authentication.Token(ctx)
	if err != nil {
		logErrorWithFields(fields, err, "error authentication", nil)
		return nil, err
	}

	req = setRequestHeader(req, token, c.session.APIVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logErrorWithFields(fields, err, "error http client", nil)
		return nil, err
	}

	return resp, nil
}

func (c *Card) GetEndpointAPI(fields logrus.Fields, cardPathUrl ...string) (string, error) {
	u, err := url.Parse(c.session.APIEndpoint)
	if err != nil {
		logErrorWithFields(fields, err, "error parsing api endpoint", nil)
		return "", err
	}

	u.Path = path.Join(u.Path, parseUrlValid(cardPathUrl))
	endpoint := u.String()
	fields["endpoint"] = endpoint
	logInfoWithFields(fields, "get endpoint sucess")
	return endpoint, nil
}

func handleResponse(resp *http.Response, fields logrus.Fields) ([]CardResponse, error) {
	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response []CardResponse

		err := json.Unmarshal(respBody, &response)
		if err != nil {
			logErrorWithFields(fields, err, "error decoding json response", nil)
			return nil, ErrDefaultCard
		}
		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var bodyErr *ErrorResponse

	err := json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		logErrorWithFields(fields, err, "error decoding json response", nil)
		return nil, ErrDefaultCard
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		err := FindError(errModel.Code, errModel.Messages...)

		var hasField = make(map[string]interface{})
		hasField["bankly_error"] = bodyErr
		logErrorWithFields(fields, err, "bankly get card error", hasField)

		return nil, err
	}

	return nil, ErrDefaultCard
}

func parseUrlValid(props []string) string {
	var result string = ""
	for _, paths := range props {
		result = result + "/" + paths
	}
	return result
}

func logErrorWithFields(fields logrus.Fields, err error, msg string, hasField map[string]interface{}) {
	if hasField != nil {
		for prop, value := range hasField {
			logrus.
				WithField(prop, value).
				WithFields(fields).
				WithError(err).
				Error(msg)
		}
	} else {
		logrus.
			WithFields(fields).
			WithError(err).
			Error(msg)
	}
}

func logInfoWithFields(fields logrus.Fields, msg string) {
	logrus.
		WithFields(fields).
		Info(msg)
}
