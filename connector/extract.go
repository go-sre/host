package connector

import (
	"errors"
	"github.com/go-sre/core/runtime"
	"github.com/go-sre/host/accessdata"
	"github.com/go-sre/host/controller"
	"golang.org/x/time/rate"
	"net/http"
	url2 "net/url"
	"reflect"
	"strings"
	"time"
)

type messageHandler func(l *accessdata.Entry) bool
type pkg struct{}

var (
	pkgPath      = reflect.TypeOf(any(pkg{})).PkgPath()
	locInit      = pkgPath + "/initialize"
	locDo        = pkgPath + "/do"
	url          string
	c            chan *accessdata.Entry
	client                      = http.DefaultClient
	handler      messageHandler = do
	errorHandler runtime.ErrorHandleFn
	operators    = []accessdata.Operator{
		{Name: "start-time", Value: accessdata.StartTimeOperator},
		{Name: "duration-ms", Value: accessdata.DurationOperator},
		{Name: "traffic", Value: accessdata.TrafficOperator},
		{Name: "route-name", Value: accessdata.RouteNameOperator},

		{Name: "region", Value: accessdata.OriginRegionOperator},
		{Name: "zone", Value: accessdata.OriginZoneOperator},
		{Name: "sub-zone", Value: accessdata.OriginSubZoneOperator},
		{Name: "service", Value: accessdata.OriginServiceOperator},
		{Name: "instance-id", Value: accessdata.OriginInstanceIdOperator},

		{Name: "method", Value: accessdata.RequestMethodOperator},
		{Name: "url", Value: accessdata.RequestUrlOperator},
		{Name: "host", Value: accessdata.RequestHostOperator},
		{Name: "path", Value: accessdata.RequestPathOperator},
		{Name: "protocol", Value: accessdata.RequestProtocolOperator},
		{Name: "request-id", Value: accessdata.RequestIdOperator},
		{Name: "forwarded", Value: accessdata.RequestForwardedForOperator},

		{Name: "status-code", Value: accessdata.ResponseStatusCodeOperator},

		{Name: "timeout-ms", Value: accessdata.TimeoutDurationOperator},
		{Name: "rate-limit", Value: accessdata.RateLimitOperator},
		{Name: "rate-burst", Value: accessdata.RateBurstOperator},
		{Name: "retry", Value: accessdata.RetryOperator},
		{Name: "proxy", Value: accessdata.ProxyOperator},
		{Name: "status-flags", Value: accessdata.StatusFlagsOperator},
	}
)

func Initialize[E runtime.ErrorHandler](uri string, newClient *http.Client) *runtime.Status {
	errorHandler = runtime.NewErrorHandler[E]()
	if uri == "" {
		return errorHandler(nil, locInit, errors.New("invalid argument: uri is empty"))
	}
	u, err1 := url2.Parse(uri)
	if err1 != nil {
		return errorHandler(nil, locInit, err1)
	}
	url = u.String()
	c = make(chan *accessdata.Entry, 100)
	go receive()
	if newClient != nil {
		client = newClient
	}
	controller.SetExtractFn(extract)
	return runtime.NewStatusOK()
}

func Shutdown() {
	if c != nil {
		close(c)
	}
}

func extract(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName string, timeout int, limit rate.Limit, burst int, retry, proxy, statusFlags string) {
	c <- accessdata.NewEntry(traffic, start, duration, req, resp, routeName, timeout, limit, burst, retry, proxy, statusFlags)
}

func do(entry *accessdata.Entry) bool {
	if entry == nil {
		errorHandler(nil, locDo, errors.New("invalid argument: access log data is nil"))
		return false
	}
	// let's not extract the extractor, the extractor, the extractor ...
	if entry.Url == url {
		return false
	}
	var req *http.Request
	var err error

	reader := strings.NewReader(accessdata.WriteJson(operators, entry))
	req, err = http.NewRequest(http.MethodPut, url, reader)
	if err == nil {
		_, err = client.Do(req)
	}
	if err != nil {
		errorHandler(nil, locDo, err)
		return false
	}
	return true
}

func receive() {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			handler(msg)
		default:
		}
	}
}
