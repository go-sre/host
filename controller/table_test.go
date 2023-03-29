package controller

import (
	"fmt"
	"net/http"
	"time"
)

func ExampleTable_SetDefaultController_Egress() {
	t := newTable(true, false)

	fmt.Printf("test: empty() -> [%v]\n", t.isEmpty())

	a := t.LookupHttp(nil)
	fmt.Printf("test: LookupHttp(nil) -> [default:%v]\n", a.(*controller).name == DefaultControllerName)
	//fmt.Printf("IsDefault : %v\n", r.(*route).name == DefaultName)

	t.SetDefaultController(newRoute("not-default"))
	a = t.LookupHttp(nil)
	fmt.Printf("test: LookupHttp(req) -> [default:%v]\n", a.(*controller).name == DefaultControllerName)

	//Output:
	//test: empty() -> [true]
	//test: LookupHttp(nil) -> [default:true]
	//test: LookupHttp(req) -> [default:false]

}

func ExampleTable_SetDefaultController_Ingress() {
	t := newTable(false, false)

	fmt.Printf("test: empty() -> [%v]\n", t.isEmpty())

	a := t.LookupHttp(nil)
	fmt.Printf("test: LookupHttp(nil) -> [default:%v]\n", a.(*controller).name == DefaultControllerName)
	//fmt.Printf("IsDefault : %v\n", r.(*route).name == DefaultName)

	err := t.SetDefaultController(newRoute("not-default"))
	//a = t.LookupHttp(nil)
	fmt.Printf("test: SetDefaultController(newRoute) -> %v\n", err) //a.(*controller).name == DefaultActuatorName)

	//Output:
	//test: empty() -> [true]
	//test: LookupHttp(nil) -> [default:true]
	//test: SetDefaultController(newRoute) -> []

}

func ExampleTable_SetHostController_Egress() {
	t := newTable(true, false)

	a := t.Host()
	fmt.Printf("test: Host() -> [name:%v] [timeout-controller:%v]\n", a.Name(), a.t().timeout)

	err := t.SetHostController(newRoute("", NewTimeoutConfig(true, 504, time.Second)))
	//a = t.Host()
	fmt.Printf("test: SetHostController(NewTimeoutConfig()) -> %v\n", err)

	//Output:
	//test: Host() -> [name:host] [timeout-controller:<nil>]
	//test: SetHostController(NewTimeoutConfig()) -> [host controller configuration is not valid for egress traffic]

}

func ExampleTable_SetHostController_Ingress() {
	t := newTable(false, false)
	a := t.Host()
	fmt.Printf("test: Host() -> [name:%v]\n", a.Name())

	err := t.SetHostController(newRoute("", NewTimeoutConfig(true, 503, time.Second)))
	fmt.Printf("test: SetHostController(newRoute) -> %v\n", err)

	err = t.SetHostController(newRoute("", NewRateLimiterConfig(true, 504, 100, 25)))
	fmt.Printf("test: SetHostController(newRoute) -> %v\n", err)

	//Output:
	//test: Host() -> [name:host]
	//test: SetHostController(newRoute) -> [host controller configuration does not allow retry, rate limiter, or failover controllers]
	//test: SetHostController(newRoute) -> []

}

func ExampleTable_Add_Exists_LookupByName() {
	name := "test-route"
	t := newTable(true, false)
	fmt.Printf("test: empty() -> [%v]\n", t.isEmpty())

	err := t.AddController(newRoute("", nil, nil, nil, nil))
	fmt.Printf("test: Add(nil) -> [err:%v] [count:%v] [exists:%v] [lookup:%v]\n", err, t.count(), t.exists(name), t.LookupByName(name))

	err = t.AddController(newRoute(name, nil, nil, nil, nil))
	fmt.Printf("test: Add(controller) -> [err:%v] [count:%v] [exists:%v] [lookup:%v]\n", err, t.count(), t.exists(name), t.LookupByName(name) != nil)

	t.remove("")
	fmt.Printf("test: remove(\"\") -> [count:%v] [exists:%v] [lookup:%v]\n", t.count(), t.exists(name), t.LookupByName(name) != nil)

	t.remove(name)
	fmt.Printf("test: remove(name) -> [count:%v] [exists:%v] [lookup:%v]\n", t.count(), t.exists(name), t.LookupByName(name))

	//Output:
	//test: empty() -> [true]
	//test: Add(nil) -> [err:[invalid argument: route name is empty]] [count:0] [exists:false] [lookup:<nil>]
	//test: Add(controller) -> [err:[]] [count:1] [exists:true] [lookup:true]
	//test: remove("") -> [count:1] [exists:true] [lookup:true]
	//test: remove(name) -> [count:0] [exists:false] [lookup:<nil>]

}

func ExampleTable_LookupHttp() {
	name := "test-route"
	t := newTable(true, false)
	fmt.Printf("test: empty() -> [%v]\n", t.isEmpty())

	r := t.LookupHttp(nil)
	fmt.Printf("test: LookupHttp(nil) -> [controller:%v]\n", r.Name())

	req, _ := http.NewRequest("", "http://localhost:8080/accesslog", nil)
	r = t.LookupHttp(req)
	fmt.Printf("test: LookupHttp(req) -> [controller:%v]\n", r.Name())

	t.SetHttpMatcher(func(req *http.Request) (string, bool) {
		return name, true
	},
	)
	ok := t.AddController(newRoute(name, NewTimeoutConfig(true, 503, 100), nil, nil, nil))
	fmt.Printf("test: Add(controller) -> [controller:%v] [count:%v] [exists:%v]\n", ok, t.count(), t.exists(name))

	r = t.LookupHttp(req)
	fmt.Printf("test: LookupHttp(req) ->  [controller:%v]\n", r.Name())

	//Output:
	//test: empty() -> [true]
	//test: LookupHttp(nil) -> [controller:*]
	//test: LookupHttp(req) -> [controller:*]
	//test: Add(controller) -> [controller:[]] [count:1] [exists:true]
	//test: LookupHttp(req) ->  [controller:test-route]

}

func ExampleTable_LookupUri() {
	name := "test-route"
	t := newTable(true, false)
	fmt.Printf("test: empty() -> [%v]\n", t.isEmpty())

	r := t.LookupUri("", "")
	fmt.Printf("test: LookupUri(nil,nil) -> [controller:%v]\n", r.Name())

	uri := "urn:postgres:query.access-log"
	r = t.LookupUri(uri, "")
	fmt.Printf("test: LookupUri(%v) -> [controller:%v]\n", uri, r.Name())

	t.SetUriMatcher(func(uri, method string) (string, bool) {
		return name, true
	},
	)
	ok := t.AddController(newRoute(name, NewTimeoutConfig(true, 503, 100), nil, nil, nil))
	fmt.Printf("test: Add(controller) -> [controller:%v] [count:%v] [exists:%v]\n", ok, t.count(), t.exists(name))

	r = t.LookupUri(uri, "")
	fmt.Printf("test: LookupUri(%v) ->  [controller:%v]\n", uri, r.Name())

	//Output:
	//test: empty() -> [true]
	//test: LookupUri(nil,nil) -> [controller:*]
	//test: LookupUri(urn:postgres:query.access-log) -> [controller:*]
	//test: Add(controller) -> [controller:[]] [count:1] [exists:true]
	//test: LookupUri(urn:postgres:query.access-log) ->  [controller:test-route]

}

func ExampleTable_Lookup_Name_Default() {
	//name := "test-route"
	t := newTable(true, true)
	fmt.Printf("test: empty() -> [%v]\n", t.isEmpty())

	act := t.LookupByName("")
	fmt.Printf("test: Lookup(nil) -> [controller:%v]\n", act != nil)

	act = t.LookupByName("test")
	fmt.Printf("test: Lookup(\"test\") -> [controller:%v]\n", act != nil)

	//Output:
	//test: empty() -> [true]
	//test: Lookup(nil) -> [controller:false]
	//test: Lookup("test") -> [controller:true]

}

func ExampleTable_Add_Ingress() {
	t := newTable(false, false)

	err := t.AddController(newRoute("valid"))
	fmt.Printf("test: AddRoute(valid) -> %v\n", err)

	err = t.AddController(newRoute("invalid", NewTimeoutConfig(true, 504, time.Second)))
	fmt.Printf("test: AddRoute(invalid) -> %v\n", err)

	//Output:
	//test: AddRoute(valid) -> []
	//test: AddRoute(invalid) -> []

}
