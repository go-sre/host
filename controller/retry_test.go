package controller

import (
	"fmt"
	"golang.org/x/time/rate"
)

func Example_newRetry() {
	//t := newTable(true, false)

	rt := newRetry("test-route", newTable(true, false), NewRetryConfig(false, 5, 10, 0, []int{504}))
	limit, burst := rt.LimitAndBurst()
	fmt.Printf("test: newRetry() -> [name:%v] [config:%v] [limit:%v] [burst:%v]\n", rt.name, rt.config, limit, burst)

	rt = newRetry("test-route2", newTable(true, false), NewRetryConfig(false, 2, 20, 0, []int{503, 504}))
	fmt.Printf("test: newRetry() -> [name:%v] [config:%v]\n", rt.name, rt.config)

	rt2 := cloneRetry(rt)
	//t2.Enable()
	rt2.Signal(EnableValues(true))
	fmt.Printf("test: cloneRetry() -> [prev-enabled:%v] [curr-enabled:%v]\n", rt.IsEnabled(), rt2.IsEnabled())

	//t = newRetry("test-route3", newTable(true), NewRetryConfig([]int{503, 504}, time.Millisecond*2000, false))
	fmt.Printf("test: retryState(nil,false,map) -> %v\n", retryState(nil, nil, false))

	fmt.Printf("test: retryState(t,false,map) -> %v\n", retryState(nil, rt, false))

	rt2 = newRetry("test-route", newTable(true, false), NewRetryConfig(false, rate.Inf, 10, 0, []int{504}))
	fmt.Printf("test: retryState(t2,true,map) -> %v\n", retryState(nil, rt2, true))

	//Output:
	//test: newRetry() -> [name:test-route] [config:{5 10 0s [504]}] [limit:5] [burst:10]
	//test: newRetry() -> [name:test-route2] [config:{2 20 0s [503 504]}]
	//test: cloneRetry() -> [prev-enabled:true] [curr-enabled:false]
	//test: retryState(nil,false,map) -> map[retry: retryBurst:-1 retryRateLimit:-1]
	//test: retryState(t,false,map) -> map[retry:false retryBurst:20 retryRateLimit:2]
	//test: retryState(t2,true,map) -> map[retry:true retryBurst:10 retryRateLimit:99999]

}

func ExampleRetry_Toggle() {
	name := "test-route"
	config := NewRetryConfig(true, 5, 10, 0, []int{504})
	t := newTable(true, false)
	err := t.AddController(newRoute(name, config))
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", err, t.count())

	ctrl := t.LookupByName(name)
	fmt.Printf("test: IsEnabled() -> [%v]\n", ctrl.Retry().IsEnabled())
	prevEnabled := ctrl.Retry().IsEnabled()

	ctrl.Retry().Signal(EnableValues(false))
	ctrl1 := t.LookupByName(name)
	fmt.Printf("test: Disable() -> [prev-enabled:%v] [curr-enabled:%v]\n", prevEnabled, ctrl1.Retry().IsEnabled())
	prevEnabled = ctrl1.Retry().IsEnabled()

	ctrl1.Retry().Signal(EnableValues(true))
	ctrl = t.LookupByName(name)
	fmt.Printf("test: Enable() -> [prev-enabled:%v] [curr-enabled:%v]\n", prevEnabled, ctrl.Retry().IsEnabled())
	prevEnabled = ctrl.Retry().IsEnabled()

	//Output:
	//test: Add() -> [[]] [count:1]
	//test: IsEnabled() -> [true]
	//test: Disable() -> [prev-enabled:true] [curr-enabled:false]
	//test: Enable() -> [prev-enabled:false] [curr-enabled:true]

}

func ExampleIsRetryable_Disabled() {
	name := "test-route"
	config := NewRetryConfig(false, 100, 10, 0, []int{503, 504})
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

func ExampleIsRetryable_StatusCode() {
	name := "test-route"
	config := NewRetryConfig(false, 100, 10, 0, []int{503, 504})
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
	config := NewRetryConfig(false, 1, 1, 0, []int{503, 504})
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

	//act.t().retry.SetRateLimiter(100, 10)
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
