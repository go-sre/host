package middleware

import (
	"fmt"
	"github.com/go-sre/host/controller"
	"io"
	"net/http"
	"net/http/httptest"
)

func init() {

}

func ExampleSignalHandler_InvalidArgument() {
	req, _ := http.NewRequest("GET", "localhost:8080/signal", nil)
	record := httptest.NewRecorder()
	SignalHandler(record, req)
	resp := record.Result()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("test: SignalHandler(nil) -> [statusCode:%v] [body:%v]\n", resp.StatusCode, string(body))

	req, _ = http.NewRequest("GET", "localhost:8080/signal?traffic=egress", nil)
	record = httptest.NewRecorder()
	SignalHandler(record, req)
	resp = record.Result()
	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("test: SignalHandler(traffic=egress) -> [statusCode:%v] [body:%v]\n", resp.StatusCode, string(body))

	req, _ = http.NewRequest("GET", "localhost:8080/signal?traffic=egress&route=timeout-route", nil)
	record = httptest.NewRecorder()
	SignalHandler(record, req)
	resp = record.Result()
	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("test: SignalHandler(traffic=egress&route=timeout-route) -> [statusCode:%v] [body:%v]\n", resp.StatusCode, string(body))

	req, _ = http.NewRequest("GET", "localhost:8080/signal?traffic=egress&route=timeout-route&behavior=proxy", nil)
	record = httptest.NewRecorder()
	SignalHandler(record, req)
	resp = record.Result()
	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("test: SignalHandler(traffic=egress&route=timeout-route&behavior=proxy) -> [statusCode:%v] [body:%v]\n", resp.StatusCode, string(body))

	//Output:
	//test: SignalHandler(nil) -> [statusCode:400] [body:invalid argument: route [] not found in [] table]
	//test: SignalHandler(traffic=egress) -> [statusCode:400] [body:invalid argument: route [] not found in [egress] table]
	//test: SignalHandler(traffic=egress&route=timeout-route) -> [statusCode:400] [body:invalid argument: behavior [] is not supported]
	//test: SignalHandler(traffic=egress&route=timeout-route&behavior=proxy) -> [statusCode:400] [body:invalid signal: proxy is not configured]

}

func ExampleSignalHandler_Timeout() {
	ctrl := controller.EgressTable().LookupByName(timeoutRoute)
	fmt.Printf("test: TimeoutController() -> [enabled:%v] [duration:%v]\n", ctrl.Timeout().IsEnabled(), ctrl.Timeout().Duration())

	req, _ := http.NewRequest("GET", "localhost:8080/signal?traffic=egress&route=timeout-route&behavior=timeout&enable=false&duration=2s", nil)
	record := httptest.NewRecorder()
	SignalHandler(record, req)
	resp := record.Result()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("test: SignalHandler(disabled,2s) -> [statusCode:%v] [body:%v]\n", resp.StatusCode, string(body))

	ctrl = controller.EgressTable().LookupByName(timeoutRoute)
	fmt.Printf("test: TimeoutController() -> [enabled:%v] [duration:%v]\n", ctrl.Timeout().IsEnabled(), ctrl.Timeout().Duration())

	//Output:
	//test: TimeoutController() -> [enabled:true] [duration:1ms]
	//test: SignalHandler(disabled,2s) -> [statusCode:200] [body:]
	//test: TimeoutController() -> [enabled:false] [duration:2s]

}

func ExampleSignalHandler_RateLimit() {
	ctrl := controller.EgressTable().LookupByName(rateLimitRoute)
	fmt.Printf("test: RateLimitController() -> [enabled:%v] [limit:%v] [burst:%v]\n", ctrl.RateLimiter().IsEnabled(), ctrl.RateLimiter().Limit(), ctrl.RateLimiter().Burst())

	req, _ := http.NewRequest("GET", "localhost:8080/signal?traffic=egress&route=rate-limit-route&behavior=rate-limit&enable=false&limit=45&burst=5", nil)
	record := httptest.NewRecorder()
	SignalHandler(record, req)
	resp := record.Result()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("test: SignalHandler(disabled,45,5) -> [statusCode:%v] [body:%v]\n", resp.StatusCode, string(body))

	ctrl = controller.EgressTable().LookupByName(rateLimitRoute)
	fmt.Printf("test: RateLimitController() -> [enabled:%v] [limit:%v] [burst:%v]\n", ctrl.RateLimiter().IsEnabled(), ctrl.RateLimiter().Limit(), ctrl.RateLimiter().Burst())

	//Output:
	//test: RateLimitController() -> [enabled:true] [limit:2000] [burst:10]
	//test: SignalHandler(disabled,45,5) -> [statusCode:200] [body:]
	//test: RateLimitController() -> [enabled:false] [limit:45] [burst:5]

}

func ExampleSignalHandler_Retry() {
	ctrl := controller.EgressTable().LookupByName(retryRoute)
	fmt.Printf("test: RetryController() -> [enabled:%v] [limit:%v] [burst:%v] [wait:%v]\n", ctrl.Retry().IsEnabled(), ctrl.Retry().Limit(), ctrl.Retry().Burst(), ctrl.Retry().Wait())

	req, _ := http.NewRequest("GET", "localhost:8080/signal?traffic=egress&route=retry-route&behavior=retry&enable=false&limit=45&burst=5&wait=100ms", nil)
	record := httptest.NewRecorder()
	SignalHandler(record, req)
	resp := record.Result()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("test: SignalHandler(disabled,45,5,100ms) -> [statusCode:%v] [body:%v]\n", resp.StatusCode, string(body))

	ctrl = controller.EgressTable().LookupByName(retryRoute)
	fmt.Printf("test: RetryController() -> [enabled:%v] [limit:%v] [burst:%v] [wait:%v]\n", ctrl.Retry().IsEnabled(), ctrl.Retry().Limit(), ctrl.Retry().Burst(), ctrl.Retry().Wait())

	//Output:
	//test: RetryController() -> [enabled:true] [limit:0] [burst:0] [wait:0s]
	//test: SignalHandler(disabled,45,5,100ms) -> [statusCode:200] [body:]
	//test: RetryController() -> [enabled:false] [limit:45] [burst:5] [wait:100ms]

}

func ExampleSignalHandler_Proxy() {
	ctrl := controller.EgressTable().LookupByName(proxyRoute)
	fmt.Printf("test: ProxyController() -> [enabled:%v] [pattern:%v]\n", ctrl.Proxy().IsEnabled(), ctrl.Proxy().Pattern())

	req, _ := http.NewRequest("GET", "localhost:8080/signal?traffic=egress&route=proxy-route&behavior=proxy&enable=false&pattern=http://localhost:8080", nil)
	record := httptest.NewRecorder()
	SignalHandler(record, req)
	resp := record.Result()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("test: SignalHandler(disabled,http://localhost:8080) -> [statusCode:%v] [body:%v]\n", resp.StatusCode, string(body))

	ctrl = controller.EgressTable().LookupByName(proxyRoute)
	fmt.Printf("test: ProxyController() -> [enabled:%v] [pattern:%v]\n", ctrl.Proxy().IsEnabled(), ctrl.Proxy().Pattern())

	//Output:
	//test: ProxyController() -> [enabled:true] [pattern:https://www.google.com/search?q=test]
	//test: SignalHandler(disabled,http://localhost:8080) -> [statusCode:200] [body:]
	//test: ProxyController() -> [enabled:false] [pattern:http://localhost:8080]

}
