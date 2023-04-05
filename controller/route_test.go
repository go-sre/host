package controller

import (
	"encoding/json"
	"fmt"
	"time"
)

func ExampleNewRoute() {
	name := "nil-config"
	route := newRoute(name)
	fmt.Printf("test: newRoute() -> [name:%v] [timeout:%v] [rateLimiter:%v] [retry:%v]\n", name,
		route.Timeout != nil, route.RateLimiter != nil, route.Retry != nil)

	name = "timeout"
	route = newRoute(name, NewTimeoutConfig(true, 504, time.Second*2))
	fmt.Printf("test: newRoute() -> [name:%v] [timeout:%v] [rateLimiter:%v] [retry:%v]\n", name,
		route.Timeout != nil, route.RateLimiter != nil, route.Retry != nil)

	name = "timeout-rateLimiter"
	route = newRoute(name, NewTimeoutConfig(true, 504, time.Second*2), NewRateLimiterConfig(true, 503, 100, 25))
	fmt.Printf("test: newRoute() -> [name:%v] [timeout:%v] [rateLimiter:%v] [retry:%v]\n", name,
		route.Timeout != nil, route.RateLimiter != nil, route.Retry != nil)

	name = "timeout-rateLimiter-retry"
	route = newRoute(name, NewTimeoutConfig(true, 504, time.Second*2), NewRateLimiterConfig(true, 503, 100, 25), NewRetryConfig(false, 100, 25, time.Second, []int{504, 503}))
	fmt.Printf("test: newRoute() -> [name:%v] [timeout:%v] [rateLimiter:%v] [retry:%v]\n", name,
		route.Timeout != nil, route.RateLimiter != nil, route.Retry != nil)

	name = "timeout-rateLimiter-nil"
	route = newRoute(name, nil, NewTimeoutConfig(true, 504, time.Second*2), nil, NewRateLimiterConfig(true, 503, 100, 25), nil)
	fmt.Printf("test: newRoute() -> [name:%v] [timeout:%v] [rateLimiter:%v] [retry:%v]\n", name,
		route.Timeout != nil, route.RateLimiter != nil, route.Retry != nil)

	//Output:
	//test: newRoute() -> [name:nil-config] [timeout:false] [rateLimiter:false] [retry:false]
	//test: newRoute() -> [name:timeout] [timeout:true] [rateLimiter:false] [retry:false]
	//test: newRoute() -> [name:timeout-rateLimiter] [timeout:true] [rateLimiter:true] [retry:false]
	//test: newRoute() -> [name:timeout-rateLimiter-retry] [timeout:true] [rateLimiter:true] [retry:true]
	//test: newRoute() -> [name:timeout-rateLimiter-nil] [timeout:true] [rateLimiter:true] [retry:false]

}

func ExampleConfig_Marshal() {
	config := Route{Name: "test-route", Pattern: "google.com", Traffic: IngressTraffic, Protocol: "HTTP11", Ping: true,
		Timeout: &TimeoutConfig{
			StatusCode: 504,
			Duration:   20000,
		},
		RateLimiter: &RateLimiterConfig{
			Limit:      100,
			Burst:      25,
			StatusCode: 503,
		},
		Retry: &RetryConfig{
			Limit:       100,
			Burst:       33,
			Wait:        500,
			StatusCodes: []int{503, 504},
		},
		Proxy: &ProxyConfig{
			Enabled: false,
			Pattern: "http:",
		},
	}
	buf, err := json.Marshal(config)
	fmt.Printf("test: Config{} -> [error:%v] %v\n", err, string(buf))

	//Output:
	//test: Config{} -> [error:<nil>] {"Name":"test-route","Pattern":"google.com","Traffic":"ingress","Ping":true,"Protocol":"HTTP11","Timeout":{"Enabled":false,"StatusCode":504,"Duration":20000},"RateLimiter":{"Enabled":false,"StatusCode":503,"Limit":100,"Burst":25},"Retry":{"Enabled":false,"Limit":100,"Burst":33,"Wait":500,"StatusCodes":[503,504]},"Proxy":{"Enabled":false,"Pattern":"http:","Headers":null,"Action":null}}

}

func ExampleNewRouteFromConfig() {
	config := RouteConfig{
		Name:    "test-route",
		Pattern: "/health/liveness",
		Timeout: &TimeoutConfigJson{
			Enabled:    true,
			Duration:   "500ms",
			StatusCode: 5040,
		},
		RateLimiter: nil,
		Retry: &RetryConfigJson{
			Limit:       100,
			Burst:       25,
			Wait:        "5x",
			StatusCodes: nil,
		},
	}
	route, err := NewRouteFromConfig(config)
	fmt.Printf("test: NewRouteFromConfig() [err:%v] [route:%v]\n", err, route)

	config.Retry.Wait = "245s"
	route, err = NewRouteFromConfig(config)
	fmt.Printf("test: NewRouteFromConfig() [err:%v] [timeout:%v] [retry:%v]\n", err, route.Timeout, route.Retry)

	config.Timeout.Duration = "x34"
	route, err = NewRouteFromConfig(config)
	fmt.Printf("test: NewRouteFromConfig() [err:%v] [route:%v]\n", err, route)

	//Output:
	//test: NewRouteFromConfig() [err:strconv.Atoi: parsing "5x": invalid syntax] [route:{   false  <nil> <nil> <nil> <nil>}]
	//test: NewRouteFromConfig() [err:<nil>] [timeout:&{true 5040 500ms}] [retry:&{false 100 25 4m5s []}]
	//test: NewRouteFromConfig() [err:strconv.Atoi: parsing "x34": invalid syntax] [route:{   false  <nil> <nil> <nil> <nil>}]

}

func _ExampleConfig_Unmarshal() {
	var config = Route{}
	s := "{\"Name\":\"test-route\",\"Timeout\":{\"StatusCode\":504,\"Timeout\":20000},\"RateLimiter\":{\"Limit\":100,\"Burst\":25,\"StatusCode\":503},\"Retry\":{\"Limit\":100,\"Burst\":33,\"Wait\":500,\"Codes\":[503,504]}}"

	err := json.Unmarshal([]byte(s), &config)

	//buf, err := json.Marshal(config)
	fmt.Printf("test: Config{} -> [error:%v] [%v]\n", err, config)

	//Output:
	//test: Config{} -> [error:<nil>] [{test-route {504 20µs} {100 25 503} {100 33 500ns [503 504]} {false <nil>}}]
}
