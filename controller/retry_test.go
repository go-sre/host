package controller

import (
	"fmt"
	"time"
)

func Example_newRetry() {
	tbl := newTable(true, false)

	rt := newRetry("test-route", tbl, NewRetryConfig(false, 100, 10, time.Second, []int{504}))
	fmt.Printf("test: newRetry() -> [name:%v] [limit:%v] [burst:%v] [wait:%v] [codes:%v]\n", rt.name, rt.config.Limit, rt.config.Burst, rt.config.Wait, rt.config.StatusCodes)

	rt = newRetry("test-route2", tbl, NewRetryConfig(false, 200, 20, time.Millisecond*500, []int{503, 504}))
	fmt.Printf("test: newRetry() -> [name:%v] [limit:%v] [burst:%v] [wait:%v] [codes:%v]\n", rt.name, rt.config.Limit, rt.config.Burst, rt.config.Wait, rt.config.StatusCodes)

	rt2 := cloneRetry(rt)
	rt2.config.Enabled = true
	rt2.config.Limit = 50
	fmt.Printf("test: cloneRetry() -> [prev-enabled:%v] [curr-enabled:%v] [prev-limit:%v] [curr-limit:%v] \n", rt.IsEnabled(), rt2.IsEnabled(), rt.config.Limit, rt2.config.Limit)

	//Output:
	//test: newRetry() -> [name:test-route] [limit:100] [burst:10] [wait:1s] [codes:[504]]
	//test: newRetry() -> [name:test-route2] [limit:200] [burst:20] [wait:500ms] [codes:[503 504]]
	//test: cloneRetry() -> [prev-enabled:false] [curr-enabled:true] [prev-limit:200] [curr-limit:50]

}

func ExampleRetry_State() {
	tbl := newTable(true, false)

	rt := newRetry("test-route3", tbl, NewRetryConfig(false, 100, 10, time.Millisecond*500, []int{503, 504}))
	fmt.Printf("test: retryState(nil,nil,false) -> %v\n", retryState(nil, nil, false))

	fmt.Printf("test: retryState(nil,rt,true) -> %v\n", retryState(nil, rt, true))

	rt.config.Enabled = true
	fmt.Printf("test: retryState(nil,rt,true) -> %v\n", retryState(nil, rt, true))

	//Output:
	//test: retryState(nil,nil,false) -> map[retry:false retryBurst:-1 retryRateLimit:-1]
	//test: retryState(nil,rt,true) -> map[retry:false retryBurst:-1 retryRateLimit:-1]
	//test: retryState(nil,rt,true) -> map[retry:true retryBurst:10 retryRateLimit:100]

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

	ctrl.Retry().Signal(enableValues(false))
	ctrl1 := t.LookupByName(name)
	fmt.Printf("test: Disable() -> [prev-enabled:%v] [curr-enabled:%v]\n", prevEnabled, ctrl1.Retry().IsEnabled())
	prevEnabled = ctrl1.Retry().IsEnabled()

	ctrl1.Retry().Signal(enableValues(true))
	ctrl = t.LookupByName(name)
	fmt.Printf("test: Enable() -> [prev-enabled:%v] [curr-enabled:%v]\n", prevEnabled, ctrl.Retry().IsEnabled())
	prevEnabled = ctrl.Retry().IsEnabled()

	//Output:
	//test: Add() -> [[]] [count:1]
	//test: IsEnabled() -> [true]
	//test: Disable() -> [prev-enabled:true] [curr-enabled:false]
	//test: Enable() -> [prev-enabled:false] [curr-enabled:true]

}

func ExampleRetry_IsRetryable_Disabled() {
	name := "test-route"
	config := NewRetryConfig(false, 100, 10, 0, []int{503, 504})
	t := newTable(true, false)
	err := t.AddController(newRoute(name, config))
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", err, t.count())

	ctrl := t.LookupByName(name)
	ok, status := ctrl.Retry().IsRetryable(200)
	fmt.Printf("test: IsRetryable(200) -> [ok:%v] [status:%v]\n", ok, status)

	ok, status = ctrl.Retry().IsRetryable(503)
	fmt.Printf("test: IsRetryable(503) -> [ok:%v] [status:%v]\n", ok, status)

	ok, status = ctrl.Retry().IsRetryable(504)
	fmt.Printf("test: IsRetryable(504) -> [ok:%v] [status:%v]\n", ok, status)

	//Output:
	//test: Add() -> [[]] [count:1]
	//test: IsRetryable(200) -> [ok:false] [status:NE]
	//test: IsRetryable(503) -> [ok:false] [status:NE]
	//test: IsRetryable(504) -> [ok:false] [status:NE]

}

func ExampleRetry_IsRetryable_StatusCode() {
	name := "test-route"
	config := NewRetryConfig(false, 100, 10, 0, []int{503, 504})
	t := newTable(true, false)
	err := t.AddController(newRoute(name, config))
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", err, t.count())

	ctrl := t.LookupByName(name)
	ctrl.Retry().Enable()
	ctrl = t.LookupByName(name)
	ok, status := ctrl.Retry().IsRetryable(200)
	fmt.Printf("test: IsRetryable(200) -> [ok:%v] [status:%v]\n", ok, status)

	ok, status = ctrl.Retry().IsRetryable(500)
	fmt.Printf("test: IsRetryable(500) -> [ok:%v] [status:%v]\n", ok, status)

	ok, status = ctrl.Retry().IsRetryable(502)
	fmt.Printf("test: IsRetryable(502) -> [ok:%v] [status:%v]\n", ok, status)

	ok, status = ctrl.Retry().IsRetryable(503)
	fmt.Printf("test: IsRetryable(503) -> [ok:%v] [status:%v]\n", ok, status)

	ok, status = ctrl.Retry().IsRetryable(504)
	fmt.Printf("test: IsRetryable(504) -> [ok:%v] [status:%v]\n", ok, status)

	ok, status = ctrl.Retry().IsRetryable(505)
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

func ExampleRetry_IsRetryable_RateLimit() {
	name := "test-route"
	config := NewRetryConfig(false, 1, 1, 0, []int{503, 504})
	t := newTable(true, false)
	err := t.AddController(newRoute(name, config))
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", err, t.count())

	ctrl := t.LookupByName(name)
	ctrl.Retry().Enable()
	ctrl = t.LookupByName(name)
	ok, status := ctrl.Retry().IsRetryable(503)
	fmt.Printf("test: IsRetryable(503) -> [ok:%v] [status:%v]\n", ok, status)

	ok, status = ctrl.Retry().IsRetryable(504)
	fmt.Printf("test: IsRetryable(504) -> [ok:%v] [status:%v]\n", ok, status)

	ctrl.Retry().Signal(rateLimiterSetValues(100, 10))
	ctrl = t.LookupByName(name)
	ok, status = ctrl.Retry().IsRetryable(503)
	fmt.Printf("test: IsRetryable(503) -> [ok:%v] [status:%v]\n", ok, status)

	ok, status = ctrl.Retry().IsRetryable(504)
	fmt.Printf("test: IsRetryable(504) -> [ok:%v] [status:%v]\n", ok, status)

	//Output:
	//test: Add() -> [[]] [count:1]
	//test: IsRetryable(503) -> [ok:true] [status:]
	//test: IsRetryable(504) -> [ok:false] [status:RL]
	//test: IsRetryable(503) -> [ok:true] [status:]
	//test: IsRetryable(504) -> [ok:true] [status:]

}
