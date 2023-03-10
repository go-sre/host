package controller

import (
	"fmt"
	"golang.org/x/time/rate"
)

func Example_newRetry() {
	t := newRetry("test-route", newTable(true, false), NewRetryConfig([]int{504}, 5, 10, 0))
	limit, burst := t.LimitAndBurst()
	fmt.Printf("test: newRetry() -> [name:%v] [config:%v] [limit:%v] [burst:%v]\n", t.name, t.config, limit, burst)

	t = newRetry("test-route2", newTable(true, false), NewRetryConfig([]int{503, 504}, 2, 20, 0))
	fmt.Printf("test: newRetry() -> [name:%v] [config:%v]\n", t.name, t.config)

	t2 := cloneRetry(t)
	t2.enabled = false
	fmt.Printf("test: cloneRetry() -> [prev-enabled:%v] [curr-enabled:%v]\n", t.enabled, t2.enabled)

	//t = newRetry("test-route3", newTable(true), NewRetryConfig([]int{503, 504}, time.Millisecond*2000, false))
	fmt.Printf("test: retryState(nil,false,map) -> %v\n", retryState(nil, nil, false))

	fmt.Printf("test: retryState(t,false,map) -> %v\n", retryState(nil, t, false))

	t2 = newRetry("test-route", newTable(true, false), NewRetryConfig([]int{504}, rate.Inf, 10, 0))
	fmt.Printf("test: retryState(t2,true,map) -> %v\n", retryState(nil, t2, true))

	//Output:
	//test: newRetry() -> [name:test-route] [config:{5 10 0s [504]}] [limit:5] [burst:10]
	//test: newRetry() -> [name:test-route2] [config:{2 20 0s [503 504]}]
	//test: cloneRetry() -> [prev-enabled:true] [curr-enabled:false]
	//test: retryState(nil,false,map) -> map[retry: retryBurst:-1 retryRateLimit:-1]
	//test: retryState(t,false,map) -> map[retry:false retryBurst:20 retryRateLimit:2]
	//test: retryState(t2,true,map) -> map[retry:true retryBurst:10 retryRateLimit:99999]

}

func Example_Status() {
	name := "test-route"
	config := NewRetryConfig([]int{504}, 5, 10, 0)
	t := newTable(true, false)
	err := t.AddController(newRoute(name, config))
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", err, t.count())

	act := t.LookupByName(name)
	fmt.Printf("test: IsEnabled() -> [%v]\n", act.t().retry.IsEnabled())
	prevEnabled := act.t().retry.IsEnabled()

	act.t().retry.Disable()
	act1 := t.LookupByName(name)
	fmt.Printf("test: Disable() -> [prev-enabled:%v] [curr-enabled:%v]\n", prevEnabled, act1.t().retry.IsEnabled())
	prevEnabled = act1.t().retry.IsEnabled()

	act1.t().retry.Enable()
	act = t.LookupByName(name)
	fmt.Printf("test: Enable() -> [prev-enabled:%v] [curr-enabled:%v]\n", prevEnabled, act.t().retry.IsEnabled())
	prevEnabled = act.t().retry.IsEnabled()

	//Output:
	//test: Add() -> [[]] [count:1]
	//test: IsEnabled() -> [true]
	//test: Disable() -> [prev-enabled:true] [curr-enabled:false]
	//test: Enable() -> [prev-enabled:false] [curr-enabled:true]

}

func Example_IsRetryable_Disabled() {
	name := "test-route"
	config := NewRetryConfig([]int{503, 504}, 100, 10, 0)
	t := newTable(true, false)
	err := t.AddController(newRoute(name, config))
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", err, t.count())

	act := t.LookupByName(name)
	act.t().retry.Disable()
	act = t.LookupByName(name)
	ok, status := act.t().retry.IsRetryable(200)
	fmt.Printf("test: IsRetryable(200) -> [ok:%v] [status:%v]\n", ok, status)

	ok, status = act.t().retry.IsRetryable(503)
	fmt.Printf("test: IsRetryable(503) -> [ok:%v] [status:%v]\n", ok, status)

	ok, status = act.t().retry.IsRetryable(504)
	fmt.Printf("test: IsRetryable(504) -> [ok:%v] [status:%v]\n", ok, status)

	//Output:
	//test: Add() -> [[]] [count:1]
	//test: IsRetryable(200) -> [ok:false] [status:NE]
	//test: IsRetryable(503) -> [ok:false] [status:NE]
	//test: IsRetryable(504) -> [ok:false] [status:NE]

}

func Example_IsRetryable_StatusCode() {
	name := "test-route"
	config := NewRetryConfig([]int{503, 504}, 100, 10, 0)
	t := newTable(true, false)
	err := t.AddController(newRoute(name, config))
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", err, t.count())

	act := t.LookupByName(name)
	act.t().retry.Enable()
	act = t.LookupByName(name)
	ok, status := act.t().retry.IsRetryable(200)
	fmt.Printf("test: IsRetryable(200) -> [ok:%v] [status:%v]\n", ok, status)

	ok, status = act.t().retry.IsRetryable(500)
	fmt.Printf("test: IsRetryable(500) -> [ok:%v] [status:%v]\n", ok, status)

	ok, status = act.t().retry.IsRetryable(502)
	fmt.Printf("test: IsRetryable(502) -> [ok:%v] [status:%v]\n", ok, status)

	ok, status = act.t().retry.IsRetryable(503)
	fmt.Printf("test: IsRetryable(503) -> [ok:%v] [status:%v]\n", ok, status)

	ok, status = act.t().retry.IsRetryable(504)
	fmt.Printf("test: IsRetryable(504) -> [ok:%v] [status:%v]\n", ok, status)

	ok, status = act.t().retry.IsRetryable(505)
	fmt.Printf("test: IsRetryable(505) -> [ok:%v] [status:%v]\n", ok, status)

	//Output:
	//test: Add() -> [[]] [count:1]
	//test: IsRetryable(200) -> [ok:false] [status:]
	//test: IsRetryable(500) -> [ok:false] [status:]
	//test: IsRetryable(502) -> [ok:false] [status:]
	//test: IsRetryable(503) -> [ok:true] [status:]
	//test: IsRetryable(504) -> [ok:true] [status:]
	//test: IsRetryable(505) -> [ok:false] [status:]

}

func Example_IsRetryable_RateLimit() {
	name := "test-route"
	config := NewRetryConfig([]int{503, 504}, 1, 1, 0)
	t := newTable(true, false)
	err := t.AddController(newRoute(name, config))
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", err, t.count())

	act := t.LookupByName(name)
	act.t().retry.Enable()
	act = t.LookupByName(name)
	ok, status := act.t().retry.IsRetryable(503)
	fmt.Printf("test: IsRetryable(503) -> [ok:%v] [status:%v]\n", ok, status)

	ok, status = act.t().retry.IsRetryable(504)
	fmt.Printf("test: IsRetryable(504) -> [ok:%v] [status:%v]\n", ok, status)

	act.t().retry.SetRateLimiter(100, 10)
	act = t.LookupByName(name)
	ok, status = act.t().retry.IsRetryable(503)
	fmt.Printf("test: IsRetryable(503) -> [ok:%v] [status:%v]\n", ok, status)

	ok, status = act.t().retry.IsRetryable(504)
	fmt.Printf("test: IsRetryable(504) -> [ok:%v] [status:%v]\n", ok, status)

	//Output:
	//test: Add() -> [[]] [count:1]
	//test: IsRetryable(503) -> [ok:true] [status:]
	//test: IsRetryable(504) -> [ok:false] [status:RL]
	//test: IsRetryable(503) -> [ok:true] [status:]
	//test: IsRetryable(504) -> [ok:true] [status:]

}
