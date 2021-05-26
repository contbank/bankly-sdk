package bankly

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/contbank/grok"
	"github.com/sirupsen/logrus"
)

// DocumentAnalysis ...
type DocumentAnalysis struct {
	session        Session
	httpClient     *http.Client
	authentication *Authentication
}

// NewDocumentAnalysis ...
func NewDocumentAnalysis(session Session) *DocumentAnalysis {
	return &DocumentAnalysis{
		session: session,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		authentication: NewAuthentication(session),
	}
}

// SendDocumentAnalysis ...
func (c *DocumentAnalysis) SendDocumentAnalysis(documentNumber string, request DocumentAnalysisRequest) (*DocumentAnalysisRequestedResponse, error) {

	endpoint, err := c.getDocumentAnalysisAPIEndpoint(documentNumber, nil, nil)
	if err != nil {
		return nil, err
	}

	// TODO Aguardando retorno do Bankly em relação ao envio em base64
	// TODO Enviando Jpeg ou PNG, está retornando "Invalid media type. Use image/png, image/jpg or image/jpeg media type"
	buffer, writer := createMultipartFormData(request,"test_images/selfie1.jpeg")

	req, err := http.NewRequest("PUT", *endpoint, &buffer)
	if err != nil {
		logrus.
			WithError(err).
			Error("error new request")
		return nil, err
	}

	token, err := c.authentication.Token()
	if err != nil {
		logrus.
			WithError(err).
			Error("error authentication")
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("api-version", c.session.APIVersion)
	req.Header.Add("Content-type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logrus.
			WithError(err).
			Error("error http client")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusCreated {
		var bodyResp *DocumentAnalysisRequestedResponse

		err = json.Unmarshal(respBody, &bodyResp)
		if err != nil {
			logrus.
				WithError(err).
				Error("error unmarshal")
			return nil, err
		}

		return bodyResp, nil
	}

	var bodyErr *ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		logrus.
			WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	if bodyErr.Errors != nil {
		return nil, FindError(bodyErr.Errors[0])
	}

	return nil, errors.New("error send document analysis")
}

// FindDocumentAnalysis ...
func (c *DocumentAnalysis) FindDocumentAnalysis(documentNumber string, documentAnalysisToken string) (*DocumentAnalysisResponse, error) {
	resultLevel := ResultLevelDetailed
	endpoint, err := c.getDocumentAnalysisAPIEndpoint(documentNumber, &resultLevel, &documentAnalysisToken)
	if err != nil {
		return nil, err
	} else if endpoint == nil {
		return nil, ErrInvalidAPIEndpoint
	}

	req, err := http.NewRequest("GET", *endpoint, nil)
	if err != nil {
		return nil, err
	}

	token, err := c.authentication.Token()
	if err != nil {
		return nil, err
	}

	req = setRequestHeader(req, token, c.session.APIVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response []*DocumentAnalysisResponse

		err = json.Unmarshal(respBody, &response)
		if err != nil {
			return nil, err
		}

		return response[0], nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var bodyErr *ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		return nil, err
	}

	if bodyErr.Errors != nil {
		return nil, FindError(bodyErr.Errors[0])
	}

	return nil, errors.New("error get document analysis")
}

// getDocumentAnalysisAPIEndpoint
func (c *DocumentAnalysis) getDocumentAnalysisAPIEndpoint(document string, resultLevel *ResultLevel,
	documentAnalysisToken *string) (*string, error) {

	u, err := url.Parse(c.session.APIEndpoint)
	if err != nil {
		logrus.
			WithError(err).
			Error("error api endpoint")
		return nil, err
	}
	u.Path = path.Join(u.Path, DocumentAnalysisPath)
	u.Path = path.Join(u.Path, grok.OnlyDigits(document))
	if documentAnalysisToken != nil {
		q := u.Query()
		q.Set("token", *documentAnalysisToken)
		u.RawQuery = q.Encode()
	}
	if resultLevel != nil {
		q := u.Query()
		q.Set("resultLevel", string(*resultLevel))
		u.RawQuery = q.Encode()
	}
	endpoint := u.String()
	return &endpoint, nil
}

// createMultipartFormData ...
func createMultipartFormData(request DocumentAnalysisRequest, fileName string) (bytes.Buffer, *multipart.Writer) {
	var b bytes.Buffer
	var err error

	w := multipart.NewWriter(&b)
	var writerImage io.Writer

	file := mustOpen(fileName)
	if writerImage, err = w.CreateFormFile("image", file.Name()); err != nil {
		logrus.WithError(err).Errorf("error creating image file writer: %v", err)
	}
	if _, err = io.Copy(writerImage, file); err != nil {
		logrus.WithError(err).Errorf("error with io.Copy: %v", err)
	}

	y1, _ := w.CreateFormField("documentType")
	y1.Write([]byte(request.DocumentType))

	y2, _ := w.CreateFormField("documentSide")
	y2.Write([]byte(request.DocumentSide))

	w.Close()

	return b, w
}

// mustOpen ...
func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		os.Getwd()
		panic(err)
	}
	return r
}