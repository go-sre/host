package middleware

import (
	"fmt"
	"github.com/gotemplates/host/controller"
	"github.com/gotemplates/host/shared"
	"net/http"
	"time"
)

var (
	accessLogging  = false
	isEnabled2     = false
	timeoutRoute   = "timeout-route"
	rateLimitRoute = "rate-limit-route"
	retryRoute     = "retry-route"
	//googleUrl      = "https://www.google.com/search?q=test"
	twitterUrl  = "https://www.twitter.com"
	facebookUrl = "https://www.facebook.com"
	//instagramUrl   = "https://www.instagram.com"

	/*
		config = []data.Operator{
			//{Value: data.StartTimeOperator},
			//{Value: data.DurationOperator},
			{Value: data.TrafficOperator},
			{Value: data.RouteNameOperator},

			{Value: data.RequestMethodOperator},
			{Value: data.RequestHostOperator},
			{Value: data.RequestPathOperator},
			{Value: data.RequestProtocolOperator},

			{Value: data.ResponseStatusCodeOperator},
			{Value: data.StatusFlagsOperator},
			{Value: data.ResponseBytesReceivedOperator},
			{Value: data.ResponseBytesSentOperator},

			{Value: data.TimeoutDurationOperator},
			{Value: data.RateLimitOperator},
			{Value: data.RateBurstOperator},
			{Value: data.RetryOperator},
			{Value: data.RetryRateLimitOperator},
			{Value: data.RetryRateBurstOperator},
			{Value: data.FailoverOperator},
		}

	*/
)

func testHttpLog(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, statusFlags string, actuatorState map[string]string) {
	s := fmt.Sprintf("\"traffic\":\"%v\","+
		"\"route_name\":\"%v\","+
		"\"method\":\"%v\","+
		"\"host\":\"%v\","+
		"\"path\":\"%v\","+
		"\"protocol\":\"%v\","+
		"\"status_code\":%v,"+
		"\"status_flags\":\"%v\","+
		"\"bytes_received\":-1,"+
		"\"bytes_sent\":0,"+
		"\"timeout_ms\":%v,"+
		"\"rate-limit\":%v,"+
		"\"rate-burst\":%v,"+
		"\"retry\":%v,"+
		"\"retry-rate-limit\":%v,"+
		"\"retry-rate-burst\":%v,"+
		"\"failover\":%v",
		traffic, actuatorState[shared.ControllerName], req.Method, req.Host, req.URL.Path, req.Proto, resp.StatusCode, statusFlags,
		actuatorState[shared.TimeoutName],
		actuatorState[shared.RateLimitName], actuatorState[shared.RateBurstName],
		actuatorState[shared.RetryName], actuatorState[shared.RetryRateLimitName], actuatorState[shared.RetryRateBurstName],
		actuatorState[shared.FailoverName])
	fmt.Printf("test: Write() -> [{%v}]\n", s)
}

func init() {
	//err := log.InitEgressOperators(config)
	//if err != nil {
	//	fmt.Printf("init() -> [:%v]\n", err)
	//}
	controller.EgressTable.SetHttpMatcher(func(req *http.Request) (string, bool) {
		if req == nil {
			return "", true
		}
		if req.URL.String() == twitterUrl {
			return rateLimitRoute, true
		}
		if req.URL.String() == googleUrl {
			return timeoutRoute, true
		}
		if req.URL.String() == facebookUrl {
			return retryRoute, true
		}
		return "", true
	})

	controller.EgressTable.AddController(controller.NewRoute(timeoutRoute, shared.EgressTraffic, "", false, controller.NewTimeoutConfig(time.Millisecond, 504)))
	controller.EgressTable.AddController(controller.NewRoute(rateLimitRoute, shared.EgressTraffic, "", false, controller.NewRateLimiterConfig(2000, 0, 503)))
	controller.EgressTable.AddController(controller.NewRoute(retryRoute, shared.EgressTraffic, "", false, controller.NewTimeoutConfig(time.Millisecond, 504), controller.NewRetryConfig([]int{503, 504}, 0, 0, 0)))

	//	controller.SetLogFn(func(entry *data.Entry) {
	//		log.Write[log.TestOutputHandler, data.JsonFormatter](entry)
	//	},
	//	)

	controller.SetLogFn(testHttpLog)

}

func Example_Controller_Default_Actuator() {
	act := controller.EgressTable.LookupHttp(nil)
	fmt.Printf("test: LookupHttp(nil) -> [name:%v]\n", act.Name())

	//Output:
	//test: LookupHttp(nil) -> [name:*]

}

func Example_Controller_No_Wrapper() {
	req, _ := http.NewRequest("GET", googleUrl, nil)

	// Testing - check for a nil wrapper or round tripper
	w := controllerWrapper{}
	resp, err := w.RoundTrip(req)
	fmt.Printf("test: RoundTrip(wrapper:nil) -> [resp:%v] [err:%v]\n", resp, err)

	// Testing - no wrapper, calling Google search
	resp, err = http.DefaultClient.Do(req)
	fmt.Printf("test: RoundTrip(handler:false) -> [status_code:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: RoundTrip(wrapper:nil) -> [resp:<nil>] [err:invalid handler round tripper configuration : http.RoundTripper is nil]
	//test: RoundTrip(handler:false) -> [status_code:200] [err:<nil>]

}

func Example_Controller_Default() {
	req, _ := http.NewRequest("GET", instagramUrl, nil)

	if !isEnabled2 {
		isEnabled2 = true
		ControllerWrapTransport(nil)
	}
	resp, err := http.DefaultClient.Do(req)
	fmt.Printf("test: RoundTrip(handler:true) -> [status_code:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: Write() -> [{"traffic":"egress","route_name":"*","method":"GET","host":"www.instagram.com","path":"","protocol":"HTTP/1.1","status_code":200,"status_flags":"","bytes_received":-1,"bytes_sent":0,"timeout_ms":-1,"rate-limit":-1,"rate-burst":-1,"retry":,"retry-rate-limit":-1,"retry-rate-burst":-1,"failover":}]
	//test: RoundTrip(handler:true) -> [status_code:200] [err:<nil>]

}

func Example_Controller_Default_Timeout() {
	req, _ := http.NewRequest("GET", googleUrl, nil)

	if !isEnabled2 {
		isEnabled2 = true
		ControllerWrapTransport(nil)
	}
	resp, err := http.DefaultClient.Do(req)
	fmt.Printf("test: RoundTrip(handler:true) -> [status_code:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: Write() -> [{"traffic":"egress","route_name":"timeout-route","method":"GET","host":"www.google.com","path":"/search","protocol":"HTTP/1.1","status_code":504,"status_flags":"UT","bytes_received":-1,"bytes_sent":0,"timeout_ms":1,"rate-limit":-1,"rate-burst":-1,"retry":,"retry-rate-limit":-1,"retry-rate-burst":-1,"failover":}]
	//test: RoundTrip(handler:true) -> [status_code:504] [err:<nil>]

}

func Example_Controller_Default_RateLimit() {
	req, _ := http.NewRequest("GET", twitterUrl, nil)

	if !isEnabled2 {
		isEnabled2 = true
		ControllerWrapTransport(nil)
	}
	resp, err := http.DefaultClient.Do(req)
	fmt.Printf("test: RoundTrip(handler:true) -> [status_code:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: Write() -> [{"traffic":"egress","route_name":"rate-limit-route","method":"GET","host":"www.twitter.com","path":"","protocol":"HTTP/1.1","status_code":503,"status_flags":"RL","bytes_received":-1,"bytes_sent":0,"timeout_ms":-1,"rate-limit":2000,"rate-burst":0,"retry":,"retry-rate-limit":-1,"retry-rate-burst":-1,"failover":}]
	//test: RoundTrip(handler:true) -> [status_code:503] [err:<nil>]

}

func Example_Controller_Default_Retry_NotEnabled() {
	req, _ := http.NewRequest("GET", facebookUrl, nil)

	if !isEnabled2 {
		isEnabled2 = true
		ControllerWrapTransport(nil)
	}
	act := controller.EgressTable.LookupByName(retryRoute)
	if act != nil {
		if c, ok := act.Retry(); ok {
			c.Disable()
		}
	}
	resp, err := http.DefaultClient.Do(req)
	fmt.Printf("test: RoundTrip(handler:true) -> [status_code:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: Write() -> [{"traffic":"egress","route_name":"retry-route","method":"GET","host":"www.facebook.com","path":"","protocol":"HTTP/1.1","status_code":504,"status_flags":"NE","bytes_received":-1,"bytes_sent":0,"timeout_ms":1,"rate-limit":-1,"rate-burst":-1,"retry":false,"retry-rate-limit":0,"retry-rate-burst":0,"failover":}]
	//test: RoundTrip(handler:true) -> [status_code:504] [err:<nil>]

}

func Example_Controller_Default_Retry_RateLimited() {
	req, _ := http.NewRequest("GET", facebookUrl, nil)

	if !isEnabled2 {
		isEnabled2 = true
		ControllerWrapTransport(nil)
	}
	act := controller.EgressTable.LookupByName(retryRoute)
	if act != nil {
		if c, ok := act.Retry(); ok {
			c.Enable()
		}
	}
	resp, err := http.DefaultClient.Do(req)
	fmt.Printf("test: RoundTrip(handler:true) -> [status_code:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: Write() -> [{"traffic":"egress","route_name":"retry-route","method":"GET","host":"www.facebook.com","path":"","protocol":"HTTP/1.1","status_code":504,"status_flags":"RL","bytes_received":-1,"bytes_sent":0,"timeout_ms":1,"rate-limit":-1,"rate-burst":-1,"retry":false,"retry-rate-limit":0,"retry-rate-burst":0,"failover":}]
	//test: RoundTrip(handler:true) -> [status_code:504] [err:<nil>]

}

func Example_Controller_Default_Retry() {
	req, _ := http.NewRequest("GET", facebookUrl, nil)

	if !isEnabled2 {
		isEnabled2 = true
		ControllerWrapTransport(nil)
	}
	act := controller.EgressTable.LookupByName(retryRoute)
	if act != nil {
		if c, ok := act.Retry(); ok {
			c.Enable()
		}
		if c, ok := act.Retry(); ok {
			c.SetRateLimiter(100, 10)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	fmt.Printf("test: RoundTrip(handler:true) -> [status_code:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: Write() -> [{"traffic":"egress","route_name":"retry-route","method":"GET","host":"www.facebook.com","path":"","protocol":"HTTP/1.1","status_code":504,"status_flags":"UT","bytes_received":-1,"bytes_sent":0,"timeout_ms":1,"rate-limit":-1,"rate-burst":-1,"retry":false,"retry-rate-limit":100,"retry-rate-burst":10,"failover":}]
	//test: Write() -> [{"traffic":"egress","route_name":"retry-route","method":"GET","host":"www.facebook.com","path":"","protocol":"HTTP/1.1","status_code":504,"status_flags":"UT","bytes_received":-1,"bytes_sent":0,"timeout_ms":1,"rate-limit":-1,"rate-burst":-1,"retry":true,"retry-rate-limit":100,"retry-rate-burst":10,"failover":}]
	//test: RoundTrip(handler:true) -> [status_code:504] [err:<nil>]

}
