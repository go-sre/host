package accessdata

import (
	"fmt"
	"golang.org/x/time/rate"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	EgressTraffic  = "egress"
	IngressTraffic = "ingress"
	PingTraffic    = "ping"

	/*
		PingName           = "ping"
		TimeoutName        = "timeout"
		ProxyName          = "proxy"
		RetryName          = "retry"
		RetryRateLimitName = "retryRateLimit"
		RetryRateBurstName = "retryBurst"
		RateLimitName      = "rateLimit"
		RateBurstName      = "burst"
		ControllerName     = "name"

	*/
)

// Accessor - function type
type Accessor func(entry *Entry)

// Entry - struct for all access logging accessdata
type Entry struct {
	Traffic   string
	Start     time.Time
	Duration  time.Duration
	RouteName string //CtrlState map[string]string

	// Request
	Url      string
	Path     string
	Host     string
	Protocol string
	Method   string
	Header   http.Header

	// Response
	StatusCode    int
	BytesSent     int64
	BytesReceived int64

	// State and
	Timeout     int
	RateLimit   rate.Limit
	RateBurst   int
	Proxied     string
	StatusFlags string
}

func NewEmptyEntry() *Entry {
	return new(Entry)
}

func NewEntry(traffic string, start time.Time, duration time.Duration, routeName string, req *http.Request, resp *http.Response, timeout int, rateLimit rate.Limit, rateBurst int, proxied string, statusFlags string) *Entry {
	e := new(Entry)
	e.Traffic = traffic
	e.Start = start
	e.Duration = duration
	e.RouteName = routeName
	//if controllerState == nil {
	//	controllerState = make(map[string]string, 1)
	//} else {
	//	if s, ok := controllerState[PingName]; ok && s == "true" {
	//		e.Traffic = PingTraffic
	//	}
	//}
	//e.CtrlState = controllerState
	e.AddRequest(req)
	e.AddResponse(resp)

	e.Timeout = timeout
	e.RateLimit = rateLimit
	e.RateBurst = rateBurst
	e.Proxied = proxied
	e.StatusFlags = statusFlags
	return e
}

// NewEgressEntry - create an Entry for egress traffic
func NewEgressEntry(start time.Time, duration time.Duration, routeName string, req *http.Request, resp *http.Response, timeout int, rateLimit rate.Limit, rateBurst int, proxied string, statusFlags string) *Entry {
	return NewEntry(EgressTraffic, start, duration, routeName, req, resp, timeout, rateLimit, rateBurst, proxied, statusFlags)
}

// NewIngressEntry - create an Entry for ingress traffic
func NewIngressEntry(start time.Time, duration time.Duration, routeName string, req *http.Request, resp *http.Response, timeout int, rateLimit rate.Limit, rateBurst int, proxied string, statusFlags string) *Entry {
	return NewEntry(IngressTraffic, start, duration, routeName, req, resp, timeout, rateLimit, rateBurst, proxied, statusFlags)
}

func (l *Entry) AddResponse(resp *http.Response) {
	if resp == nil {
		return
	}
	l.StatusCode = resp.StatusCode
	l.BytesReceived = resp.ContentLength
}

func (l *Entry) AddUrl(uri string) {
	if uri == "" {
		return
	}
	u, err := url.Parse(uri)
	if err != nil {
		l.Url = err.Error()
		return
	}
	if u.Scheme == "urn" && u.Host == "" {
		l.Url = uri
		l.Protocol = u.Scheme
		t := strings.Split(u.Opaque, ":")
		if len(t) == 1 {
			l.Host = t[0]
		} else {
			l.Host = t[0]
			l.Path = t[1]
		}
	} else {
		l.Protocol = u.Scheme
		l.Url = u.String()
		l.Path = u.Path
		l.Host = u.Host
	}
}

func (l *Entry) AddRequest(req *http.Request) {
	if req == nil {
		return
	}
	l.Protocol = req.Proto
	l.Method = req.Method
	if req.Header != nil {
		l.Header = req.Header.Clone()
	}
	if req.URL == nil {
		return
	}
	if req.URL.Scheme == "urn" {
		l.AddUrl(req.URL.String())
	} else {
		l.Url = req.URL.String()
		l.Path = req.URL.Path
		if req.Host == "" {
			l.Host = req.URL.Host
		} else {
			l.Host = req.Host
		}
	}
}

func (l *Entry) Value(value string) string {
	switch value {
	case TrafficOperator:
		return l.Traffic
	case StartTimeOperator:
		return FmtTimestamp(l.Start)
	case DurationOperator:
		d := int(l.Duration / time.Duration(1e6))
		return strconv.Itoa(d)
	case DurationStringOperator:
		return l.Duration.String()

		// Origin
	case OriginRegionOperator, OriginZoneOperator, OriginSubZoneOperator, OriginServiceOperator, OriginInstanceIdOperator:
		return originValue(value)

		// Request
	case RequestMethodOperator:
		return l.Method
	case RequestProtocolOperator:
		return l.Protocol
	case RequestPathOperator:
		return l.Path
	case RequestUrlOperator:
		return l.Url
	case RequestHostOperator:
		return l.Host
	case RequestIdOperator:
		return l.Header.Get(RequestIdHeaderName)
	case RequestFromRouteOperator:
		return l.Header.Get(FromRouteHeaderName)
	case RequestUserAgentOperator:
		return l.Header.Get(UserAgentHeaderName)
	case RequestAuthorityOperator:
		return ""
	case RequestForwardedForOperator:
		return l.Header.Get(ForwardedForHeaderName)

		// Response
	case StatusFlagsOperator:
		return l.StatusFlags
	case ResponseBytesReceivedOperator:
		return strconv.Itoa(int(l.BytesReceived))
	case ResponseBytesSentOperator:
		return fmt.Sprintf("%v", l.BytesSent)
	case ResponseStatusCodeOperator:
		return strconv.Itoa(l.StatusCode)

	// Controller State
	case RouteNameOperator:
		return l.RouteName
	case TimeoutDurationOperator:
		return strconv.Itoa(l.Timeout)
	case RateLimitOperator:
		return fmt.Sprintf("%v", l.RateLimit)
	case RateBurstOperator:
		return strconv.Itoa(l.RateBurst)
	case ProxyOperator:
		return l.Proxied
		//case RetryOperator:
		//	return l.CtrlState[RetryName]
		//case RetryRateLimitOperator:
		//		return l.CtrlState[RetryRateLimitName]
		//	case RetryRateBurstOperator:
		//		return l.CtrlState[RetryRateBurstName]
	}
	if strings.HasPrefix(value, RequestReferencePrefix) {
		name := requestOperatorHeaderName(value)
		return l.Header.Get(name)
	}
	if !strings.HasPrefix(value, OperatorPrefix) {
		return value
	}
	return ""
}

func (l *Entry) String() string {
	return fmt.Sprintf( //"start:%v ,"+
		//"duration:%v ,"+
		"traffic:%v, "+
			"route:%v, "+
			"request-id:%v, "+
			"status-code:%v, "+
			"protocol:%v, "+
			"method:%v, "+
			"url:%v, "+
			"host:%v, "+
			"path:%v, "+
			"timeout:%v, "+
			"rate-limit:%v, "+
			"rate-burst:%v, "+
			"proxy:%v, "+
			"status-flags:%v",
		//l.Value(StartTimeOperator),
		//l.Value(DurationOperator),
		l.Value(TrafficOperator),
		l.Value(RouteNameOperator),

		l.Value(RequestIdOperator),
		l.Value(ResponseStatusCodeOperator),
		l.Value(RequestProtocolOperator),
		l.Value(RequestMethodOperator),
		l.Value(RequestUrlOperator),
		l.Value(RequestHostOperator),
		l.Value(RequestPathOperator),

		l.Value(TimeoutDurationOperator),
		l.Value(RateLimitOperator),
		l.Value(RateBurstOperator),

		l.Value(ProxyOperator),

		l.Value(StatusFlagsOperator),
	)
}
