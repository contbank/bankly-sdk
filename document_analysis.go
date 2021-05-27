package bankly

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
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

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open("/Users/firmiano/Downloads/WhatsApp Image 2021-05-24 at 21.39.36 (5).jpeg")
	defer file.Close()

	contentType, _ := getFileContentType(file)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
		escapeQuotes("image"), escapeQuotes("/Users/firmiano/Downloads/WhatsApp Image 2021-05-24 at 21.39.36 (5).jpeg")))
	h.Set("Content-Type", contentType)

	part1, errFile1 := writer.CreatePart(h)
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		fmt.Println(errFile1)
	}
	_ = writer.WriteField("documentType", "RG")
	_ = writer.WriteField("documentSide", "FRONT")
	err = writer.Close()
	if err != nil {
		fmt.Println(err)
	}

	req, err := http.NewRequest("PUT", *endpoint, payload)
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
	req.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logrus.
			WithError(err).
			Error("error http client")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusAccepted {
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

	payload := &bytes.Buffer{}
	w := multipart.NewWriter(payload)
	var writerImage io.Writer

	file, _ := os.Open(fileName)
	defer file.Close()

	if writerImage, err = w.CreateFormFile("image", filepath.Base(fileName)); err != nil {
		logrus.WithError(err).Errorf("error creating image file writer: %v", err)
	}
	if _, err = io.Copy(writerImage, file); err != nil {
		logrus.WithError(err).Errorf("error with io.Copy: %v", err)
	}

	_ = w.WriteField("documentType", "RG")
	_ = w.WriteField("documentSide", "FRONT")
	w.Close()

	return b, w
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func getFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
