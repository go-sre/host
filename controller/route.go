package controller

import (
	"golang.org/x/time/rate"
	"strconv"
	"strings"
	"time"
)

// Route - route data
type Route struct {
	Name        string
	Pattern     string
	Traffic     string // egress/ingress
	Ping        bool   // health traffic
	Protocol    string // gRPC, HTTP10, HTTP11, HTTP2, HTTP3gRPC, HTTP
	Timeout     *TimeoutConfig
	RateLimiter *RateLimiterConfig
	Retry       *RetryConfig
	Failover    *FailoverConfig
}

type TimeoutConfigJson struct {
	Duration   string
	StatusCode int
}

type RetryConfigJson struct {
	Limit rate.Limit
	Burst int
	Wait  string
	Codes []int
}

type RouteConfig struct {
	Name        string
	Pattern     string
	Traffic     string // Egress/Ingress
	Ping        bool   // Health traffic
	Protocol    string // gRPC, HTTP10, HTTP11, HTTP2, HTTP3gRPC, HTTP
	Timeout     *TimeoutConfigJson
	RateLimiter *RateLimiterConfig
	Retry       *RetryConfigJson
	Failover    *FailoverConfig
}

func newRoute(name string, config ...any) Route {
	return NewRoute(name, "", "", false, config...)
}

// NewRoute - creates a new route
func NewRoute(name string, traffic, protocol string, ping bool, config ...any) Route {
	route := Route{}
	route.Name = name
	route.Traffic = traffic
	route.Protocol = protocol
	route.Ping = ping
	for _, cfg := range config {
		if cfg == nil {
			continue
		}
		switch c := cfg.(type) {
		case *TimeoutConfig:
			route.Timeout = c
		case *RateLimiterConfig:
			route.RateLimiter = c
		case *FailoverConfig:
			route.Failover = c
		case *RetryConfig:
			route.Retry = c
		}
	}
	return route
}

// NewRouteFromConfig - creates a new route from configuration
func NewRouteFromConfig(config RouteConfig) (Route, error) {
	route := Route{}
	route.Name = config.Name
	route.Pattern = config.Pattern
	route.Traffic = config.Traffic
	route.Ping = config.Ping
	route.Protocol = config.Protocol
	route.Failover = config.Failover
	route.RateLimiter = config.RateLimiter
	if config.Timeout != nil {
		duration, err := ConvertDuration(config.Timeout.Duration)
		if err != nil {
			return Route{}, err
		}
		route.Timeout = NewTimeoutConfig(duration, config.Timeout.StatusCode)
	}
	if config.Retry != nil {
		duration, err := ConvertDuration(config.Retry.Wait)
		if err != nil {
			return Route{}, err
		}
		route.Retry = NewRetryConfig(config.Retry.Codes, config.Retry.Limit, config.Retry.Burst, duration)
	}
	return route, nil
}

func (r Route) IsConfigured() bool {
	return r.Retry != nil || r.Timeout != nil || r.RateLimiter != nil || r.Failover != nil
}

func ConvertDuration(s string) (time.Duration, error) {
	if s == "" {
		return 0, nil
	}
	tokens := strings.Split(s, "ms")
	if len(tokens) == 2 {
		val, err := strconv.Atoi(tokens[0])
		if err != nil {
			return 0, err
		}
		return time.Duration(val) * time.Millisecond, nil
	}
	tokens = strings.Split(s, "µs")
	if len(tokens) == 2 {
		val, err := strconv.Atoi(tokens[0])
		if err != nil {
			return 0, err
		}
		return time.Duration(val) * time.Microsecond, nil
	}
	tokens = strings.Split(s, "m")
	if len(tokens) == 2 {
		val, err := strconv.Atoi(tokens[0])
		if err != nil {
			return 0, err
		}
		return time.Duration(val) * time.Minute, nil
	}
	// Assume seconds
	tokens = strings.Split(s, "s")
	if len(tokens) == 2 {
		s = tokens[0]
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return time.Duration(val) * time.Second, nil
}
