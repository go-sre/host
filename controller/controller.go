package controller

import (
	"errors"
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
	RateLimiter() (RateLimiter, bool)
	Retry() (Retry, bool)
	Failover() (Failover, bool)
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
	timeout     *timeout
	rateLimiter *rateLimiter
	failover    *failover
	retry       *retry
	proxy       *proxy
}

func cloneController[T *timeout | *rateLimiter | *retry | *proxy | *failover](curr *controller, item T) *controller {
	newC := new(controller)
	*newC = *curr
	switch i := any(item).(type) {
	case *timeout:
		newC.timeout = i
	case *rateLimiter:
		newC.rateLimiter = i
	case *failover:
		newC.failover = i
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
	ctrl := new(controller)
	ctrl.name = route.Name
	ctrl.ping = route.Ping
	if route.Timeout != nil {
		ctrl.timeout = newTimeout(route.Name, t, route.Timeout)
		err = ctrl.timeout.validate()
		if err != nil {
			errs = append(errs, err)
		}
	} else {
		ctrl.timeout = disabledTimeout
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
	if route.Failover != nil {
		ctrl.failover = newFailover(route.Name, t, route.Failover)
		err = ctrl.failover.validate()
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
	} else {
		ctrl.proxy = disabledProxy
	}
	return ctrl, errs
}

func newDefaultController(name string) *controller {
	ctrl := new(controller)
	ctrl.name = name
	ctrl.timeout = disabledTimeout
	ctrl.proxy = disabledProxy
	return ctrl
}

func newNilController(name string) *controller {
	ctrl := new(controller)
	ctrl.name = name
	ctrl.timeout = disabledTimeout
	ctrl.proxy = disabledProxy
	return ctrl
}

func (c *controller) validate(egress bool) error {
	if !egress {
		if c.failover != nil {
			return errors.New("invalid configuration: Failover is not valid for ingress traffic")
		}
		if c.retry != nil {
			return errors.New("invalid configuration: Retry is not valid for ingress traffic")
		}
		if c.name == HostControllerName {
			if c.timeout != nil {
				return errors.New("invalid configuration: Timeout is not valid for host controller")
			}
		} else {
			if c.rateLimiter != nil {
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

func (c *controller) RateLimiter() (RateLimiter, bool) {
	if c.rateLimiter == nil {
		return nil, false
	}
	return c.rateLimiter, true
}

func (c *controller) Retry() (Retry, bool) {
	if c.retry == nil {
		return nil, false
	}
	return c.retry, true
}

func (c *controller) Failover() (Failover, bool) {
	if c.failover == nil {
		return nil, false
	}
	return c.failover, true
}

func (c *controller) Proxy() Proxy {
	return c.proxy
}

func (c *controller) Signal(values url.Values) error {
	return nil
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
	failoverState(state, c.failover)
	retryState(state, c.retry, retry)
	proxyState(state, c.proxy)

	defaultLogFn(EgressTraffic, start, duration, req, resp, statusFlags, state)
}

func (c *controller) LogEgress(start time.Time, duration time.Duration, statusCode int, uri, requestId, method, statusFlags string) {
	state := c.state()
	failoverState(state, c.failover)
	//retryState(state, c.retry, false)
	//proxyState(state, c.proxy)

	req, _ := http.NewRequest(method, uri, nil)
	req.Header.Add(RequestIdHeaderName, requestId)

	resp := new(http.Response)
	resp.StatusCode = statusCode
	defaultLogFn(EgressTraffic, start, duration, req, resp, statusFlags, state)
}
