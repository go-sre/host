package accessdata

import (
	"fmt"
	"net/http"
	"time"
)

func _Example_Value_Duration() {
	start := time.Now()

	time.Sleep(time.Second * 2)
	data := &Entry{}
	data.Duration = time.Since(start)
	fmt.Printf("test: Value(\"Duration\") -> [%v]\n", data.Value(DurationOperator))
	fmt.Printf("test: Value(\"DurationString\") -> [%v]\n", data.Value(DurationStringOperator))

	//Output:
	//test: Value("Duration") -> [2011]
	//test: Value("DurationString") -> [2.0117949s]

}

func Example_Value_Controller() {
	name := "test-route"
	op := RouteNameOperator
	start := time.Now().UTC()

	data := Entry{}
	fmt.Printf("test: Value(\"%v\") -> [%v]\n", name, data.Value(op))

	data = Entry{RouteName: name}
	fmt.Printf("test: Value(\"%v\") -> [route_name:%v]\n", name, data.Value(op))

	data1 := NewEntry(PingTraffic, start, time.Since(start), nil, nil, name, -1, -1, -1, "", "", "")
	fmt.Printf("test: Value(\"%v\") -> [traffic:%v]\n", name, data1.Value(TrafficOperator))

	data = Entry{Timeout: 500}
	fmt.Printf("test: Value(\"%v\") -> [timeout:%v]\n", name, data.Value(TimeoutDurationOperator))

	//Output:
	//test: Value("test-route") -> []
	//test: Value("test-route") -> [route_name:test-route]
	//test: Value("test-route") -> [traffic:ping]
	//test: Value("test-route") -> [timeout:500]

}

func Example_Value_Request() {
	op := RequestMethodOperator

	data := &Entry{}
	fmt.Printf("test: Value(\"method\") -> [%v]\n", data.Value(op))

	req, _ := http.NewRequest("POST", "www.google.com", nil)
	//req.Header.Add(RequestIdHeaderName, uuid.New().String())
	req.Header.Add(RequestIdHeaderName, "123-456-789")
	req.Header.Add(FromRouteHeaderName, "calling-route")
	data = &Entry{}
	data.AddRequest(req)
	fmt.Printf("test: Value(\"method\") -> [%v]\n", data.Value(op))

	fmt.Printf("test: Value(\"headers\") -> [request-id:%v] [from-route:%v]\n", data.Value(RequestIdOperator), data.Value(RequestFromRouteOperator))

	//Output:
	//test: Value("method") -> []
	//test: Value("method") -> [POST]
	//test: Value("headers") -> [request-id:123-456-789] [from-route:calling-route]
}

func Example_Value_Response() {
	op := ResponseStatusCodeOperator

	data := &Entry{}
	fmt.Printf("test: Value(\"code\") -> [%v]\n", data.Value(op))

	resp := &http.Response{StatusCode: 200}
	data = &Entry{}
	data.AddResponse(resp)
	fmt.Printf("test: Value(\"code\") -> [%v]\n", data.Value(op))

	//Output:
	//test: Value("code") -> [0]
	//test: Value("code") -> [200]
}

func Example_Value_Request_Header() {
	req, _ := http.NewRequest("", "www.google.com", nil)
	req.Header.Add("customer", "Ted's Bait & Tackle")
	data := Entry{}
	data.AddRequest(req)
	fmt.Printf("test: Value(\"customer\") -> [%v]\n", data.Value("%REQ(customer)%"))

	//Output:
	//test: Value("customer") -> [Ted's Bait & Tackle]
}

func Example_EgressEntry() {
	var start time.Time

	url := "urn:postgres:query.access-accesslog"
	requestId := "123-456-789"
	req, _ := http.NewRequest("QUERY", url, nil)
	req.Header.Add(RequestIdHeaderName, requestId)

	resp := new(http.Response)
	resp.StatusCode = 201

	e := NewEgressEntry(start, 0, req, resp, "egress-route", -1, -1, -1, "", "", "RL")
	fmt.Printf("test: String() -> {%v}\n", e)

	//Output:
	//test: String() -> {traffic:egress, route:egress-route, request-id:123-456-789, status-code:201, protocol:urn, method:QUERY, url:urn:postgres:query.access-accesslog, host:postgres, path:query.access-accesslog, timeout:-1, rate-limit:-1, rate-burst:-1, retry:, proxy:, status-flags:RL}

}
