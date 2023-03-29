package controller

import (
	"fmt"
	"golang.org/x/time/rate"
)

func Example_newRateLimiter() {
	t := newRateLimiter("test-route", newTable(true, false), NewRateLimiterConfig(true, 503, 1, 100))
	limit, burst := t.LimitAndBurst()
	fmt.Printf("test: newRateLimiter() -> [name:%v] [limit:%v] [burst:%v] [statusCode:%v]\n", t.name, limit, burst, t.StatusCode())

	t = newRateLimiter("test-route2", newTable(true, false), NewRateLimiterConfig(true, 429, rate.Inf, DefaultBurst))
	limit, burst = t.LimitAndBurst()
	fmt.Printf("test: newRateLimiter() -> [name:%v] [limit:%v] [burst:%v] [statusCode:%v]\n", t.name, limit, burst, t.StatusCode())

	t2 := cloneRateLimiter(t)
	t2.config.Limit = 123
	fmt.Printf("test: cloneRateLimiter() -> [prev-limit:%v] [prev-name:%v] [curr-limit:%v] [curr-name:%v]\n", t.config.Limit, t.name, t2.config.Limit, t2.name)

	//Output:
	//test: newRateLimiter() -> [name:test-route] [limit:1] [burst:100] [statusCode:503]
	//test: newRateLimiter() -> [name:test-route2] [limit:1.7976931348623157e+308] [burst:1] [statusCode:429]
	//test: cloneRateLimiter() -> [prev-limit:1.7976931348623157e+308] [prev-name:test-route2] [curr-limit:123] [curr-name:test-route2]

}

/*
func Example_RateLimiter_Status() {
	name := "test-route"
	config := NewRateLimiterConfig(10, 100, 503)
	t := newTable(true)
	err := t.Add(name, config)
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", err, t.count())

	act := t.LookupByName(name)
	fmt.Printf("test: IsEnabled() -> [%v]\n", act.RateLimiter().IsEnabled())
	prevEnabled := act.RateLimiter().IsEnabled()

	act.RateLimiter().Disable()
	act1 := t.LookupByName(name)
	fmt.Printf("test: Disable() -> [prev-enabled:%v] [curr-enabled:%v]\n", prevEnabled, act1.RateLimiter().IsEnabled())
	prevEnabled = act1.RateLimiter().IsEnabled()

	act1.RateLimiter().Enable()
	act = t.LookupByName(name)
	fmt.Printf("test: Enable() -> [prev-enabled:%v] [curr-enabled:%v]\n", prevEnabled, act.RateLimiter().IsEnabled())
	prevEnabled = act.RateLimiter().IsEnabled()

	//Output:
	//test: Add() -> [<nil>] [count:1]
	//test: IsEnabled() -> [true]
	//test: Disable() -> [prev-enabled:true] [curr-enabled:false]
	//test: Enable() -> [prev-enabled:false] [curr-enabled:true]

}
*/

func Example_RateLimiter_Set() {
	name := "test-route"
	config := NewRateLimiterConfig(true, 503, 10, 100)
	t := newTable(true, false)
	err := t.AddController(newRoute(name, config))
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", err, t.count())

	act := t.LookupByName(name)
	fmt.Printf("test: rateLimiterState(map,t) -> %v\n", rateLimiterState(nil, act.t().rateLimiter))

	act.t().rateLimiter.SetLimit(rate.Inf)
	act1 := t.LookupByName(name)
	fmt.Printf("test: SetLimit(rate.Inf) -> %v\n", rateLimiterState(nil, act1.t().rateLimiter))

	act1.t().rateLimiter.SetBurst(1)
	act = t.LookupByName(name)
	fmt.Printf("test: SetBurst(1) -> %v\n", rateLimiterState(nil, act.t().rateLimiter))

	//Output:
	//test: Add() -> [[]] [count:1]
	//test: rateLimiterState(map,t) -> map[burst:100 rateLimit:10]
	//test: SetLimit(rate.Inf) -> map[burst:100 rateLimit:99999]
	//test: SetBurst(1) -> map[burst:1 rateLimit:99999]

}

func Example_RateLimiter_Adjust() {
	name := "test-route"
	config := NewRateLimiterConfig(true, 503, 10, 1)
	t := newTable(true, false)
	err := t.AddController(newRoute(name, config))
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", err, t.count())

	act := t.LookupByName(name)
	fmt.Printf("test: rateLimiterState(map,t) -> %v\n", rateLimiterState(nil, act.t().rateLimiter))

	ok := act.t().rateLimiter.AdjustRateLimiter(10)
	act1 := t.LookupByName(name)
	fmt.Printf("test: AdjustRateLimiter(10) -> [%v] [state:%v]\n", ok, rateLimiterState(nil, act1.t().rateLimiter))

	ok = act1.t().rateLimiter.AdjustRateLimiter(1)
	act = t.LookupByName(name)
	fmt.Printf("test: AdjustRateLimiter(1) -> [%v] [state:%v]\n", ok, rateLimiterState(nil, act.t().rateLimiter))

	act1.t().rateLimiter.SetRateLimiter(100, 25)
	act = t.LookupByName(name)
	fmt.Printf("test: rateLimiterState(map,t) -> %v\n", rateLimiterState(nil, act.t().rateLimiter))

	ok = act.t().rateLimiter.AdjustRateLimiter(10)
	act1 = t.LookupByName(name)
	fmt.Printf("test: AdjustRateLimiter(10) -> [%v] [state:%v]\n", ok, rateLimiterState(nil, act1.t().rateLimiter))

	ok = act1.t().rateLimiter.AdjustRateLimiter(-10)
	act = t.LookupByName(name)
	fmt.Printf("test: AdjustRateLimiter(-10) -> [%v] [state:%v]\n", ok, rateLimiterState(nil, act.t().rateLimiter))

	//fmt.Printf("test: AdjustRateLimiter(10) -> [%v] [state:%v]\n", ok, rateLimiterState(nil, act1.t().rateLimiter))

	//Output:
	//test: Add() -> [[]] [count:1]
	//test: rateLimiterState(map,t) -> map[burst:1 rateLimit:10]
	//test: AdjustRateLimiter(10) -> [false] [state:map[burst:1 rateLimit:10]]
	//test: AdjustRateLimiter(1) -> [false] [state:map[burst:1 rateLimit:10]]
	//test: rateLimiterState(map,t) -> map[burst:25 rateLimit:100]
	//test: AdjustRateLimiter(10) -> [true] [state:map[burst:28 rateLimit:110]]
	//test: AdjustRateLimiter(-10) -> [true] [state:map[burst:25 rateLimit:99]]

}
