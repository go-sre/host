package connector

import (
	"fmt"
	"github.com/go-sre/core/runtime"
	"github.com/go-sre/host/accessdata"
	"net/http"
	"time"
)

/*
func setTestErrorHandler() {
	opt.handler = func(err error) {
		fmt.Printf("test: connector(logd) -> [err:%v]\n", err)
	}
}

*/

func Example_Initialize_Url() {
	status := InitializePush[runtime.DebugError]("", nil)
	fmt.Printf("test: Initialize(\"\") -> [%v] [url:%v]\n", status, pushUrl)

	status = InitializePush[runtime.DebugError]("test", nil)
	fmt.Printf("test: Initialize(\"\") -> [%v] [url:%v]\n", status, pushUrl)

	//Output:
	//[[] github.com/go-sre/host/connector/initialize-push [invalid argument: uri is empty]]
	//test: Initialize("") -> [Internal] [url:]
	//test: Initialize("") -> [OK] [url:test]

}

func Example_Handler_NotProcessed() {
	pushUrl = "http://localhost:8080/accesslog"

	status := pushHandler(nil)
	fmt.Printf("test: handler(nil) -> [%v]\n", status)

	req, _ := http.NewRequest("post", "http://localhost:8080/accesslog", nil)
	data := new(accessdata.Entry)
	data.AddRequest(req)
	status = pushHandler(data)
	fmt.Printf("test: handler(data) -> [%v]\n", status)

	//Output:
	//[[] github.com/go-sre/host/connector/do [invalid argument: access log data is nil]]
	//test: handler(nil) -> [false]
	//test: handler(data) -> [false]

}

func Example_Handler_ConnectFailure() {
	pushUrl = "http://localhost:8080/accesslog"

	req, _ := http.NewRequest("post", "localhost:8081/accesslog", nil)
	data := new(accessdata.Entry)
	data.AddRequest(req)
	status := pushHandler(data)
	fmt.Printf("test: handler(data) -> [%v]\n", status)

	//Output:
	//[[] github.com/go-sre/host/connector/do [Put "http://localhost:8080/accesslog": dial tcp [::1]:8080: connectex: No connection could be made because the target machine actively refused it.]]
	//test: handler(data) -> [false]

}

func Example_Handler_Processed() {
	// Override the message handler
	//handler = func(l *accessdata.Entry) bool {
	//	fmt.Printf("test: handler(logd) -> [%v]\n", accessdata.WriteJson(operators, l))
	//	return true
	//}

	status := InitializePush[runtime.DebugError]("http://localhost:8086/access-log", nil)
	fmt.Printf("test: Initialize() -> [%v]\n", status)

	req, _ := http.NewRequest("POST", "http://localhost:8081/accesslog", nil)
	req.Header.Set("X-Request-ID", "1234-56-7890")
	resp := &http.Response{StatusCode: 200, Proto: "HTTP/1.1", ProtoMajor: 1, ProtoMinor: 1, Header: http.Header{}, Body: nil, ContentLength: 0, TransferEncoding: nil, Close: false, Uncompressed: false, Trailer: http.Header{}, Request: req, TLS: nil}

	extract("egress", time.Now(), time.Millisecond*450, req, resp, "test-route", -1, 50, 5, "false", "true", "RL")
	time.Sleep(time.Second * 2)
	ShutdownPush()

	//Output:
	//test: Initialize() -> [OK]

}
