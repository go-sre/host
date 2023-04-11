package middleware

import (
	"fmt"
	"github.com/go-sre/host/controller"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
)

func init() {

}

func ExampleActuatorHandler_InvalidArgument() {
	req, _ := http.NewRequest("GET", "http://localhost:8080/actuator", nil)
	record := httptest.NewRecorder()
	ActuatorHandler(record, req)
	resp := record.Result()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("test: ActuatorHandler(nil) -> [statusCode:%v] [body:%v]\n", resp.StatusCode, string(body))

	/*
		req, _ = http.NewRequest("GET", "http://localhost:8080/actuator?enabled=false", nil)
		record = httptest.NewRecorder()
		ActuatorHandler(record, req)
		resp = record.Result()
		body, _ = io.ReadAll(resp.Body)
		fmt.Printf("test: ActuatorHandler(traffic=egress) -> [statusCode:%v] [body:%v]\n", resp.StatusCode, string(body))

		req, _ = http.NewRequest("GET", "http://localhost:8080/signal?traffic=egress&route=timeout-route", nil)
		record = httptest.NewRecorder()
		ActuatorHandler(record, req)
		resp = record.Result()
		body, _ = io.ReadAll(resp.Body)
		fmt.Printf("test: ActuatorHandler(traffic=egress&route=timeout-route) -> [statusCode:%v] [body:%v]\n", resp.StatusCode, string(body))

		req, _ = http.NewRequest("GET", "localhost:8080/signal?traffic=egress&route=timeout-route&behavior=proxy", nil)
		record = httptest.NewRecorder()
		ActuatorHandler(record, req)
		resp = record.Result()
		body, _ = io.ReadAll(resp.Body)
		fmt.Printf("test: ActuatorHandler(traffic=egress&route=timeout-route&behavior=proxy) -> [statusCode:%v] [body:%v]\n", resp.StatusCode, string(body))

	*/

	//Output:
	//test: ActuatorHandler(nil) -> [statusCode:400] [body:invalid argument: request URL does not contain any query arguments]

}

func ExampleActuatorHandler_Timeout() {
	ctrl := controller.EgressTable().LookupByName(timeoutRoute)
	fmt.Printf("test: TimeoutController() -> [enabled:%v] [duration:%v]\n", ctrl.Timeout().IsEnabled(), ctrl.Timeout().Duration())

	req, _ := http.NewRequest("GET", "http://localhost:8080/actuator/egress/timeout-route/timeout?enabled=false&duration=2s", nil)
	record := httptest.NewRecorder()
	ActuatorHandler(record, req)
	resp := record.Result()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("test: ActuatorHandler(disabled,2s) -> [statusCode:%v] [body:%v]\n", resp.StatusCode, string(body))

	ctrl = controller.EgressTable().LookupByName(timeoutRoute)
	fmt.Printf("test: TimeoutController() -> [enabled:%v] [duration:%v]\n", ctrl.Timeout().IsEnabled(), ctrl.Timeout().Duration())

	//Output:
	//test: TimeoutController() -> [enabled:true] [duration:1ms]
	//test: ActuatorHandler(disabled,2s) -> [statusCode:200] [body:]
	//test: TimeoutController() -> [enabled:false] [duration:2s]

}

func ExampleActuatorHandler_RateLimit() {
	ctrl := controller.EgressTable().LookupByName(rateLimitRoute)
	fmt.Printf("test: RateLimitController() -> [enabled:%v] [limit:%v] [burst:%v]\n", ctrl.RateLimiter().IsEnabled(), ctrl.RateLimiter().Limit(), ctrl.RateLimiter().Burst())

	req, _ := http.NewRequest("GET", "http://localhost:8080/actuator/egress/rate-limit-route/rate-limit?enabled=false&limit=45&burst=5", nil)
	record := httptest.NewRecorder()
	ActuatorHandler(record, req)
	resp := record.Result()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("test: ActuatorHandler(disabled,45,5) -> [statusCode:%v] [body:%v]\n", resp.StatusCode, string(body))

	ctrl = controller.EgressTable().LookupByName(rateLimitRoute)
	fmt.Printf("test: RateLimitController() -> [enabled:%v] [limit:%v] [burst:%v]\n", ctrl.RateLimiter().IsEnabled(), ctrl.RateLimiter().Limit(), ctrl.RateLimiter().Burst())

	//Output:
	//test: RateLimitController() -> [enabled:true] [limit:2000] [burst:10]
	//test: ActuatorHandler(disabled,45,5) -> [statusCode:200] [body:]
	//test: RateLimitController() -> [enabled:false] [limit:45] [burst:5]

}

func ExampleActuatorHandler_Retry() {
	ctrl := controller.EgressTable().LookupByName(retryRoute)
	fmt.Printf("test: RetryController() -> [enabled:%v] [limit:%v] [burst:%v] [wait:%v]\n", ctrl.Retry().IsEnabled(), ctrl.Retry().Limit(), ctrl.Retry().Burst(), ctrl.Retry().Wait())

	req, _ := http.NewRequest("GET", "http://localhost:8080/actuator/egress/retry-route/retry?enabled=false&limit=45&burst=5&wait=100ms", nil)
	record := httptest.NewRecorder()
	ActuatorHandler(record, req)
	resp := record.Result()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("test: ActuatorHandler(disabled,45,5,100ms) -> [statusCode:%v] [body:%v]\n", resp.StatusCode, string(body))

	ctrl = controller.EgressTable().LookupByName(retryRoute)
	fmt.Printf("test: RetryController() -> [enabled:%v] [limit:%v] [burst:%v] [wait:%v]\n", ctrl.Retry().IsEnabled(), ctrl.Retry().Limit(), ctrl.Retry().Burst(), ctrl.Retry().Wait())

	//Output:
	//test: RetryController() -> [enabled:true] [limit:0] [burst:0] [wait:0s]
	//test: ActuatorHandler(disabled,45,5,100ms) -> [statusCode:200] [body:]
	//test: RetryController() -> [enabled:false] [limit:45] [burst:5] [wait:100ms]

}

func ExampleActuatorHandler_Proxy() {
	ctrl := controller.EgressTable().LookupByName(proxyRoute)
	fmt.Printf("test: ProxyController() -> [enabled:%v] [pattern:%v]\n", ctrl.Proxy().IsEnabled(), ctrl.Proxy().Pattern())

	req, _ := http.NewRequest("GET", "http://localhost:8080/actuator/egress/proxy-route/proxy?enabled=false&pattern=http://localhost:8080", nil)
	record := httptest.NewRecorder()
	ActuatorHandler(record, req)
	resp := record.Result()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("test: ActuatorHandler(disabled,http://localhost:8080) -> [statusCode:%v] [body:%v]\n", resp.StatusCode, string(body))

	ctrl = controller.EgressTable().LookupByName(proxyRoute)
	fmt.Printf("test: ProxyController() -> [enabled:%v] [pattern:%v]\n", ctrl.Proxy().IsEnabled(), ctrl.Proxy().Pattern())

	//Output:
	//test: ProxyController() -> [enabled:true] [pattern:https://www.google.com/search?q=test]
	//test: ActuatorHandler(disabled,http://localhost:8080) -> [statusCode:200] [body:]
	//test: ProxyController() -> [enabled:false] [pattern:http://localhost:8080]

}

func ExampleActuatorParse_Url() {
	_, _, _, err0 := parseUrl(nil)
	fmt.Printf("test: parseUrl(nil) -> [err:%v]\n", err0)

	url, _ := url.Parse("http://localhost:8080")
	traffic, route, behavior, err := parseUrl(url)
	fmt.Printf("test: parseUrl(%v) -> [t%v] [r:%v] [b:%v] [err:%v]\n", url, traffic, route, behavior, err)

	url, _ = url.Parse("http://localhost:8080/test")
	traffic, route, behavior, err = parseUrl(url)
	fmt.Printf("test: parseUrl(%v) -> [t%v] [r:%v] [b:%v] [err:%v]\n", url, traffic, route, behavior, err)

	url, _ = url.Parse("http://localhost:8080/actuator")
	traffic, route, behavior, err = parseUrl(url)
	fmt.Printf("test: parseUrl(%v) -> [t:%v] [r:%v] [b:%v] [err:%v]\n", url, traffic, route, behavior, err)

	url, _ = url.Parse("http://localhost:8080/actuator/invalid-traffictype")
	traffic, route, behavior, err = parseUrl(url)
	fmt.Printf("test: parseUrl(%v) -> [t:%v] [r:%v] [b:%v] [err:%v]\n", url, traffic, route, behavior, err)

	url, _ = url.Parse("http://localhost:8080/actuator/egress")
	traffic, route, behavior, err = parseUrl(url)
	fmt.Printf("test: parseUrl(%v) -> [t:%v] [r:%v] [b:%v] [err:%v]\n", url, traffic, route, behavior, err)

	url, _ = url.Parse("http://localhost:8080/actuator/egress/test-route")
	traffic, route, behavior, err = parseUrl(url)
	fmt.Printf("test: parseUrl(%v) -> [t:%v] [r:%v] [b:%v] [err:%v]\n", url, traffic, route, behavior, err)

	url, _ = url.Parse("http://localhost:8080/actuator/egress/test-route/timeout")
	traffic, route, behavior, err = parseUrl(url)
	fmt.Printf("test: parseUrl(%v) -> [t:%v] [r:%v] [b:%v] [err:%v]\n", url, traffic, route, behavior, err)

	//Output:
	//test: parseUrl(nil) -> [err:invalid argument: request URL is nil]
	//test: parseUrl(http://localhost:8080) -> [t] [r:] [b:] [err:invalid argument: request URL path is empty]
	//test: parseUrl(http://localhost:8080/test) -> [t] [r:] [b:] [err:invalid argument: request URL path does not start with 'actuator']
	//test: parseUrl(http://localhost:8080/actuator) -> [t:] [r:] [b:] [err:invalid argument: request URL path does not contain traffic type]
	//test: parseUrl(http://localhost:8080/actuator/invalid-traffictype) -> [t:invalid-traffictype] [r:] [b:] [err:invalid argument: request URL path does not contain valid traffic type [invalid-traffictype]]
	//test: parseUrl(http://localhost:8080/actuator/egress) -> [t:egress] [r:] [b:] [err:invalid argument: request URL path does not contain route name]
	//test: parseUrl(http://localhost:8080/actuator/egress/test-route) -> [t:egress] [r:test-route] [b:] [err:invalid argument: request URL path does not contain behavior name]
	//test: parseUrl(http://localhost:8080/actuator/egress/test-route/timeout) -> [t:egress] [r:test-route] [b:timeout] [err:<nil>]

}
