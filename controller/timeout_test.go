package controller

import (
	"fmt"
	"net/url"
	"time"
)

func Example_newTimeout() {
	tbl := newTable(true, false)
	t := newTimeout("test-route", tbl, NewTimeoutConfig(true, 0, 100))
	fmt.Printf("test: newTimeout() -> [name:%v] [current:%v]\n", t.name, t.config.Duration)

	t = newTimeout("test-route2", tbl, NewTimeoutConfig(true, 503, time.Millisecond*2000))
	fmt.Printf("test: newTimeout() -> [name:%v] [current:%v]\n", t.name, t.config.Duration)

	t2 := cloneTimeout(t)
	t2.config.Duration = time.Millisecond * 1000
	fmt.Printf("test: cloneTimeout() -> [prev-config:%v] [prev-name:%v] [curr-config:%v] [curr-name:%v]\n", t.config, t.name, t2.config, t2.name)

	//Output:
	//test: newTimeout() -> [name:test-route] [current:100ns]
	//test: newTimeout() -> [name:test-route2] [current:2s]
	//test: cloneTimeout() -> [prev-config:{true 503 2s}] [prev-name:test-route2] [curr-config:{true 503 1s}] [curr-name:test-route2]

}

func ExampleTimeout_State() {
	t := newTimeout("test-route", newTable(true, false), NewTimeoutConfig(true, 0, time.Millisecond*2000))
	fmt.Printf("test: newTimeout() -> [name:%v] [state:%v]\n", t.name, t.config)

	m := make(map[string]string, 16)
	timeoutState(m, t)
	fmt.Printf("test: timeoutState(map,t) -> [enabled:%v] %v\n", t.IsEnabled(), m)

	t.config.Enabled = false
	m = make(map[string]string, 16)
	timeoutState(m, t)
	fmt.Printf("test: timeoutState(map,t) -> [enabled:%v] %v\n", t.IsEnabled(), m)

	//Output:
	//test: newTimeout() -> [name:test-route] [state:{true 504 2s}]
	//test: timeoutState(map,t) -> [enabled:true] map[timeout:2000]
	//test: timeoutState(map,t) -> [enabled:false] map[timeout:-1]

}

func ExampleTimeout_SetTimeout() {
	var v = make(url.Values)
	name := "test-route"
	config := NewTimeoutConfig(true, 504, time.Millisecond*1500)
	t := newTable(true, false)

	errs := t.AddController(newRoute(name, config))
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", errs, t.count())

	ctrl := t.LookupByName(name)
	d := ctrl.Timeout().Duration()
	fmt.Printf("test: Duration() -> [%v]\n", d)

	v.Add(DurationKey, "3s")
	ctrl.Timeout().Signal(v)
	ctrl = t.LookupByName(name)
	d = ctrl.Timeout().Duration()
	fmt.Printf("test: Duration() -> [%v]\n", d)

	//Output:
	//test: Add() -> [[]] [count:1]
	//test: Duration() -> [1.5s]
	//test: Duration() -> [3s]

}

func ExampleTimeout_Toggle() {
	name := "test-route"
	config := NewTimeoutConfig(true, 504, time.Millisecond*1500)
	t := newTable(true, false)

	errs := t.AddController(newRoute(name, config))
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", errs, t.count())

	ctrl := t.LookupByName(name)
	fmt.Printf("test: IsEnabled() -> [%v]\n", ctrl.Timeout().IsEnabled())

	ctrl.Timeout().Signal(EnableValues(false))
	ctrl = t.LookupByName(name)
	fmt.Printf("test: IsEnabled() -> [%v]\n", ctrl.Timeout().IsEnabled())

	ctrl.Timeout().Signal(EnableValues(true))
	ctrl = t.LookupByName(name)
	fmt.Printf("test: IsEnabled() -> [%v]\n", ctrl.Timeout().IsEnabled())

	//Output:
	//test: Add() -> [[]] [count:1]
	//test: IsEnabled() -> [true]
	//test: IsEnabled() -> [false]
	//test: IsEnabled() -> [true]

}
