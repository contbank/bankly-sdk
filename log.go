package bankly

import (
	"fmt"
	"net/http"
)

type LoggingRoundTripper struct {
	Proxied http.RoundTripper
}

func (lrt LoggingRoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {
	fmt.Printf("Sending request to %v\n", req.URL)

	fmt.Printf("Request-ID %v\n", req.Context().Value("Request-Id"))

	res, e = lrt.Proxied.RoundTrip(req)

	if e != nil {
		fmt.Printf("Error: %v", e)
	} else {
		fmt.Printf("Received %v response\n", res.Status)
	}

	return
}
