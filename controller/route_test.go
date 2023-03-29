package controller

import (
	"encoding/json"
	"fmt"
	"time"
)

func ExampleNewRoute() {
	name := "nil-config"
	route := newRoute(name)
	fmt.Printf("test: newRoute() -> [name:%v] [timeout:%v] [rateLimiter:%v] [retry:%v] [failover:%v]\n", name,
		route.Timeout != nil, route.RateLimiter != nil, route.Retry != nil, route.Failover != nil)

	name = "timeout"
	route = newRoute(name, NewTimeoutConfig(time.Second*2, 504))
	fmt.Printf("test: newRoute() -> [name:%v] [timeout:%v] [rateLimiter:%v] [retry:%v] [failover:%v]\n", name,
		route.Timeout != nil, route.RateLimiter != nil, route.Retry != nil, route.Failover != nil)

	name = "timeout-rateLimiter"
	route = newRoute(name, NewTimeoutConfig(time.Second*2, 504), NewRateLimiterConfig(100, 25, 503))
	fmt.Printf("test: newRoute() -> [name:%v] [timeout:%v] [rateLimiter:%v] [retry:%v] [failover:%v]\n", name,
		route.Timeout != nil, route.RateLimiter != nil, route.Retry != nil, route.Failover != nil)

	name = "timeout-rateLimiter-retry"
	route = newRoute(name, NewTimeoutConfig(time.Second*2, 504), NewRateLimiterConfig(100, 25, 503), NewRetryConfig([]int{504, 503}, 100, 25, time.Second))
	fmt.Printf("test: newRoute() -> [name:%v] [timeout:%v] [rateLimiter:%v] [retry:%v] [failover:%v]\n", name,
		route.Timeout != nil, route.RateLimiter != nil, route.Retry != nil, route.Failover != nil)

	name = "timeout-rateLimiter-retry-failover"
	route = newRoute(name, NewTimeoutConfig(time.Second*2, 504), NewRateLimiterConfig(100, 25, 503), NewRetryConfig([]int{504, 503}, 100, 25, time.Second), NewFailoverConfig(nil))
	fmt.Printf("test: newRoute() -> [name:%v] [timeout:%v] [rateLimiter:%v] [retry:%v] [failover:%v]\n", name,
		route.Timeout != nil, route.RateLimiter != nil, route.Retry != nil, route.Failover != nil)

	name = "timeout-rateLimiter-nil"
	route = newRoute(name, nil, NewTimeoutConfig(time.Second*2, 504), nil, NewRateLimiterConfig(100, 25, 503), nil)
	fmt.Printf("test: newRoute() -> [name:%v] [timeout:%v] [rateLimiter:%v] [retry:%v] [failover:%v]\n", name,
		route.Timeout != nil, route.RateLimiter != nil, route.Retry != nil, route.Failover != nil)

	//Output:
	//test: newRoute() -> [name:nil-config] [timeout:false] [rateLimiter:false] [retry:false] [failover:false]
	//test: newRoute() -> [name:timeout] [timeout:true] [rateLimiter:false] [retry:false] [failover:false]
	//test: newRoute() -> [name:timeout-rateLimiter] [timeout:true] [rateLimiter:true] [retry:false] [failover:false]
	//test: newRoute() -> [name:timeout-rateLimiter-retry] [timeout:true] [rateLimiter:true] [retry:true] [failover:false]
	//test: newRoute() -> [name:timeout-rateLimiter-retry-failover] [timeout:true] [rateLimiter:true] [retry:true] [failover:true]
	//test: newRoute() -> [name:timeout-rateLimiter-nil] [timeout:true] [rateLimiter:true] [retry:false] [failover:false]

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
			Limit: 100,
			Burst: 33,
			Wait:  500,
			Codes: []int{503, 504},
		},
		//Failover: &FailoverConfig{
		//	Enabled: false,
		//	invoke:  nil,
		//},
		Proxy: &ProxyConfig{
			Enabled: false,
			Pattern: "http:",
		},
	}
	buf, err := json.Marshal(config)
	fmt.Printf("test: Config{} -> [error:%v] %v\n", err, string(buf))

	//list := []Route{config, config} //{Package: "package-one", Config: config}, {Package: "package-two", Config: config}}
	//buf, err = json.Marshal(list)
	//fmt.Printf("test: []Route -> [error:%v] %v\n", err, string(buf))

	//Output:
	//test: Config{} -> [error:<nil>] {"Name":"test-route","Pattern":"google.com","Traffic":"ingress","Ping":true,"Protocol":"HTTP11","Timeout":{"Duration":20000,"StatusCode":504},"RateLimiter":{"Limit":100,"Burst":25,"StatusCode":503},"Retry":{"Limit":100,"Burst":33,"Wait":500,"Codes":[503,504]},"Failover":null,"Proxy":{"Enabled":false,"Pattern":"http:","Headers":null}}
	
}

func ExampleNewRouteFromConfig() {
	config := RouteConfig{
		Name:    "test-route",
		Pattern: "/health/liveness",
		Timeout: &TimeoutConfigJson{
			Duration:   "500ms",
			StatusCode: 5040,
		},
		RateLimiter: nil,
		Retry: &RetryConfigJson{
			Limit: 100,
			Burst: 25,
			Wait:  "5x",
			Codes: nil,
		},
		Failover: nil,
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
	//test: NewRouteFromConfig() [err:strconv.Atoi: parsing "5x": invalid syntax] [route:{   false  <nil> <nil> <nil> <nil> <nil>}]
	//test: NewRouteFromConfig() [err:<nil>] [timeout:&{500ms 5040}] [retry:&{100 25 4m5s []}]
	//test: NewRouteFromConfig() [err:strconv.Atoi: parsing "x34": invalid syntax] [route:{   false  <nil> <nil> <nil> <nil> <nil>}]

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
