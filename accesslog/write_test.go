package accesslog

import (
	"fmt"
	"github.com/go-sre/host/accessdata"
	"net/http"
	"reflect"
	"time"
)

func ExampleLog_Error() {
	start := time.Now()

	Write[TestOutputHandler, accessdata.TextFormatter](nil)
	Write[TestOutputHandler, accessdata.JsonFormatter](accessdata.NewEgressEntry(start, time.Since(start), nil, nil, "egress-route", -1, -1, -1, "", "", ""))

	//Output:
	//test: Write() -> [access data entry is nil]
	//test: Write() -> [{"error":"egress accesslog entries are empty"}]

}

func ExampleLog_Origin() {
	name := "ingress-origin-route"
	start := time.Now()

	accessdata.SetOrigin("us-west", "dfw", "cluster-1", "test-service", "123456-7890-1234")
	err := InitIngressOperators([]accessdata.Operator{{Value: accessdata.StartTimeOperator}, {Value: accessdata.DurationOperator, Name: "duration_ms"},
		{Value: accessdata.TrafficOperator}, {Value: accessdata.RouteNameOperator}, {Value: accessdata.OriginRegionOperator},
		{Value: accessdata.OriginZoneOperator}, {Value: accessdata.OriginSubZoneOperator}, {Value: accessdata.OriginServiceOperator},
		{Value: accessdata.OriginInstanceIdOperator},
	})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	var start1 time.Time
	entry := accessdata.NewIngressEntry(start1, time.Since(start), nil, nil, name, 500, 100, 10, "false", "false", "")
	Write[TestOutputHandler, accessdata.JsonFormatter](entry)
	Write[TestOutputHandler, accessdata.TextFormatter](entry)

	//Output:
	//test: Write() -> [{"start-time":"0001-01-01 00:00:00.000000","duration_ms":0,"traffic":"ingress","route-name":"ingress-origin-route","region":"us-west","zone":"dfw","sub-zone":"cluster-1","service":"test-service","instance-id":"123456-7890-1234"}]
	//test: Write() -> [0001-01-01 00:00:00.000000,0,ingress,ingress-origin-route,us-west,dfw,cluster-1,test-service,123456-7890-1234]

}

/*
func ExampleLog_Ping() {
	name := "ingress-ping-route"
	url := "https://www.google.com/search"

	req, _ := http.NewRequest("", url, nil)
	data.SetPingRoutes([]data.PingRoute{{Traffic: "ingress", Pattern: "/search"}})
	start := time.Now()
	err := InitIngressOperators([]data.Operator{{Value: data.StartTimeOperator}, {Value: data.DurationOperator, Name: "duration_ms"},
		{Value: data.TrafficOperator}, {Value: data.RouteNameOperator}})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	var start1 time.Time
	entry := data.NewHttpIngressEntry(start1, time.Since(start), req, nil, "", map[string]string{data.ActName: name})
	Write[TestOutputHandler, data.JsonFormatter](entry)

	//Output:
	//test: Write() -> [{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"traffic":"ping","route_name":"ingress-ping-route"}]

}


*/
func ExampleLog_Timeout() {
	start := time.Now()

	err := InitEgressOperators([]accessdata.Operator{{Value: accessdata.StartTimeOperator}, {Name: "duration_ms", Value: accessdata.DurationOperator},
		{Value: accessdata.TrafficOperator}, {Value: accessdata.RouteNameOperator}, {Value: accessdata.TimeoutDurationOperator}, {Name: "static", Value: "value"}})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	var start1 time.Time
	Write[TestOutputHandler, accessdata.JsonFormatter](accessdata.NewEgressEntry(start1, time.Since(start), nil, nil, "handler-route", 5000, -1, -1, "", "", ""))

	//Output:
	//test: Write() -> [{"start-time":"0001-01-01 00:00:00.000000","duration_ms":0,"traffic":"egress","route-name":"handler-route","timeout-ms":5000,"static":"value"}]

}

func ExampleLog_RateLimiter_500() {
	start := time.Now()

	err := InitEgressOperators([]accessdata.Operator{{Value: accessdata.StartTimeOperator}, {Name: "duration", Value: accessdata.DurationOperator},
		{Value: accessdata.TrafficOperator}, {Value: accessdata.RouteNameOperator}, {Value: accessdata.RateLimitOperator}, {Value: accessdata.RateBurstOperator}, {Name: "static2", Value: "value2"}})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	var start1 time.Time
	Write[TestOutputHandler, accessdata.JsonFormatter](accessdata.NewEgressEntry(start1, time.Since(start), nil, nil, "handler-route", -1, 500, 10, "", "", ""))

	//Output:
	//test: Write() -> [{"start-time":"0001-01-01 00:00:00.000000","duration":0,"traffic":"egress","route-name":"handler-route","rate-limit":500,"rate-burst":10,"static2":"value2"}]

}

func ExampleLog_Retry_Proxy() {
	start := time.Now()

	err := InitEgressOperators([]accessdata.Operator{{Value: accessdata.StartTimeOperator}, {Value: accessdata.DurationOperator, Name: "duration_ms"},
		{Value: accessdata.TrafficOperator}, {Value: accessdata.RouteNameOperator}, {Value: accessdata.RateLimitOperator},
		{Value: accessdata.RateBurstOperator}, {Value: accessdata.RetryOperator}, {Value: accessdata.ProxyOperator}})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	var start1 time.Time
	Write[TestOutputHandler, accessdata.JsonFormatter](accessdata.NewEgressEntry(start1, time.Since(start), nil, nil, "handler-route", -1, 123, 67, "", "", ""))
	Write[TestOutputHandler, accessdata.JsonFormatter](accessdata.NewEgressEntry(start1, time.Since(start), nil, nil, "handler-route", -1, 123, 67, "true", "false", ""))

	//Output:
	//test: Write() -> [{"start-time":"0001-01-01 00:00:00.000000","duration_ms":0,"traffic":"egress","route-name":"handler-route","rate-limit":123,"rate-burst":67,"retry":null,"proxy":null}]
	//test: Write() -> [{"start-time":"0001-01-01 00:00:00.000000","duration_ms":0,"traffic":"egress","route-name":"handler-route","rate-limit":123,"rate-burst":67,"retry":true,"proxy":false}]

}

func ExampleLog_Request() {
	req, _ := http.NewRequest("", "www.google.com/search/documents", nil)
	req.Header.Add("customer", "Ted's Bait & Tackle")

	var start time.Time
	err := InitEgressOperators([]accessdata.Operator{{Value: accessdata.RequestProtocolOperator}, {Value: accessdata.RequestMethodOperator}, {Value: accessdata.RequestUrlOperator},
		{Value: accessdata.RequestPathOperator}, {Value: accessdata.RequestHostOperator}, {Value: "%REQ(customer)%"}})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	Write[TestOutputHandler, accessdata.JsonFormatter](accessdata.NewEgressEntry(start, time.Since(start), nil, nil, "handler-route", -1, -1, -1, "", "", ""))
	Write[TestOutputHandler, accessdata.JsonFormatter](accessdata.NewEgressEntry(start, time.Since(start), req, nil, "handler-route", -1, -1, -1, "", "", ""))

	//Output:
	//test: Write() -> [{"protocol":null,"method":null,"url":null,"path":null,"host":null,"customer":null}]
	//test: Write() -> [{"protocol":"HTTP/1.1","method":"GET","url":"www.google.com/search/documents","path":"www.google.com/search/documents","host":null,"customer":"Ted's Bait & Tackle"}]

}

func ExampleLog_Response() {
	resp := &http.Response{StatusCode: 404, ContentLength: 1234}

	err := InitEgressOperators([]accessdata.Operator{{Value: accessdata.ResponseStatusCodeOperator}, {Value: accessdata.ResponseBytesReceivedOperator}, {Value: accessdata.StatusFlagsOperator}})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	var start time.Time
	Write[TestOutputHandler, accessdata.JsonFormatter](accessdata.NewEgressEntry(start, time.Since(start), nil, nil, "handler-route", -1, -1, -1, "", "", "UT"))
	Write[TestOutputHandler, accessdata.JsonFormatter](accessdata.NewEgressEntry(start, time.Since(start), nil, resp, "handler-route", -1, -1, -1, "", "", "UT"))

	//Output:
	//test: Write() -> [{"status-code":0,"bytes-received":0,"status-flags":"UT"}]
	//test: Write() -> [{"status-code":404,"bytes-received":1234,"status-flags":"UT"}]

}

func _Example_Log_State() {
	t := time.Duration(time.Millisecond * 500)
	i := reflect.TypeOf(t)
	a := any(t)

	fmt.Printf("test 1 -> %v\n", a)

	fmt.Printf("test 2 -> %v\n", i)

	//Output:
	//fail
}
