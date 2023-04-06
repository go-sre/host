package controller

import (
	"net/http"
	"time"
)

func _ExampleLog() {
	start := time.Now().UTC()
	time.Sleep(time.Second)

	req, _ := http.NewRequest("GET", "http://www.google.com/search?t=test", nil)
	req.Header.Add(RequestIdHeaderName, "1234-56-7890")

	resp := new(http.Response)
	resp.StatusCode = 404
	state := make(map[string]string)
	state[ControllerName] = "test-route"
	state[TimeoutName] = "500"

	state[RateLimitName] = "100"
	state[RateBurstName] = "10"

	state[RetryName] = "true"
	state[RetryRateLimitName] = "10"
	state[RetryRateBurstName] = "1"

	state[ProxyName] = "true"

	defaultLogFn("egress", start, time.Since(start), req, resp, "UT", state)

	//Output:
	//{start:2023-02-25 14:57:37.040782 ,duration:1013 ,traffic:egress, route:test-route, request-id:1234-56-7890, protocol:HTTP/1.1, method:GET, url:http://www.google.com/search?t=test, host:www.google.com, path:/search, status-code:404, timeout_ms:500, rate-limit:100, rate-burst:10, retry:true, retry-rate-limit:10, retry-rate-burst:1, status-flags:UT}

}
