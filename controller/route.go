package controller

import (
	"golang.org/x/time/rate"
)

// Route - route data
type Route struct {
	Name        string
	Pattern     string
	Traffic     string // egress/ingress
	Ping        bool   // health traffic
	Protocol    string // gRPC, HTTP10, HTTP11, HTTP2, HTTP3
	Timeout     *TimeoutConfig
	RateLimiter *RateLimiterConfig
	Retry       *RetryConfig
	Failover    *FailoverConfig
	Proxy       *ProxyConfig
}

type TimeoutConfigJson struct {
	Enabled    bool
	StatusCode int
	Duration   string
}

type RetryConfigJson struct {
	Enabled bool
	Limit   rate.Limit
	Burst   int
	Wait    string
	Codes   []int
}

type RouteConfig struct {
	Name        string
	Pattern     string
	Traffic     string // Egress/Ingress
	Ping        bool   // Health traffic
	Protocol    string // gRPC, HTTP10, HTTP11, HTTP2, HTTP3
	Timeout     *TimeoutConfigJson
	RateLimiter *RateLimiterConfig
	Retry       *RetryConfigJson
	Failover    *FailoverConfig
	Proxy       *ProxyConfig
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
		case *ProxyConfig:
			route.Proxy = c
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
	route.Proxy = config.Proxy
	route.RateLimiter = config.RateLimiter
	if config.Timeout != nil {
		duration, err := ParseDuration(config.Timeout.Duration)
		if err != nil {
			return Route{}, err
		}
		route.Timeout = NewTimeoutConfig(config.Timeout.Enabled, config.Timeout.StatusCode, duration)
	}
	if config.Retry != nil {
		duration, err := ParseDuration(config.Retry.Wait)
		if err != nil {
			return Route{}, err
		}
		route.Retry = NewRetryConfig(config.Retry.Enabled, config.Retry.Limit, config.Retry.Burst, duration, config.Retry.Codes)
	}
	return route, nil
}

func (r Route) IsConfigured() bool {
	return r.Retry != nil || r.Timeout != nil || r.RateLimiter != nil || r.Failover != nil || r.Proxy != nil
}
