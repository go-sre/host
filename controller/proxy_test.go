package controller

import (
	"errors"
	"fmt"
	"net/url"
)

type testAction struct{}

func (testAction) Signal(values url.Values) error {
	return errors.New("test action error")
}

type testAction2 struct{}

func (testAction2) Signal(values url.Values) error {
	return errors.New("test action 2 error")
}

func Example_newProxy() {
	t := newTable(true, false)
	p := newProxy("test-route", t, NewProxyConfig(false, "http://localhost:8080", []Header{{"name", "value"}, {"name2", "value2"}}, nil))
	fmt.Printf("test: newProxy() -> [name:%v] [current:%v] [headers:%v]\n", p.name, p.config.Pattern, p.config.Headers)

	p = newProxy("test-route2", t, NewProxyConfig(false, "https://google.com", nil, nil))
	fmt.Printf("test: newProxy() -> [name:%v] [current:%v]\n", p.name, p.config.Pattern)

	err := nilProxy.validate()
	fmt.Printf("test: validate() -> [name:%v] [error:%v]\n", nilProxy.name, err)
	err = p.validate()
	fmt.Printf("test: validate() -> [name:%v] [error:%v]\n", p.name, err)

	p2 := cloneProxy(p)
	p2.config.Pattern = "urn:test"
	fmt.Printf("test: cloneProxy() -> [prev-config:%v] [prev-name:%v] [curr-config:%v] [curr-name:%v]\n", p.config.Pattern, p.name, p2.config.Pattern, p2.name)

	//Output:
	//test: newProxy() -> [name:test-route] [current:http://localhost:8080] [headers:[{name value} {name2 value2}]]
	//test: newProxy() -> [name:test-route2] [current:https://google.com]
	//test: validate() -> [name:!] [error:<nil>]
	//test: validate() -> [name:test-route2] [error:<nil>]
	//test: cloneProxy() -> [prev-config:https://google.com] [prev-name:test-route2] [curr-config:urn:test] [curr-name:test-route2]

}

func ExampleProxy_State() {
	t := newTable(true, false)
	p := newProxy("test-route", t, NewProxyConfig(false, "http://localhost:8080", nil, nil))

	//m := make(map[string]string, 16)
	//proxyState(p)
	fmt.Printf("test: proxyState(p) -> [enabled:%v] [proxied:%v]\n", p.IsEnabled(), proxyState(p))
	//m = make(map[string]string, 16)
	p.config.Enabled = true

	fmt.Printf("test: proxyState(p) -> [enabled:%v] [proxied:%v]\n", p.IsEnabled(), proxyState(p))

	//Output:
	//test: proxyState(p) -> [enabled:false] [proxied:]
	//test: proxyState(p) -> [enabled:true] [proxied:true]

}

func ExampleProxy_BuildUrl() {
	t := newTable(true, false)
	uri, _ := url.Parse("https://localhost:8080/basePath/resource?first=false")
	c := newProxy("proxy-route", t, NewProxyConfig(false, "http:", nil, nil))

	fmt.Printf("test: InputUrl() -> %v\n", uri.String())

	uri2 := c.BuildUrl(uri)
	fmt.Printf("test: BuildUrl(%v) -> %v\n", c.config.Pattern, uri2.String())

	c = newProxy("proxy-route", t, NewProxyConfig(false, "http://google.com", nil, nil))
	uri2 = c.BuildUrl(uri)
	fmt.Printf("test: BuildUrl(%v) -> %v\n", c.config.Pattern, uri2.String())

	c = newProxy("proxy-route", t, NewProxyConfig(false, "http://google.com/search", nil, nil))
	uri2 = c.BuildUrl(uri)
	fmt.Printf("test: BuildUrl(%v) -> %v\n", c.config.Pattern, uri2.String())

	c = newProxy("proxy-route", t, NewProxyConfig(false, "http://google.com/search?q=test", nil, nil))
	uri2 = c.BuildUrl(uri)
	fmt.Printf("test: BuildUrl(%v) -> %v\n", c.config.Pattern, uri2.String())

	//Output:
	//test: InputUrl() -> https://localhost:8080/basePath/resource?first=false
	//test: BuildUrl(http:) -> http://localhost:8080/basePath/resource?first=false
	//test: BuildUrl(http://google.com) -> http://google.com/basePath/resource?first=false
	//test: BuildUrl(http://google.com/search) -> http://google.com/search?first=false
	//test: BuildUrl(http://google.com/search?q=test) -> http://google.com/search?q=test
}

func ExampleProxy_SetPattern() {
	name := "test-route"
	config := NewProxyConfig(false, "http://localhost:8080", nil, nil)
	t := newTable(true, false)

	errs := t.AddController(newRoute(name, config))
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", errs, t.count())

	ctrl := t.LookupByName(name)
	fmt.Printf("test: Pattern() -> [%v]\n", ctrl.Proxy().Pattern())
	prevPattern := ctrl.Proxy().Pattern()

	err := ctrl.Proxy().Signal(NewValues(PatternKey, ""))
	fmt.Printf("test: Signal() -> [error:%v]\n", err)

	err = ctrl.Proxy().Signal(NewValues(PatternKey, "https://google.com/{0x34567"))
	fmt.Printf("test: Signal() -> [error:%v]\n", err)

	err = ctrl.Proxy().Signal(NewValues(PatternKey, "https://google.com"))

	ctrl1 := t.LookupByName(name)
	fmt.Printf("test: SetPattern(https://google.com) -> [prev-pattern:%v] [curr-pattern:%v] [error:%v]\n", prevPattern, ctrl1.Proxy().Pattern(), err)

	//Output:
	//test: Add() -> [[]] [count:1]
	//test: Pattern() -> [http://localhost:8080]
	//test: Signal() -> [error:invalid configuration: proxy pattern is empty]
	//test: Signal() -> [error:<nil>]
	//test: SetPattern(https://google.com) -> [prev-pattern:http://localhost:8080] [curr-pattern:https://google.com] [error:<nil>]

}

func ExampleProxy_Toggle() {
	name := "test-route"
	config := NewProxyConfig(false, "http://localhost:8080", nil, nil)
	t := newTable(true, false)

	ok := t.AddController(newRoute(name, config))
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", ok, t.count())

	ctrl := t.LookupByName(name)
	fmt.Printf("test: IsEnabled() -> [%v]\n", ctrl.t().proxy.IsEnabled())
	prevEnabled := ctrl.Proxy().IsEnabled()

	ctrl.Proxy().Signal(enableValues(true))
	ctrl1 := t.LookupByName(name)
	fmt.Printf("test: Enable() -> [prev-enabled:%v] [curr-enabled:%v]\n", prevEnabled, ctrl1.Proxy().IsEnabled())

	//Output:
	//test: Add() -> [[]] [count:1]
	//test: IsEnabled() -> [false]
	//test: Enable() -> [prev-enabled:false] [curr-enabled:true]

}

func ExampleProxy_Action() {
	var action testAction
	var action2 testAction2
	name := "test-route"
	config := NewProxyConfig(true, "urn:postgresql:host:path", nil, action)
	t := newTable(true, false)

	ok := t.AddController(newRoute(name, config))
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", ok, t.count())

	ctrl := t.LookupByName(name)
	fmt.Printf("test: IsEnabled() -> [%v] [action:%v]\n", ctrl.Proxy().IsEnabled(), ctrl.Proxy().Action() != nil)
	fmt.Printf("test: Action().Signal(nil) -> [%v]\n", ctrl.Proxy().Action().Signal(nil))

	ctrl.Proxy().SetAction(action2)
	ctrl = t.LookupByName(name)
	fmt.Printf("test: Action2().Signal(nil) -> [%v]\n", ctrl.Proxy().Action().Signal(nil))

	//Output:
	//test: Add() -> [[]] [count:1]
	//test: IsEnabled() -> [true] [action:true]
	//test: Action().Signal(nil) -> [test action error]
	//test: Action2().Signal(nil) -> [test action 2 error]

}

func ExampleProxy_SetAction() {
	var action testAction
	name := "test-route"
	config := NewProxyConfig(true, "urn:postgresql:host:path", nil, nil)
	t := newTable(true, false)

	ok := t.AddController(newRoute(name, config))
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", ok, t.count())

	ctrl := t.LookupByName(name)
	fmt.Printf("test: IsEnabled() -> [%v] [action:%v]\n", ctrl.Proxy().IsEnabled(), ctrl.Proxy().Action() != nil)

	t.SetAction(name, action)
	ctrl = t.LookupByName(name)
	fmt.Printf("test: IsEnabled() -> [%v] [action:%v]\n", ctrl.Proxy().IsEnabled(), ctrl.Proxy().Action() != nil)

	fmt.Printf("test: Action().Signal(nil) -> [%v]\n", ctrl.Proxy().Action().Signal(nil))

	//Output:
	//test: Add() -> [[]] [count:1]
	//test: IsEnabled() -> [true] [action:false]
	//test: IsEnabled() -> [true] [action:true]
	//test: Action().Signal(nil) -> [test action error]

}

func ExampleProxy_SignalAction() {
	var action testAction
	name := "test-route"
	config := NewProxyConfig(true, "urn:postgresql:host:path", nil, action)
	t := newTable(true, false)

	ok := t.AddController(newRoute(name, config))
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", ok, t.count())

	ctrl := t.LookupByName(name)
	fmt.Printf("test: IsEnabled() -> [%v] [action:%v] [pattern:%v]\n", ctrl.Proxy().IsEnabled(), ctrl.Proxy().Action() != nil, ctrl.Proxy().Pattern())

	err := ctrl.Proxy().Signal(NewValues(PatternKey, "urn:postgresql:host:path"))
	ctrl = t.LookupByName(name)
	fmt.Printf("test: Signal() -> [%v] [action:%v] [pattern:%v] [error:%v]\n", ctrl.Proxy().IsEnabled(), ctrl.Proxy().Action() != nil, ctrl.Proxy().Pattern(), err)

	err = ctrl.Proxy().Signal(NewValues(PatternKey, "urn:postgresql:host2:path"))
	ctrl = t.LookupByName(name)
	fmt.Printf("test: Signal() -> [%v] [action:%v] [pattern:%v] [error:%v]\n", ctrl.Proxy().IsEnabled(), ctrl.Proxy().Action() != nil, ctrl.Proxy().Pattern(), err)

	//Output:
	//test: Add() -> [[]] [count:1]
	//test: IsEnabled() -> [true] [action:true] [pattern:urn:postgresql:host:path]
	//test: Signal() -> [true] [action:true] [pattern:urn:postgresql:host:path] [error:<nil>]
	//test: Signal() -> [true] [action:true] [pattern:urn:postgresql:host2:path] [error:test action error]

}
