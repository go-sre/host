package middleware

import (
	"errors"
	"github.com/go-sre/host/accessdata"
	"net/http"
	"time"
)

type accessWrapper struct {
	rt http.RoundTripper
}

// RoundTrip - implementation of the RoundTrip interface for a transport, also logs an access entry
func (w *accessWrapper) RoundTrip(req *http.Request) (*http.Response, error) {
	var start = time.Now().UTC()

	// !panic
	if w == nil || w.rt == nil {
		return nil, errors.New("invalid handler round tripper configuration : http.RoundTripper is nil")
	}
	resp, err := w.rt.RoundTrip(req)
	if err != nil {
		return resp, err
	}
	entry := accessdata.NewEgressEntry(start, time.Since(start), req, resp, "", -1, -1, -1, "", "", "")
	defaultLogFn(entry)
	return resp, nil
}

// AccessLogWrapTransport - provides a RoundTrip wrapper that applies controller controllers
func AccessLogWrapTransport(client *http.Client) {
	if client == nil || client == http.DefaultClient {
		if http.DefaultClient.Transport == nil {
			http.DefaultClient.Transport = &accessWrapper{http.DefaultTransport}
		} else {
			http.DefaultClient.Transport = AccessWrapRoundTripper(http.DefaultClient.Transport)
		}
	} else {
		if client.Transport == nil {
			client.Transport = &accessWrapper{http.DefaultTransport}
		} else {
			client.Transport = AccessWrapRoundTripper(client.Transport)
		}
	}
}

// AccessWrapRoundTripper - wrap a round tripper
func AccessWrapRoundTripper(rt http.RoundTripper) http.RoundTripper {
	return &accessWrapper{rt}
}
