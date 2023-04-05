package controller

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"time"
)

const (
	RateLimitInfValue = 99999

	HostControllerName    = "host"
	DefaultControllerName = "*"
	NilControllerName     = "!"
	FromRouteHeaderName   = "from-route"

	RateLimitFlag       = "RL"
	UpstreamTimeoutFlag = "UT"
	HostTimeoutFlag     = "HT"
	NotEnabledFlag      = "NE"
)

// State - defines enabled state
type State interface {
	IsEnabled() bool
	Enable()
	Disable()
}

// Controller - definition for properties of a controller
type Controller interface {
	Actuator
	Name() string
	Timeout() Timeout
	RateLimiter() RateLimiter
	Retry() Retry
	Proxy() Proxy
	UpdateHeaders(req *http.Request)
	LogHttpIngress(start time.Time, duration time.Duration, req *http.Request, statusCode int, written int64, statusFlags string)
	LogHttpEgress(start time.Time, duration time.Duration, req *http.Request, resp *http.Response, statusFlags string, retry bool)
	LogEgress(start time.Time, duration time.Duration, statusCode int, uri, requestId, method, statusFlags string)
	t() *controller
}

type controller struct {
	name        string
	ping        bool
	tbl         *table
	timeout     *timeout
	rateLimiter *rateLimiter
	retry       *retry
	proxy       *proxy
}

func cloneController[T *timeout | *rateLimiter | *retry | *proxy](curr *controller, item T) *controller {
	newC := new(controller)
	*newC = *curr
	switch i := any(item).(type) {
	case *timeout:
		newC.timeout = i
	case *rateLimiter:
		newC.rateLimiter = i
	case *proxy:
		newC.proxy = i
	case *retry:
		newC.retry = i
	default:
	}
	return newC
}

func newController(route Route, t *table) (*controller, []error) {
	var errs []error
	var err error
	ctrl := newDefaultController(route.Name)
	ctrl.ping = route.Ping
	ctrl.tbl = t
	if route.Timeout != nil {
		ctrl.timeout = newTimeout(route.Name, t, route.Timeout)
		err = ctrl.timeout.validate()
		if err != nil {
			errs = append(errs, err)
		}
	}
	if route.RateLimiter != nil {
		ctrl.rateLimiter = newRateLimiter(route.Name, t, route.RateLimiter)
		err = ctrl.rateLimiter.validate()
		if err != nil {
			errs = append(errs, err)
		}
	}
	if route.Retry != nil {
		ctrl.retry = newRetry(route.Name, t, route.Retry)
		err = ctrl.retry.validate()
		if err != nil {
			errs = append(errs, err)
		}
	}
	if route.Proxy != nil {
		ctrl.proxy = newProxy(route.Name, t, route.Proxy)
		err = ctrl.proxy.validate()
		if err != nil {
			errs = append(errs, err)
		}
	}
	return ctrl, errs
}

func newDefaultController(name string) *controller {
	ctrl := new(controller)
	ctrl.name = name
	ctrl.timeout = disabledTimeout
	ctrl.proxy = disabledProxy
	ctrl.rateLimiter = disabledRateLimiter
	ctrl.retry = disabledRetry
	return ctrl
}

func (c *controller) validate(egress bool) error {
	if !egress {
		if c.retry.IsEnabled() {
			return errors.New("invalid configuration: Retry is not valid for ingress traffic")
		}
		if c.name == HostControllerName {
			if c.timeout.IsEnabled() {
				return errors.New("invalid configuration: Timeout is not valid for host controller")
			}
		} else {
			if c.rateLimiter.IsEnabled() {
				return errors.New("invalid configuration: RateLimiter is not valid for ingress traffic")
			}
		}
	}
	return nil
}

func (c *controller) Name() string {
	return c.name
}

func (c *controller) Timeout() Timeout {
	return c.timeout
}

func (c *controller) RateLimiter() RateLimiter {
	return c.rateLimiter
}

func (c *controller) Retry() Retry {
	return c.retry
}

func (c *controller) Proxy() Proxy {
	return c.proxy
}

func (c *controller) Signal(values url.Values) error {
	if values == nil {
		return nil
	}
	switch values.Get(BehaviorKey) {
	case BehaviorTimeout:
		return c.Timeout().Signal(values)
		break
	case BehaviorRetry:
		return c.Retry().Signal(values)
		break
	case BehaviorRateLimit:
		return c.RateLimiter().Signal(values)
		break
	case BehaviorProxy:
		return c.Proxy().Signal(values)
		break
	}
	return errors.New(fmt.Sprintf("invalid argument: behavior [%s] is not supported", values.Get(BehaviorKey)))
}

func (c *controller) t() *controller {
	return c
}

func (c *controller) state() map[string]string {
	state := make(map[string]string, 12)
	state[ControllerName] = c.Name()
	if c.ping {
		state[PingName] = "true"
	}
	timeoutState(state, c.timeout)
	rateLimiterState(state, c.rateLimiter)
	return state
}

func (c *controller) UpdateHeaders(req *http.Request) {
	if req == nil || req.Header == nil {
		return
	}
	req.Header.Add(FromRouteHeaderName, c.name)
	if req.Header.Get(RequestIdHeaderName) == "" {
		req.Header.Add(RequestIdHeaderName, uuid.New().String())
	}
}

func (c *controller) LogHttpIngress(start time.Time, duration time.Duration, req *http.Request, statusCode int, written int64, statusFlags string) {
	if c.name == NilControllerName {
		return
	}
	resp := new(http.Response)
	resp.StatusCode = statusCode
	resp.ContentLength = written
	defaultLogFn("ingress", start, duration, req, resp, statusFlags, c.state())
}

func (c *controller) LogHttpEgress(start time.Time, duration time.Duration, req *http.Request, resp *http.Response, statusFlags string, retry bool) {
	if c.name == NilControllerName {
		return
	}
	state := c.state()
	retryState(state, c.retry, retry)
	proxyState(state, c.proxy)

	defaultLogFn(EgressTraffic, start, duration, req, resp, statusFlags, state)
}

func (c *controller) LogEgress(start time.Time, duration time.Duration, statusCode int, uri, requestId, method, statusFlags string) {
	state := c.state()
	//retryState(state, c.retry, false)
	//proxyState(state, c.proxy)

	req, _ := http.NewRequest(method, uri, nil)
	req.Header.Add(RequestIdHeaderName, requestId)

	resp := new(http.Response)
	resp.StatusCode = statusCode
	defaultLogFn(EgressTraffic, start, duration, req, resp, statusFlags, state)
}
