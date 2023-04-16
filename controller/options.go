package controller

import (
	"fmt"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

// HttpMatcher - type for Ingress/Egress table lookups by request
type HttpMatcher func(req *http.Request) (routeName string, ok bool)

// UriMatcher - type for Ingress/Egress table lookups by uri
type UriMatcher func(uri string, method string) (routeName string, ok bool)

// OutputHandler - type for output handling
type OutputHandler func(traffic string, start time.Time, duration time.Duration, routeName string, req *http.Request, resp *http.Response, timeout int, rateLimit rate.Limit, rateBurst int, proxied string, statusFlags string)

// SetLogFn - configuration for logging function
func SetLogFn(fn OutputHandler) {
	if fn != nil {
		defaultLogFn = fn
	}
}

var defaultLogFn = func(traffic string, start time.Time, duration time.Duration, routeName string, req *http.Request, resp *http.Response, timeout int, rateLimit rate.Limit, rateBurst int, proxied string, statusFlags string) {
	s := FmtLog(traffic, start, duration, routeName, req, resp, timeout, rateLimit, rateBurst, proxied, statusFlags)
	fmt.Printf("{%v}\n", s)
}

// SetExtractFn - configuration for extract function
func SetExtractFn(fn OutputHandler) {
	if fn != nil {
		defaultExtractFn = fn
	}
}

var defaultExtractFn OutputHandler
