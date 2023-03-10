package middleware

import (
	"fmt"
	"github.com/gotemplates/host/accessdata"
	"github.com/gotemplates/host/accesslog"
	"net/http"
)

var (
	isEnabled    = false
	googleUrl    = "https://www.google.com/search?q=test"
	instagramUrl = "https://www.instagram.com"
	config       = []accessdata.Operator{
		{Value: accessdata.StartTimeOperator},
		{Value: accessdata.DurationOperator},
		{Value: accessdata.TrafficOperator},

		{Value: accessdata.RequestMethodOperator},
		{Value: accessdata.RequestHostOperator},
		{Value: accessdata.RequestPathOperator},
		{Value: accessdata.RequestProtocolOperator},

		{Value: accessdata.ResponseStatusCodeOperator},
		{Value: accessdata.StatusFlagsOperator},
		{Value: accessdata.ResponseBytesSentOperator},
	}
)

func init() {
	accesslog.InitEgressOperators(config)

}

func Example_AccessLog_No_Wrapper() {
	req, _ := http.NewRequest("GET", googleUrl, nil)

	// Testing - check for a nil wrapper or round tripper
	w := accessWrapper{}
	resp, err := w.RoundTrip(req)
	fmt.Printf("test: RoundTrip(wrapper:nil) -> [resp:%v] [err:%v]\n", resp, err)

	// Testing - no wrapper, calling Google search
	resp, err = http.DefaultClient.Do(req)
	fmt.Printf("test: RoundTrip(handler:false) -> [status_code:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: RoundTrip(wrapper:nil) -> [resp:<nil>] [err:invalid handler round tripper configuration : http.RoundTripper is nil]
	//test: RoundTrip(handler:false) -> [status_code:200] [err:<nil>]

}

func Example_AccessLog_Default() {
	req, _ := http.NewRequest("GET", instagramUrl, nil)

	if !isEnabled {
		isEnabled = true
		AccessLogWrapTransport(nil)
	}
	resp, err := http.DefaultClient.Do(req)
	fmt.Printf("test: RoundTrip(handler:true) -> [status_code:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: RoundTrip(handler:true) -> [status_code:200] [err:<nil>]

}
