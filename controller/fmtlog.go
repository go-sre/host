package controller

import (
	"fmt"
	"golang.org/x/time/rate"
	"net/http"
	"strconv"
	"time"
)

func FmtLog(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName string, timeout int, rateLimit rate.Limit, rateBurst int, rateThreshold, retry, proxy, proxyThreshold, statusFlags string) string {
	d := int(duration / time.Duration(1e6))
	s := fmt.Sprintf("start:%v ,"+
		"duration:%v ,"+
		"traffic:%v, "+
		"route:%v, "+
		"request-id:%v, "+
		"protocol:%v, "+
		"method:%v, "+
		"url:%v, "+
		"host:%v, "+
		"path:%v, "+
		"status-code:%v, "+
		"timeout-ms:%v, "+
		"rate-limit:%v, "+
		"rate-burst:%v, "+
		"rate-threshold:%v, "+
		"retry:%v, "+
		"proxy:%v, "+
		"proxy-threshold:%v, "+
		"status-flags:%v",
		FmtTimestamp(start), //l.Value(StartTimeOperator),
		strconv.Itoa(d),     //l.Value(DurationOperator),
		traffic,             //l.Value(TrafficOperator),
		routeName,           //l.Value(RouteNameOperator),

		req.Header.Get(RequestIdHeaderName), //l.Value(RequestIdOperator),
		req.Proto,                           //l.Value(RequestProtocolOperator),
		req.Method,                          //l.Value(RequestMethodOperator),
		req.URL.String(),                    //l.Value(RequestUrlOperator),
		req.URL.Host,                        //l.Value(RequestHostOperator),
		req.URL.Path,                        //l.Value(RequestPathOperator),

		resp.StatusCode, //l.Value(ResponseStatusCodeOperator),

		timeout, //Tl.Value(TimeoutDurationOperator),

		rateLimit, //l.Value(RateLimitOperator),
		rateBurst, //l.Value(RateBurstOperator),
		rateThreshold,

		retry,
		proxy,
		proxyThreshold,
		statusFlags, //l.Value(StatusFlagsOperator),
	)

	return s
}
