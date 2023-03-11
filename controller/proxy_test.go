package controller

import (
	"fmt"
	"net/url"
)

func ExampleProxy_BuildUrl() {
	t := newTable(true, false)
	uri, _ := url.Parse("https://localhost:8080/basePath/resource?first=false")
	c := newProxy("proxy-route", t, NewProxyConfig("http:"))

	//fmt.Printf("test: newProxy() -> [pattern:%v]\n", c.pattern)
	fmt.Printf("test: InputUrl() -> %v\n", uri.String())

	uri2 := c.BuildUrl(uri)
	fmt.Printf("test: BuildUrl(%v) -> %v\n", c.pattern, uri2.String())

	c = newProxy("proxy-route", t, NewProxyConfig("http://google.com"))
	uri2 = c.BuildUrl(uri)
	fmt.Printf("test: BuildUrl(%v) -> %v\n", c.pattern, uri2.String())

	c = newProxy("proxy-route", t, NewProxyConfig("http://google.com/search"))
	uri2 = c.BuildUrl(uri)
	fmt.Printf("test: BuildUrl(%v) -> %v\n", c.pattern, uri2.String())

	c = newProxy("proxy-route", t, NewProxyConfig("http://google.com/search?q=test"))
	uri2 = c.BuildUrl(uri)
	fmt.Printf("test: BuildUrl(%v) -> %v\n", c.pattern, uri2.String())

	//Output:
	//test: InputUrl() -> https://localhost:8080/basePath/resource?first=false
	//test: BuildUrl(http:) -> http://localhost:8080/basePath/resource?first=false
	//test: BuildUrl(http://google.com) -> http://google.com/basePath/resource?first=false
	//test: BuildUrl(http://google.com/search) -> http://google.com/search?first=false
	//test: BuildUrl(http://google.com/search?q=test) -> http://google.com/search?q=test
}
