package controller

import "fmt"

func ExampleParseUri_Url() {
	uri := "www.google.com"
	scheme, host, path := ParseUri(uri)
	fmt.Printf("test: ParseUri(%v) -> [scheme:%v] [host:%v] [path:%v]\n", uri, scheme, host, path)

	uri = "https://www.google.com"
	scheme, host, path = ParseUri(uri)
	fmt.Printf("test: ParseUri(%v) -> [scheme:%v] [host:%v] [path:%v]\n", uri, scheme, host, path)

	uri = "https://www.google.com/search?q=test"
	scheme, host, path = ParseUri(uri)
	fmt.Printf("test: ParseUri(%v) -> [scheme:%v] [host:%v] [path:%v]\n", uri, scheme, host, path)

	//Output:
	//test: ParseUri(www.google.com) -> [scheme:] [host:] [path:www.google.com]
	//test: ParseUri(https://www.google.com) -> [scheme:https] [host:www.google.com] [path:]
	//test: ParseUri(https://www.google.com/search?q=test) -> [scheme:https] [host:www.google.com] [path:/search]

}

func ExampleParseUri_Urn() {
	uri := "urn"
	scheme, host, path := ParseUri(uri)
	fmt.Printf("test: ParseUri(%v) -> [scheme:%v] [host:%v] [path:%v]\n", uri, scheme, host, path)

	uri = "urn:postgres"
	scheme, host, path = ParseUri(uri)
	fmt.Printf("test: ParseUri(%v) -> [scheme:%v] [host:%v] [path:%v]\n", uri, scheme, host, path)

	uri = "urn:postgres:query.access-log"
	scheme, host, path = ParseUri(uri)
	fmt.Printf("test: ParseUri(%v) -> [scheme:%v] [host:%v] [path:%v]\n", uri, scheme, host, path)

	//Output:
	//test: ParseUri(urn) -> [scheme:] [host:] [path:urn]
	//test: ParseUri(urn:postgres) -> [scheme:urn] [host:postgres] [path:]
	//test: ParseUri(urn:postgres:query.access-log) -> [scheme:urn] [host:postgres] [path:query.access-log]

}
