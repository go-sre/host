package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gotemplates/host/shared"
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
	config := Route{Name: "test-route", Pattern: "google.com", Traffic: shared.IngressTraffic, Protocol: "HTTP11", Ping: true,
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
	}
	buf, err := json.Marshal(config)
	fmt.Printf("test: Config{} -> [error:%v] %v\n", err, string(buf))

	//list := []Route{config, config} //{Package: "package-one", Config: config}, {Package: "package-two", Config: config}}

	//buf, err = json.Marshal(list)
	//fmt.Printf("test: []Route -> [error:%v] %v\n", err, string(buf))

	//Output:
	//test: Config{} -> [error:<nil>] {"Name":"test-route","Pattern":"google.com","Traffic":"ingress","Ping":true,"Protocol":"HTTP11","Timeout":{"Duration":20000,"StatusCode":504},"RateLimiter":{"Limit":100,"Burst":25,"StatusCode":503},"Retry":{"Limit":100,"Burst":33,"Wait":500,"Codes":[503,504]},"Failover":null}

}

func ExampleConvertDuration() {
	s := ""
	duration, err := ConvertDuration(s)
	fmt.Printf("test: ConvertDuration(\"%v\") [err:%v] [duration:%v]\n", s, err, duration)

	s = "  "
	duration, err = ConvertDuration(s)
	fmt.Printf("test: ConvertDuration(\"%v\") [err:%v] [duration:%v]\n", s, err, duration)

	s = "12as"
	duration, err = ConvertDuration(s)
	fmt.Printf("test: ConvertDuration(\"%v\") [err:%v] [duration:%v]\n", s, err, duration)

	s = "1000"
	duration, err = ConvertDuration(s)
	fmt.Printf("test: ConvertDuration(\"%v\") [err:%v] [duration:%v]\n", s, err, duration)

	s = "1000s"
	duration, err = ConvertDuration(s)
	fmt.Printf("test: ConvertDuration(\"%v\") [err:%v] [duration:%v]\n", s, err, duration)

	s = "1000m"
	duration, err = ConvertDuration(s)
	fmt.Printf("test: ConvertDuration(\"%v\") [err:%v] [duration:%v]\n", s, err, duration)

	s = "1m"
	duration, err = ConvertDuration(s)
	fmt.Printf("test: ConvertDuration(\"%v\") [err:%v] [duration:%v]\n", s, err, duration)

	s = "10ms"
	duration, err = ConvertDuration(s)
	fmt.Printf("test: ConvertDuration(\"%v\") [err:%v] [duration:%v]\n", s, err, duration)

	//t := time.Microsecond * 100
	//fmt.Printf("test: time.String %v\n", t.String())

	s = "10µs"
	duration, err = ConvertDuration(s)
	fmt.Printf("test: ConvertDuration(\"%v\") [err:%v] [duration:%v]\n", s, err, duration)

	//Output:
	//test: ConvertDuration("") [err:<nil>] [duration:0s]
	//test: ConvertDuration("  ") [err:strconv.Atoi: parsing "  ": invalid syntax] [duration:0s]
	//test: ConvertDuration("12as") [err:strconv.Atoi: parsing "12a": invalid syntax] [duration:0s]
	//test: ConvertDuration("1000") [err:<nil>] [duration:16m40s]
	//test: ConvertDuration("1000s") [err:<nil>] [duration:16m40s]
	//test: ConvertDuration("1000m") [err:<nil>] [duration:16h40m0s]
	//test: ConvertDuration("1m") [err:<nil>] [duration:1m0s]
	//test: ConvertDuration("10ms") [err:<nil>] [duration:10ms]
	//test: ConvertDuration("10µs") [err:<nil>] [duration:10µs]

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
	//test: NewRouteFromConfig() [err:strconv.Atoi: parsing "5x": invalid syntax] [route:{   false  <nil> <nil> <nil> <nil>}]
	//test: NewRouteFromConfig() [err:<nil>] [timeout:&{500ms 5040}] [retry:&{100 25 4m5s []}]
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
