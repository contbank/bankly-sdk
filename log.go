package bankly

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// LoggingRoundTripper ...
type LoggingRoundTripper struct {
	Proxied http.RoundTripper
}

// RoundTrip ...
func (lrt LoggingRoundTripper) RoundTrip(req *http.Request) (res *http.Response, err error) {

	fields := logrus.Fields{
		"request_id": req.Context().Value("Request-Id"),
		"worker_request_id": req.Context().Value("Worker-Request-Id"),
	}

	now := time.Now()

	fields["request"] = request(req)

	logrus.WithFields(fields).Infof("sending request to %v", req.URL)

	res, err = lrt.Proxied.RoundTrip(req)

	elapsed := time.Since(now)

	if err != nil {
		logrus.
			WithError(err).
			WithFields(fields).
			Error("error while receiving response")
		return
	}

	fields["response"] = response(res)
	fields["latency"] = elapsed.Seconds()

	logrus.
		WithFields(fields).
		Print("request completed successfully")

	return
}

func request(request *http.Request) interface{} {
	r := make(map[string]interface{})

	if request.Body != nil {
		bodyCopy := new(bytes.Buffer)
		io.Copy(bodyCopy, request.Body)
		bodyData := bodyCopy.Bytes()

		var body map[string]interface{}
		json.Unmarshal(bodyData, &body)

		r["body"] = body
		request.Body = ioutil.NopCloser(bytes.NewReader(bodyData))
	}

	r["host"] = request.Host
	r["form"] = request.Form
	r["path"] = request.URL.Path
	r["method"] = request.Method
	r["url"] = request.URL.String()
	r["post_form"] = request.PostForm
	r["remote_addr"] = request.RemoteAddr
	r["query_string"] = request.URL.Query()

	return r
}

func response(response *http.Response) interface{} {
	r := make(map[string]interface{})

	bodyCopy := new(bytes.Buffer)
	io.Copy(bodyCopy, response.Body)
	bodyData := bodyCopy.Bytes()

	var body map[string]interface{}
	json.Unmarshal(bodyData, &body)

	r["body"] = body
	r["status"] = response.StatusCode

	response.Body = ioutil.NopCloser(bytes.NewReader(bodyData))

	return r
}
