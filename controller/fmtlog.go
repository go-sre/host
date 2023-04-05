package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func FmtLog(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, statusFlags string, controllerState map[string]string) string {
	if controllerState == nil {
		controllerState = make(map[string]string)
	}
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
		"timeout_ms:%v, "+
		"rate-limit:%v, "+
		"rate-burst:%v, "+
		"retry:%v, "+
		"retry-rate-limit:%v, "+
		"retry-rate-burst:%v, "+
		"status-flags:%v",
		FmtTimestamp(start),             //l.Value(StartTimeOperator),
		strconv.Itoa(d),                 //l.Value(DurationOperator),
		traffic,                         //l.Value(TrafficOperator),
		controllerState[ControllerName], //l.Value(RouteNameOperator),

		req.Header.Get(RequestIdHeaderName), //l.Value(RequestIdOperator),
		req.Proto,                           //l.Value(RequestProtocolOperator),
		req.Method,                          //l.Value(RequestMethodOperator),
		req.URL.String(),                    //l.Value(RequestUrlOperator),
		req.URL.Host,                        //l.Value(RequestHostOperator),
		req.URL.Path,                        //l.Value(RequestPathOperator),

		resp.StatusCode, //l.Value(ResponseStatusCodeOperator),

		controllerState[TimeoutName], //Tl.Value(TimeoutDurationOperator),

		controllerState[RateLimitName], //l.Value(RateLimitOperator),
		controllerState[RateBurstName], //l.Value(RateBurstOperator),

		controllerState[RetryName],          //l.Value(RetryOperator),
		controllerState[RetryRateLimitName], //l.Value(RetryRateLimitOperator),
		controllerState[RetryRateBurstName], //l.Value(RetryRateBurstOperator),

		statusFlags, //l.Value(StatusFlagsOperator),
	)

	return s
}
