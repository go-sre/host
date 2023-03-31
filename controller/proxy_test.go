package controller

import (
	"fmt"
	"net/url"
)

func Example_newProxy() {
	t := newTable(true, false)
	p := newProxy("test-route", t, NewProxyConfig(false, "http://localhost:8080", []Header{{"name", "value"}, {"name2", "value2"}}))
	fmt.Printf("test: newProxy() -> [name:%v] [current:%v] [headers:%v]\n", p.name, p.pattern, p.headers)

	p = newProxy("test-route2", t, NewProxyConfig(false, "https://google.com", nil))
	fmt.Printf("test: newProxy() -> [name:%v] [current:%v]\n", p.name, p.pattern)

	p2 := cloneProxy(p)
	p2.pattern = "urn:test"
	fmt.Printf("test: cloneProxy() -> [prev-config:%v] [prev-name:%v] [curr-config:%v] [curr-name:%v]\n", p.pattern, p.name, p2.pattern, p2.name)

	//Output:
	//test: newProxy() -> [name:test-route] [current:http://localhost:8080] [headers:[{name value} {name2 value2}]]
	//test: newProxy() -> [name:test-route2] [current:https://google.com]
	//test: cloneProxy() -> [prev-config:https://google.com] [prev-name:test-route2] [curr-config:urn:test] [curr-name:test-route2]

}

func ExampleProxy_BuildUrl() {
	t := newTable(true, false)
	uri, _ := url.Parse("https://localhost:8080/basePath/resource?first=false")
	c := newProxy("proxy-route", t, NewProxyConfig(false, "http:", nil))

	fmt.Printf("test: InputUrl() -> %v\n", uri.String())

	uri2 := c.BuildUrl(uri)
	fmt.Printf("test: BuildUrl(%v) -> %v\n", c.pattern, uri2.String())

	c = newProxy("proxy-route", t, NewProxyConfig(false, "http://google.com", nil))
	uri2 = c.BuildUrl(uri)
	fmt.Printf("test: BuildUrl(%v) -> %v\n", c.pattern, uri2.String())

	c = newProxy("proxy-route", t, NewProxyConfig(false, "http://google.com/search", nil))
	uri2 = c.BuildUrl(uri)
	fmt.Printf("test: BuildUrl(%v) -> %v\n", c.pattern, uri2.String())

	c = newProxy("proxy-route", t, NewProxyConfig(false, "http://google.com/search?q=test", nil))
	uri2 = c.BuildUrl(uri)
	fmt.Printf("test: BuildUrl(%v) -> %v\n", c.pattern, uri2.String())

	//Output:
	//test: InputUrl() -> https://localhost:8080/basePath/resource?first=false
	//test: BuildUrl(http:) -> http://localhost:8080/basePath/resource?first=false
	//test: BuildUrl(http://google.com) -> http://google.com/basePath/resource?first=false
	//test: BuildUrl(http://google.com/search) -> http://google.com/search?first=false
	//test: BuildUrl(http://google.com/search?q=test) -> http://google.com/search?q=test
}

func ExampleProxy_State() {
	t := newTable(true, false)
	p := newProxy("test-route", t, NewProxyConfig(false, "http://localhost:8080", nil))

	m := make(map[string]string, 16)
	proxyState(m, nil)
	fmt.Printf("test: proxyState(map,nil) -> %v\n", m)
	m = make(map[string]string, 16)
	proxyState(m, p)
	fmt.Printf("test: proxyState(map,p) -> %v\n", m)

	//Output:
	//test: proxyState(map,nil) -> map[proxy:]
	//test: proxyState(map,p) -> map[proxy:false]

}

func ExampleProxy_SetPattern() {
	name := "test-route"
	config := NewProxyConfig(false, "http://localhost:8080", nil)
	t := newTable(true, false)

	errs := t.AddController(newRoute(name, config))
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", errs, t.count())

	ctrl := t.LookupByName(name)
	//pattern := ctrl.t().proxy.Pattern()
	fmt.Printf("test: Pattern() -> [%v]\n", ctrl.Proxy().Pattern())
	prevPattern := ctrl.Proxy().Pattern() //ctrl.(*controller).proxy.Pattern()

	var v = make(url.Values)
	v.Add("pattern", "https://google.com")
	//pxy, _ := ctrl.Proxy()
	ctrl.Proxy().Signal(v)

	ctrl1 := t.LookupByName(name)
	//pattern = ctrl1.t().proxy.Pattern()
	fmt.Printf("test: SetPattern(https://google.com) -> [prev-pattern:%v] [curr-pattern:%v]\n", prevPattern, ctrl1.Proxy().Pattern())

	//Output:
	//test: Add() -> [[]] [count:1]
	//test: Pattern() -> [http://localhost:8080]
	//test: SetPattern(https://google.com) -> [prev-pattern:http://localhost:8080] [curr-pattern:https://google.com]

}

func ExampleProxy_Toggle() {
	name := "test-route"
	config := NewProxyConfig(false, "http://localhost:8080", nil)
	t := newTable(true, false)

	ok := t.AddController(newRoute(name, config))
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", ok, t.count())

	ctrl := t.LookupByName(name)
	enabled := ctrl.t().proxy.IsEnabled()
	fmt.Printf("test: IsEnabled() -> [%v]\n", enabled)
	prevEnabled := ctrl.(*controller).proxy.IsEnabled()

	ctrl.t().proxy.Enable()
	ctrl1 := t.LookupByName(name)
	enabled = ctrl1.t().proxy.IsEnabled()
	fmt.Printf("test: Enable() -> [prev-enabled:%v] [curr-enabled:%v]\n", prevEnabled, enabled)

	//Output:
	//test: Add() -> [[]] [count:1]
	//test: IsEnabled() -> [false]
	//test: Enable() -> [prev-enabled:false] [curr-enabled:true]

}
