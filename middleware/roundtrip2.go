package middleware

import (
	"context"
	"errors"
	"github.com/go-sre/host/controller"
	"net/http"
	"time"
)

type controllerWrapper struct {
	rt http.RoundTripper
}

// RoundTrip - implementation of the RoundTrip interface for a transport, also logs an access entry
func (w *controllerWrapper) RoundTrip(req *http.Request) (*http.Response, error) {
	var start = time.Now().UTC()
	var retry = false

	// !panic
	if w == nil || w.rt == nil {
		return nil, errors.New("invalid handler round tripper configuration : http.RoundTripper is nil")
	}
	ctrl := controller.EgressTable().LookupHttp(req)
	ctrl.UpdateHeaders(req)
	if rlc := ctrl.RateLimiter(); rlc.IsEnabled() && !rlc.Allow() {
		resp := &http.Response{Request: req, StatusCode: rlc.StatusCode()}
		ctrl.LogHttpEgress(start, time.Since(start), req, resp, controller.RateLimitFlag, false)
		return resp, nil
	}
	if pc := ctrl.Proxy(); pc.IsEnabled() {
		req.URL = pc.BuildUrl(req.URL)
		if req.URL != nil {
			req.Host = req.URL.Host
		}
		for _, header := range pc.Headers() {
			req.Header.Add(header.Name, header.Value)
		}
	}
	resp, err, statusFlags := w.exchange(ctrl.Timeout(), req)
	if err != nil {
		return resp, err
	}
	if rc := ctrl.Retry(); rc.IsEnabled() {
		prevFlags := statusFlags
		retry, statusFlags = rc.IsRetryable(resp.StatusCode)
		if retry {
			ctrl.LogHttpEgress(start, time.Since(start), req, resp, prevFlags, false)
			start = time.Now()
			resp, err, statusFlags = w.exchange(ctrl.Timeout(), req)
		}
	}
	ctrl.LogHttpEgress(start, time.Since(start), req, resp, statusFlags, retry)
	return resp, err
}

func (w *controllerWrapper) exchange(tc controller.Timeout, req *http.Request) (resp *http.Response, err error, statusFlags string) {
	if tc == nil || !tc.IsEnabled() {
		resp, err = w.rt.RoundTrip(req)
		return
	}
	ctx, cancel := context.WithTimeout(req.Context(), tc.Duration())
	//defer cancel()
	req = req.Clone(ctx)
	resp, err = w.rt.RoundTrip(req)
	if w.deadlineExceeded(err) {
		resp = &http.Response{Request: req, StatusCode: tc.StatusCode()}
		err = nil
		statusFlags = controller.UpstreamTimeoutFlag
		cancel()
	}
	return
}

func (w *controllerWrapper) deadlineExceeded(err error) bool {
	return err != nil && errors.As(err, &context.DeadlineExceeded)
}

// ControllerWrapTransport - provides a RoundTrip wrapper that applies controller controllers
func ControllerWrapTransport(client *http.Client) {
	if client == nil || client == http.DefaultClient {
		if http.DefaultClient.Transport == nil {
			http.DefaultClient.Transport = &controllerWrapper{http.DefaultTransport}
		} else {
			http.DefaultClient.Transport = ControllerWrapRoundTripper(http.DefaultClient.Transport)
		}
	} else {
		if client.Transport == nil {
			client.Transport = &controllerWrapper{http.DefaultTransport}
		} else {
			client.Transport = ControllerWrapRoundTripper(client.Transport)
		}
	}
}

// ControllerWrapRoundTripper - provides a RoundTrip wrapper that applies controller controllers
func ControllerWrapRoundTripper(rt http.RoundTripper) http.RoundTripper {
	return &controllerWrapper{rt}
}
