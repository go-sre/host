package connector

import (
	"errors"
	"github.com/go-sre/core/runtime"
	"github.com/go-sre/host/accessdata"
	"github.com/go-sre/host/controller"
	"golang.org/x/time/rate"
	"net/http"
	url2 "net/url"
	"strings"
	"time"
)

type messageHandler func(l *accessdata.Entry) bool

var (
	pushLocInit      = pkgPath + "/initialize-push"
	pushLocDo        = pkgPath + "/do"
	pushUrl          string
	pushC            chan *accessdata.Entry
	pushClient                      = http.DefaultClient
	pushHandler      messageHandler = pushDo
	pushErrorHandler runtime.ErrorHandleFn
	operators        = []accessdata.Operator{
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
		{Name: "rate-threshold", Value: accessdata.RateThresholdOperator},
		{Name: "retry", Value: accessdata.RetryOperator},
		{Name: "proxy", Value: accessdata.ProxyOperator},
		{Name: "proxy-threshold", Value: accessdata.ProxyThresholdOperator},
		{Name: "status-flags", Value: accessdata.StatusFlagsOperator},
	}
)

func InitializePush[E runtime.ErrorHandler](uri string, newClient *http.Client) *runtime.Status {
	pushErrorHandler = runtime.NewErrorHandler[E]()
	if uri == "" {
		return pushErrorHandler(nil, pushLocInit, errors.New("invalid argument: uri is empty"))
	}
	u, err1 := url2.Parse(uri)
	if err1 != nil {
		return pushErrorHandler(nil, pushLocInit, err1)
	}
	pushUrl = u.String()
	pushC = make(chan *accessdata.Entry, 100)
	go pushReceive()
	if newClient != nil {
		pushClient = newClient
	}
	controller.SetExtractFn(extract)
	return runtime.NewStatusOK()
}

func ShutdownPush() {
	if pushC != nil {
		close(pushC)
	}
}

func extract(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName string, timeout int, limit rate.Limit, burst int, rateThreshold, retry, proxy, proxyThreshold, statusFlags string) {
	pushC <- accessdata.NewEntry(traffic, start, duration, req, resp, routeName, timeout, limit, burst, rateThreshold, retry, proxy, proxyThreshold, statusFlags)
}

func pushDo(entry *accessdata.Entry) bool {
	if entry == nil {
		pushErrorHandler(nil, pushLocDo, errors.New("invalid argument: access log data is nil"))
		return false
	}
	// let's not extract the extractor, the extractor, the extractor ...
	if entry.Url == pushUrl {
		return false
	}
	var req *http.Request
	var err error

	reader := strings.NewReader(accessdata.WriteJson(operators, entry))
	req, err = http.NewRequest(http.MethodPut, pushUrl, reader)
	if err == nil {
		_, err = pushClient.Do(req)
	}
	if err != nil {
		pushErrorHandler(nil, pushLocDo, err)
		return false
	}
	return true
}

func pushReceive() {
	for {
		select {
		case msg, open := <-pushC:
			if !open {
				return
			}
			pushHandler(msg)
		default:
		}
	}
}
