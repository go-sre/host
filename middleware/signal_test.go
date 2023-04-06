package middleware

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

func init() {

}

func ExampleSignalHandler_InvalidArgument() {
	req, _ := http.NewRequest("GET", "localhost:8080/signal", nil)
	record := httptest.NewRecorder()
	SignalHandler(record, req)
	resp := record.Result()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("test: SignalHandler(nil) -> [statusCode:%v] [body:%v]\n", resp.StatusCode, string(body))

	req, _ = http.NewRequest("GET", "localhost:8080/signal?traffic=egress", nil)
	record = httptest.NewRecorder()
	SignalHandler(record, req)
	resp = record.Result()
	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("test: SignalHandler(traffic=egress) -> [statusCode:%v] [body:%v]\n", resp.StatusCode, string(body))

	req, _ = http.NewRequest("GET", "localhost:8080/signal?traffic=egress&route=timeout-route", nil)
	record = httptest.NewRecorder()
	SignalHandler(record, req)
	resp = record.Result()
	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("test: SignalHandler(traffic=egress&route=timeout-route) -> [statusCode:%v] [body:%v]\n", resp.StatusCode, string(body))

	req, _ = http.NewRequest("GET", "localhost:8080/signal?traffic=egress&route=timeout-route&behavior=proxy", nil)
	record = httptest.NewRecorder()
	SignalHandler(record, req)
	resp = record.Result()
	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("test: SignalHandler(traffic=egress&route=timeout-route&behavior=proxy) -> [statusCode:%v] [body:%v]\n", resp.StatusCode, string(body))

	//Output:
	//test: SignalHandler(nil) -> [statusCode:400] [body:invalid argument: route [] not found in [] table]
	//test: SignalHandler(traffic=egress) -> [statusCode:400] [body:invalid argument: route [] not found in [egress] table]
	//test: SignalHandler(traffic=egress&route=timeout-route) -> [statusCode:400] [body:invalid argument: behavior [] is not supported]
	//test: SignalHandler(traffic=egress&route=timeout-route&behavior=proxy) -> [statusCode:400] [body:invalid signal: proxy is not configured]
	
}

func _ExampleSignalHandler() {
	req, _ := http.NewRequest("GET", "localhost:8080/signal?traffic=egress&route=timeout-route&behavoir=", nil)
	record := httptest.NewRecorder()
	SignalHandler(record, req)
	resp := record.Result()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("test: SignalHandler(nil) -> [statusCode:%v] [body:%v]\n", resp.StatusCode, string(body))

	//Output:
	//test: SignalHandler(nil) -> [statusCode:400] [body:invalid argument: route [] not found in [] table]
	//test: SignalHandler(traffic=egress) -> [statusCode:400] [body:invalid argument: route [] not found in [egress] table]

}
