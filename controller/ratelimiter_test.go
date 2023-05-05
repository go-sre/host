package controller

import (
	"fmt"
	"golang.org/x/time/rate"
	"net/url"
	"strconv"
)

func rateLimiterSetValues(limit rate.Limit,
	burst int) url.Values {
	v := make(url.Values)
	if limit != -2 {
		v.Add(RateLimitKey, fmt.Sprintf("%v", limit))
	}
	if burst != -2 {
		v.Add(RateBurstKey, strconv.Itoa(burst))
	}
	return v
}

func Example_newRateLimiter() {
	t := newRateLimiter("test-route", newTable(true, false), NewRateLimiterConfig(true, 503, 1, 100, ""))
	fmt.Printf("test: newRateLimiter() -> [name:%v] [limit:%v] [burst:%v] [statusCode:%v]\n", t.name, t.config.Limit, t.config.Burst, t.StatusCode())

	t = newRateLimiter("test-route2", newTable(true, false), NewRateLimiterConfig(true, 429, rate.Inf, DefaultBurst, ""))
	fmt.Printf("test: newRateLimiter() -> [name:%v] [limit:%v] [burst:%v] [statusCode:%v]\n", t.name, t.config.Limit, t.config.Burst, t.StatusCode())

	t2 := cloneRateLimiter(t)
	t2.config.Limit = 123
	fmt.Printf("test: cloneRateLimiter() -> [prev-limit:%v] [prev-name:%v] [curr-limit:%v] [curr-name:%v]\n", t.config.Limit, t.name, t2.config.Limit, t2.name)

	//Output:
	//test: newRateLimiter() -> [name:test-route] [limit:1] [burst:100] [statusCode:503]
	//test: newRateLimiter() -> [name:test-route2] [limit:1.7976931348623157e+308] [burst:1] [statusCode:429]
	//test: cloneRateLimiter() -> [prev-limit:1.7976931348623157e+308] [prev-name:test-route2] [curr-limit:123] [curr-name:test-route2]

}

func ExampleRateLimiter_State() {
	tbl := newTable(true, false)
	t := newRateLimiter("test-route", tbl, NewRateLimiterConfig(true, 503, 1, 100, "95/200s"))
	fmt.Printf("test: newRateLimiter() -> [name:%v] [limit:%v] [burst:%v] [threshold:%v] [statusCode:%v]\n", t.name, t.config.Limit, t.config.Burst, t.config.Threshold, t.StatusCode())

	limit, burst, threshold := rateLimiterState(t)
	fmt.Printf("test: rateLimiterState(t) -> [enabled:%v] [limit:%v] [burst:%v] [threshold:%v]\n", t.IsEnabled(), limit, burst, threshold)

	t.config.Enabled = false

	limit, burst, threshold = rateLimiterState(t)
	fmt.Printf("test: rateLimiterState(t) -> [enabled:%v] [limit:%v] [burst:%v] [threshold:%v]\n", t.IsEnabled(), limit, burst, threshold)

	//Output:
	//test: newRateLimiter() -> [name:test-route] [limit:1] [burst:100] [threshold:95/200s] [statusCode:503]
	//test: rateLimiterState(t) -> [enabled:true] [limit:1] [burst:100] [threshold:95/200s]
	//test: rateLimiterState(t) -> [enabled:false] [limit:-1] [burst:-1] [threshold:]

}

func ExampleRateLimiter_Toggle() {
	name := "test-route"
	config := NewRateLimiterConfig(true, 503, 100, 10, "")
	t := newTable(true, false)
	err := t.AddController(newRoute(name, config))
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", err, t.count())

	ctrl := t.LookupByName(name)
	fmt.Printf("test: IsEnabled() -> [%v]\n", ctrl.RateLimiter().IsEnabled())
	prevEnabled := ctrl.RateLimiter().IsEnabled()

	ctrl.RateLimiter().Signal(enableValues(false))
	ctrl1 := t.LookupByName(name)
	fmt.Printf("test: Disable() -> [prev-enabled:%v] [curr-enabled:%v]\n", prevEnabled, ctrl1.RateLimiter().IsEnabled())
	prevEnabled = ctrl1.RateLimiter().IsEnabled()

	ctrl1.RateLimiter().Signal(enableValues(true))
	ctrl = t.LookupByName(name)
	fmt.Printf("test: Enable() -> [prev-enabled:%v] [curr-enabled:%v]\n", prevEnabled, ctrl.RateLimiter().IsEnabled())

	//Output:
	//test: Add() -> [[]] [count:1]
	//test: IsEnabled() -> [true]
	//test: Disable() -> [prev-enabled:true] [curr-enabled:false]
	//test: Enable() -> [prev-enabled:false] [curr-enabled:true]

}

func _ExampleRateLimiter_Signal() {
	name := "test-route"
	config := NewRateLimiterConfig(true, 503, 100, 10, "")
	t := newTable(true, false)
	errs := t.AddController(newRoute(name, config))
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", errs, t.count())

	ctrl := t.LookupByName(name)
	limit, burst, _ := rateLimiterState(ctrl.t().rateLimiter)
	fmt.Printf("test: rateLimiterState(map,t) -> [limit:%v] [burst:%v]\n", limit, burst)

	err := ctrl.RateLimiter().Signal(nil)
	fmt.Printf("test: Signal(nil) -> [nil:%v] [empty:%v] \n", ctrl.RateLimiter().Signal(nil), ctrl.RateLimiter().Signal(make(url.Values)))

	err = ctrl.RateLimiter().Signal(rateLimiterSetValues(100, 0))
	ctrl1 := t.LookupByName(name)
	limit, burst, _ = rateLimiterState(ctrl1.t().rateLimiter)
	fmt.Printf("test: Signal(100,0) -> [error:%v] [limit:%v] [burst:%v]\n", err, limit, burst)

	err = ctrl1.RateLimiter().Signal(rateLimiterSetValues(-1, 10))
	ctrl = t.LookupByName(name)
	limit, burst, _ = rateLimiterState(ctrl.t().rateLimiter)
	fmt.Printf("test: Signal(-1,10) -> [error:%v] [limit:%v] [burst:%v]\n", err, limit, burst)

	err = ctrl.RateLimiter().Signal(rateLimiterSetValues(100, 10))
	ctrl1 = t.LookupByName(name)
	limit, burst, _ = rateLimiterState(ctrl1.t().rateLimiter)
	fmt.Printf("test: Signal(100,10) -> [error:%v] [limit:%v] [burst:%v]\n", err, limit, burst)

	err = ctrl1.RateLimiter().Signal(rateLimiterSetValues(100, 8))
	ctrl = t.LookupByName(name)
	limit, burst, _ = rateLimiterState(ctrl.t().rateLimiter)
	fmt.Printf("test: Signal(100,8) -> [error:%v] [limit:%v] [burst:%v]\n", err, limit, burst)

	err = ctrl.RateLimiter().Signal(rateLimiterSetValues(99, 8))
	ctrl1 = t.LookupByName(name)
	limit, burst, _ = rateLimiterState(ctrl1.t().rateLimiter)
	fmt.Printf("test: Signal(99,8) -> [error:%v] [limit:%v] [burst:%v]\n", err, limit, burst)

	err = ctrl1.RateLimiter().Signal(rateLimiterSetValues(-2, 5))
	ctrl = t.LookupByName(name)
	limit, burst, _ = rateLimiterState(ctrl.t().rateLimiter)
	fmt.Printf("test: Signal(99,5) -> [error:%v] [limit:%v] [burst:%v]\n", err, limit, burst)

	err = ctrl.RateLimiter().Signal(rateLimiterSetValues(88, -2))
	ctrl1 = t.LookupByName(name)
	limit, burst, _ = rateLimiterState(ctrl1.t().rateLimiter)
	fmt.Printf("test: Signal(88,5) -> [error:%v] [limit:%v] [burst:%v]\n", err, limit, burst)

	//Output:
	//test: Add() -> [[]] [count:1]
	//test: rateLimiterState(map,t) -> map[burst:10 rateLimit:100]
	//test: Signal(nil) -> [nil:<nil>] [empty:<nil>]
	//test: Signal(100,0) -> [error:invalid argument: burst value is <= 0 [0]] [state:map[burst:10 rateLimit:100]]
	//test: Signal(-1,10) -> [error:invalid argument: limit value is <= 0 [-1]] [state:map[burst:10 rateLimit:100]]
	//test: Signal(100,10) -> [error:<nil>] [state:map[burst:10 rateLimit:100]]
	//test: Signal(100,8) -> [error:<nil>] [state:map[burst:8 rateLimit:100]]
	//test: Signal(99,8) -> [error:<nil>] [state:map[burst:8 rateLimit:99]]
	//test: Signal(99,5) -> [error:<nil>] [state:map[burst:5 rateLimit:99]]
	//test: Signal(88,5) -> [error:<nil>] [state:map[burst:5 rateLimit:88]]

}
