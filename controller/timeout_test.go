package controller

import (
	"fmt"
	"time"
)

func Example_newTimeout() {
	t := newTimeout("test-route", newTable(true, false), NewTimeoutConfig(true, 0, 100))
	fmt.Printf("test: newTimeout() -> [name:%v] [current:%v]\n", t.name, t.config.Duration)

	t = newTimeout("test-route2", newTable(true, false), NewTimeoutConfig(true, 503, time.Millisecond*2000))
	fmt.Printf("test: newTimeout() -> [name:%v] [current:%v]\n", t.name, t.config.Duration)

	t2 := cloneTimeout(t)
	t2.config.Duration = time.Millisecond * 1000
	fmt.Printf("test: cloneTimeout() -> [prev-config:%v] [prev-name:%v] [curr-config:%v] [curr-name:%v]\n", t.config, t.name, t2.config, t2.name)

	//Output:
	//test: newTimeout() -> [name:test-route] [current:100ns]
	//test: newTimeout() -> [name:test-route2] [current:2s]
	//test: cloneTimeout() -> [prev-config:{false 503 2s}] [prev-name:test-route2] [curr-config:{false 503 1s}] [curr-name:test-route2]

}

func Example_Timeout_State() {
	t := newTimeout("test-route", newTable(true, false), NewTimeoutConfig(true, 0, time.Millisecond*2000))

	d := t.Duration()
	fmt.Printf("test: Duration() -> [%v]\n", d)

	t = newTimeout("test-route", newTable(true, false), NewTimeoutConfig(true, 0, time.Millisecond*2000))

	m := make(map[string]string, 16)
	timeoutState(m, nil)
	fmt.Printf("test: timeoutState(map,nil) -> %v\n", m)
	m = make(map[string]string, 16)
	timeoutState(m, t)
	fmt.Printf("test: timeoutState(map,t) -> %v\n", m)

	//Output:
	//test: Duration() -> [2s]
	//test: timeoutState(map,nil) -> map[timeout:-1]
	//test: timeoutState(map,t) -> map[timeout:2000]

}

func Example_Timeout_SetTimeout() {
	name := "test-route"
	config := NewTimeoutConfig(true, 0, time.Millisecond*1500)
	t := newTable(true, false)

	ok := t.AddController(newRoute(name, config))
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", ok, t.count())

	ctrl := t.LookupByName(name)
	d := ctrl.t().timeout.Duration()
	fmt.Printf("test: Duration() -> [%v]\n", d)
	prevDuration := ctrl.(*controller).timeout.Duration()

	ctrl.t().timeout.SetTimeout(time.Second * 2)
	ctrl1 := t.LookupByName(name)
	d = ctrl1.t().timeout.Duration()
	fmt.Printf("test: SetTimeout(2s) -> [prev-duration:%v] [curr-duration:%v]\n", prevDuration, d)
	prevDuration = ctrl1.t().timeout.Duration()

	m := make(map[string]string, 16)
	timeoutState(m, ctrl1.t().timeout)
	fmt.Printf("test: timeoutState(map,t) -> %v\n", m)

	//Output:
	//test: Add() -> [[]] [count:1]
	//test: Duration() -> [1.5s]
	//test: SetTimeout(2s) -> [prev-duration:1.5s] [curr-duration:2s]
	//test: timeoutState(map,t) -> map[timeout:2000]

}
