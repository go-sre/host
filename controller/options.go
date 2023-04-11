package controller

import (
	"fmt"
	"net/http"
	"time"
)

// HttpMatcher - type for Ingress/Egress table lookups by request
type HttpMatcher func(req *http.Request) (routeName string, ok bool)

// UriMatcher - type for Ingress/Egress table lookups by uri
type UriMatcher func(uri string, method string) (routeName string, ok bool)

// OutputHandler - type for output handling
type OutputHandler func(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, statusFlags string, controllerState map[string]string)

// SetLogFn - configuration for logging function
func SetLogFn(fn OutputHandler) {
	if fn != nil {
		defaultLogFn = fn
	}
}

var defaultLogFn = func(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, statusFlags string, controllerState map[string]string) {
	s := FmtLog(traffic, start, duration, req, resp, statusFlags, controllerState)
	fmt.Printf("{%v}\n", s)
}

// SetExtractFn - configuration for extract function
func SetExtractFn(fn OutputHandler) {
	if fn != nil {
		defaultExtractFn = fn
	}
}

var defaultExtractFn OutputHandler
