package httpProxy

import (
	"fmt"
	"net/http"
)

type LoggingRoundTripper struct {
	Proxied http.RoundTripper
}

func (lrt LoggingRoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {
	// before request
	fmt.Printf("[Logging] Sending request with body: %v\n", req.Body)

	// request sending execution
	resp, e := lrt.Proxied.RoundTrip(req)

	// handle result
	fmt.Printf("[Logging] Receiving responds with body: %v\n", resp.Body)
	return resp, e
}
