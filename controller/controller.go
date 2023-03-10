package controller

import (
	"errors"
	"github.com/google/uuid"
	"github.com/gotemplates/host/accessdata"
	"net/http"
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

// Controller - definition for properties of a controller
type Controller interface {
	Name() string
	Timeout() (Timeout, bool)
	RateLimiter() (RateLimiter, bool)
	Retry() (Retry, bool)
	Failover() (Failover, bool)
	UpdateHeaders(req *http.Request)
	LogHttpIngress(start time.Time, duration time.Duration, req *http.Request, statusCode int, written int64, statusFlags string)
	LogHttpEgress(start time.Time, duration time.Duration, req *http.Request, resp *http.Response, statusFlags string, retry bool)
	LogEgress(start time.Time, duration time.Duration, statusCode int, uri, requestId, method, statusFlags string)
	t() *controller
}

// Configuration - configuration for actuators
type Configuration interface {
	SetHttpMatcher(fn HttpMatcher)
	SetUriMatcher(fn UriMatcher)
	SetDefaultController(route Route) []error
	SetHostController(route Route) []error
	AddController(route Route) []error
}

// Controllers - public interface
type Controllers interface {
	Host() Controller
	Default() Controller
	LookupHttp(req *http.Request) Controller
	LookupUri(urn string, method string) Controller
	LookupByName(name string) Controller
}

// Table - controller table
type Table interface {
	Configuration
	Controllers
}

// IngressTable - table for ingress controllers
var IngressTable = NewIngressTable()

// EgressTable - table for egress controllers
var EgressTable = NewEgressTable()

type controller struct {
	name        string
	ping        bool
	timeout     *timeout
	rateLimiter *rateLimiter
	failover    *failover
	retry       *retry
}

func cloneController[T *timeout | *rateLimiter | *retry | *failover](curr *controller, item T) *controller {
	newC := new(controller)
	*newC = *curr
	switch i := any(item).(type) {
	case *timeout:
		newC.timeout = i
	case *rateLimiter:
		newC.rateLimiter = i
	case *failover:
		newC.failover = i
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
	return ctrl, errs
}

func newDefaultController(name string) *controller {
	return &controller{name: name}
}

func newNilController(name string) *controller {
	return &controller{name: name}
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

func (c *controller) Timeout() (Timeout, bool) {
	if c.timeout == nil {
		return nil, false
	}
	return c.timeout, true
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

func (c *controller) t() *controller {
	return c
}

func (c *controller) state() map[string]string {
	state := make(map[string]string, 12)
	state[accessdata.ControllerName] = c.Name()
	if c.ping {
		state[accessdata.PingName] = "true"
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
	if req.Header.Get(accessdata.RequestIdHeaderName) == "" {
		req.Header.Add(accessdata.RequestIdHeaderName, uuid.New().String())
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

	defaultLogFn("egress", start, duration, req, resp, statusFlags, state)
}

func (c *controller) LogEgress(start time.Time, duration time.Duration, statusCode int, uri, requestId, method, statusFlags string) {
	state := c.state()
	failoverState(state, c.failover)
	retryState(state, c.retry, false)

	req, _ := http.NewRequest(method, uri, nil)
	req.Header.Add(accessdata.RequestIdHeaderName, requestId)

	resp := new(http.Response)
	resp.StatusCode = statusCode
	defaultLogFn("egress", start, duration, req, resp, statusFlags, state)
}
