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
type OutputHandler func(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName string, timeout int, rateLimit rate.Limit, rateBurst int, rateThreshold, retry, proxy, proxyThreshold, statusFlags string)

// SetLogFn - configuration for logging function
func SetLogFn(fn OutputHandler) {
	if fn != nil {
		defaultLogFn = fn
	}
}

var defaultLogFn = func(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName string, timeout int, rateLimit rate.Limit, rateBurst int, rateThreshold, retry, proxy, proxyThreshold, statusFlags string) {
	s := FmtLog(traffic, start, duration, req, resp, routeName, timeout, rateLimit, rateBurst, rateThreshold, retry, proxy, proxyThreshold, statusFlags)
	fmt.Printf("{%v}\n", s)
}

// SetExtractFn - configuration for connector function
func SetExtractFn(fn OutputHandler) {
	if fn != nil {
		defaultExtractFn = fn
	}
}

var defaultExtractFn OutputHandler
