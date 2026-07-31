package authserver

import (
	"log"
	"net/http"
)

// VerboseRoundTripper wraps an http.RoundTripper and logs the URL of each
// outgoing HTTP request along with the status code of its response.
type VerboseRoundTripper struct {
	Transport http.RoundTripper
}

func (v *VerboseRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	transport := v.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	resp, err := transport.RoundTrip(req)
	if err != nil {
		log.Printf("%s %s -> error: %v", req.Method, req.URL.String(), err)
		return resp, err
	}

	log.Printf("%s %s -> %d", req.Method, req.URL.String(), resp.StatusCode)
	return resp, nil
}
