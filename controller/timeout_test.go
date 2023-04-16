package controller

import (
	"fmt"
	"net/url"
	"time"
)

func enableValues(enable bool) url.Values {
	v := make(url.Values)
	if enable {
		v.Add(EnabledKey, TrueValue)
	} else {
		v.Add(EnabledKey, FalseValue)
	}
	return v
}

func Example_newTimeout() {
	tbl := newTable(true, false)

	//errt := newTimeout("error-timeout", tbl, NewTimeoutConfig(true, 0, 0))
	//err := errt.validate()
	//fmt.Printf("test: validate() -> [name:%v] [error:%v]\n", errt.name, err)

	err := nilTimeout.validate()
	fmt.Printf("test: validate() -> [name:%v] [error:%v]\n", nilTimeout.name, err)

	t := newTimeout("test-route", tbl, NewTimeoutConfig(true, 0, 100))
	fmt.Printf("test: newTimeout() -> [name:%v] [current:%v]\n", t.name, t.config.Duration)

	t = newTimeout("test-route2", tbl, NewTimeoutConfig(true, 503, time.Millisecond*2000))
	fmt.Printf("test: newTimeout() -> [name:%v] [current:%v]\n", t.name, t.config.Duration)

	t2 := cloneTimeout(t)
	t2.config.Duration = time.Millisecond * 1000
	fmt.Printf("test: cloneTimeout() -> [prev-config:%v] [prev-name:%v] [curr-config:%v] [curr-name:%v]\n", t.config, t.name, t2.config, t2.name)

	//Output:
	//test: validate() -> [name:!] [error:<nil>]
	//test: newTimeout() -> [name:test-route] [current:100ns]
	//test: newTimeout() -> [name:test-route2] [current:2s]
	//test: cloneTimeout() -> [prev-config:{true 503 2s}] [prev-name:test-route2] [curr-config:{true 503 1s}] [curr-name:test-route2]

}

func ExampleTimeout_State() {
	t := newTimeout("test-route", newTable(true, false), NewTimeoutConfig(true, 0, time.Millisecond*2000))
	fmt.Printf("test: newTimeout() -> [name:%v] [state:%v]\n", t.name, t.config)

	//m := make(map[string]string, 16)

	fmt.Printf("test: timeoutState(map,t) -> [enabled:%v] [timeout:%v]\n", t.IsEnabled(), timeoutState(t))

	t.config.Enabled = false
	//m = make(map[string]string, 16)
	//timeoutState(m, t)
	fmt.Printf("test: timeoutState(map,t) -> [enabled:%v] [timeout:%v]\n", t.IsEnabled(), timeoutState(t))

	//Output:
	//test: newTimeout() -> [name:test-route] [state:{true 504 2s}]
	//test: timeoutState(map,t) -> [enabled:true] [timeout:2000]
	//test: timeoutState(map,t) -> [enabled:false] [timeout:-1]

}

func ExampleTimeout_SetTimeout() {
	var v = make(url.Values)
	name := "test-route"
	config := NewTimeoutConfig(true, 504, time.Millisecond*1500)
	t := newTable(true, false)

	errs := t.AddController(newRoute(name, config))
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", errs, t.count())

	ctrl := t.LookupByName(name)
	v.Add(DurationKey, "0")
	err := ctrl.Timeout().Signal(v)
	fmt.Printf("test: Signal(0) -> [%v]\n", err)

	d := ctrl.Timeout().Duration()
	fmt.Printf("test: Duration() -> [%v]\n", d)

	v.Del(DurationKey)
	v.Add(DurationKey, "3s")
	err = ctrl.Timeout().Signal(v)
	ctrl = t.LookupByName(name)
	d = ctrl.Timeout().Duration()
	fmt.Printf("test: Duration() -> [%v][error:%v]\n", d, err)

	//Output:
	//test: Add() -> [[]] [count:1]
	//test: Signal(0) -> [invalid configuration: timeout duration is <= 0]
	//test: Duration() -> [1.5s]
	//test: Duration() -> [3s][error:<nil>]

}

func ExampleTimeout_Toggle() {
	name := "test-route"
	config := NewTimeoutConfig(true, 504, time.Millisecond*1500)
	t := newTable(true, false)

	errs := t.AddController(newRoute(name, config))
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", errs, t.count())

	ctrl := t.LookupByName(name)
	fmt.Printf("test: IsEnabled() -> [%v]\n", ctrl.Timeout().IsEnabled())

	ctrl.Timeout().Signal(enableValues(false))
	ctrl = t.LookupByName(name)
	fmt.Printf("test: IsEnabled() -> [%v]\n", ctrl.Timeout().IsEnabled())

	ctrl.Timeout().Signal(enableValues(true))
	ctrl = t.LookupByName(name)
	fmt.Printf("test: IsEnabled() -> [%v]\n", ctrl.Timeout().IsEnabled())

	//Output:
	//test: Add() -> [[]] [count:1]
	//test: IsEnabled() -> [true]
	//test: IsEnabled() -> [false]
	//test: IsEnabled() -> [true]

}
