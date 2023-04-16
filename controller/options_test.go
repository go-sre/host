package controller

import (
	"net/http"
	"time"
)

func ExampleLog() {
	start := time.Now().UTC()
	time.Sleep(time.Second)

	req, _ := http.NewRequest("GET", "http://www.google.com/search?t=test", nil)
	req.Header.Add(RequestIdHeaderName, "1234-56-7890")

	resp := new(http.Response)
	resp.StatusCode = 404

	defaultLogFn("egress", start, time.Since(start), "test-route", req, resp, 500, 100, 10, "true", "UT")

	//Output:
	//{traffic:egress ,route:test-route ,request-id:1234-56-7890, status-code:404, method:GET, url:http://www.google.com/search?t=test, host:www.google.com, path:/search, timeout:500, rate-limit:100, rate-burst:10, proxy:true, status-flags:UT}

}
