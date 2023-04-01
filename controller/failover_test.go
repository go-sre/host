package controller

import (
	"fmt"
)

var failoverFn FailoverInvoke = func(name string, failover bool) { fmt.Printf("test: Invoke(%v,%v)\n", name, failover) }

func Example_newFailover() {
	name := "failover-test"

	f := newFailover(name, nil, nil)
	fmt.Printf("test: newFailover(nil) -> [enabled:%v] [validate:%v]\n", f.enabled, f.validate())

	f = newFailover(name, nil, NewFailoverConfig(false, failoverFn))
	fmt.Printf("test: newFailover(testFn) -> [enabled:%v] [validate:%v]\n", f.enabled, f.validate())

	f2 := cloneFailover(f)
	f2.enabled = true
	fmt.Printf("test: cloneFailover(f1) -> [f2-enabled:%v] [f2-validate:%v]\n", f2.enabled, f2.validate())

	//Output:
	//test: newFailover(nil) -> [enabled:false] [validate:invalid configuration: Failover FailureInvoke function is nil]
	//test: newFailover(testFn) -> [enabled:false] [validate:<nil>]
	//test: cloneFailover(f1) -> [f2-enabled:true] [f2-validate:<nil>]

}

func ExampleFailover_State() {
	name := "failover-test"
	f := newFailover(name, nil, NewFailoverConfig(false, failoverFn))

	m := make(map[string]string, 16)
	failoverState(m, f)
	fmt.Printf("test: failoverState(map,f1) -> [enabled:%v] %v\n", f.IsEnabled(), m)

	m = make(map[string]string, 16)
	f.enabled = true
	failoverState(m, f)
	fmt.Printf("test: failoverState(map,f2) -> [enabled:%v] %v\n", f.IsEnabled(), m)

	//Output:
	//test: failoverState(map,f1) -> [enabled:false] map[failover:false]
	//test: failoverState(map,f2) -> [enabled:true] map[failover:true]

}

func ExampleFailover_Toggle() {
	prevEnabled := false
	name := "failover-test"
	t := newTable(true, false)

	errs := t.AddController(newRoute(name, NewFailoverConfig(false, failoverFn)))
	fmt.Printf("test: Add() -> [error:%v] [count:%v]\n", errs, t.count())

	ctrl := t.LookupByName(name)
	fmt.Printf("test: IsEnabled() -> [%v]\n", ctrl.Failover().IsEnabled())
	prevEnabled = ctrl.Failover().IsEnabled()

	ctrl.Failover().Signal(EnableValues(false))
	ctrl2 := t.LookupByName(name)
	fmt.Printf("test: Disable() -> [prev-enabled:%v] [curr-enabled:%v]\n", prevEnabled, ctrl2.Failover().IsEnabled())
	prevEnabled = ctrl2.Failover().IsEnabled()

	ctrl2.Failover().Signal(EnableValues(true))
	ctrl = t.LookupByName(name)
	fmt.Printf("test: Enable() -> [prev-enabled:%v] [curr-enabled:%v]\n", prevEnabled, ctrl.Failover().IsEnabled())

	//Output:
	//test: Add() -> [error:[]] [count:1]
	//test: IsEnabled() -> [false]
	//test: Disable() -> [prev-enabled:false] [curr-enabled:false]
	//test: Enable() -> [prev-enabled:false] [curr-enabled:true]

}

func ExampleFailover_Invoke() {
	name := "failover-test"
	t := newTable(true, false)
	err := t.AddController(newRoute(name, NewFailoverConfig(false, failoverFn)))
	fmt.Printf("test: Add() -> [error:%v] [count:%v]\n", err, t.count())

	f := t.LookupByName(name)
	f.t().failover.Invoke(true)
	fmt.Printf("test: Invoke(true) -> []\n")

	f.t().failover.Invoke(false)
	fmt.Printf("test: Invoke(false) -> []\n")

	//Output:
	//test: Add() -> [error:[]] [count:1]
	//test: Invoke(failover-test,true)
	//test: Invoke(true) -> []
	//test: Invoke(failover-test,false)
	//test: Invoke(false) -> []
}
