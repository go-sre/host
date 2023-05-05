package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func behaviorValues(behavior string) url.Values {
	v := make(url.Values)
	v.Add(BehaviorKey, behavior)
	return v
}

func ExampleController_newController() {
	t := newTable(true, false)
	route := NewRoute("test", EgressTraffic, "", false, NewTimeoutConfig(true, 0, time.Millisecond*1500), NewRateLimiterConfig(true, 503, 100, 10, ""))

	ctrl, _ := newController(route, t)
	fmt.Printf("test: newController() -> [timeout:%v] [rateLimit:%v] [retry:%v]\n", ctrl.Timeout().IsEnabled(), ctrl.RateLimiter().IsEnabled(), ctrl.Retry().IsEnabled())

	d := ctrl.timeout.Duration()
	a1 := cloneController[*timeout](ctrl, newTimeout("new-timeout", t, NewTimeoutConfig(true, http.StatusGatewayTimeout, time.Millisecond*500)))

	d1 := a1.timeout.Duration()
	fmt.Printf("test: cloneController() -> [prev-duration:%v] [curr-duration:%v]\n", d, d1)

	//Output:
	//test: newController() -> [timeout:true] [rateLimit:true] [retry:false]
	//test: cloneController() -> [prev-duration:1.5s] [curr-duration:500ms]

}

func ExampleController_newController_config() {
	t := newTable(true, false)
	route := NewRoute("test", EgressTraffic, "", false, NewTimeoutConfig(true, 0, time.Millisecond*1500), nil, NewRateLimiterConfig(true, 503, 100, 10, ""), nil)

	ctrl, _ := newController(route, t)
	fmt.Printf("test: newController() -> [timeout:%v] [rateLimit:%v] [retry:%v]\n", ctrl.Timeout().IsEnabled(), ctrl.RateLimiter().IsEnabled(), ctrl.Retry().IsEnabled())

	//d := ctrl.timeout.Duration()
	//ctrl1 := cloneController[*timeout](ctrl, newTimeout("new-timeout", t, NewTimeoutConfig(time.Millisecond*500, http.StatusGatewayTimeout)))

	//d1 := ctrl1.timeout.Duration()
	//fmt.Printf("test: cloneController() -> [prev-duration:%v] [curr-duration:%v]\n", d, d1)

	//ctrl.Actuate(nil)

	//Output:
	//test: newController() -> [timeout:true] [rateLimit:true] [retry:false]

}

func ExampleController_newController_Error() {
	t := newTable(false, false)
	route := NewRoute("test", IngressTraffic, "", false, NewTimeoutConfig(true, 0, time.Millisecond*1500), NewRateLimiterConfig(true, 503, 100, 10, ""))

	_, errs := newController(route, t)
	fmt.Printf("test: newController() -> [errs:%v]\n", errs)

	route = NewRoute("test", IngressTraffic, "", false, NewTimeoutConfig(true, 0, time.Millisecond*1500), NewRetryConfig(false, 100, 10, 0, nil))
	_, errs = newController(route, t)
	fmt.Printf("test: newController() -> [errs:%v]\n", errs)

	route = NewRoute("test", IngressTraffic, "", false, NewTimeoutConfig(true, 0, -1))
	_, errs = newController(route, t)
	fmt.Printf("test: newController() -> [errs:%v]\n", errs)

	route = NewRoute("test", IngressTraffic, "", false, NewTimeoutConfig(true, 0, 10))
	_, errs = newController(route, t)
	fmt.Printf("test: newController() -> [errs:%v]\n", errs)

	route = newRoute("test", NewRateLimiterConfig(true, 504, -1, 10, ""))
	_, errs = newController(route, t)
	fmt.Printf("test: newController() -> [errs:%v]\n", errs)

	//Output:
	//test: newController() -> [errs:[]]
	//test: newController() -> [errs:[invalid configuration: retry status codes are empty [test]]]
	//test: newController() -> [errs:[invalid configuration: Timeout duration is < 0 [test]]]
	//test: newController() -> [errs:[]]
	//test: newController() -> [errs:[invalid configuration: RateLimiter limit is < 0 [test]]]

}

func ExampleController_Signal() {
	t := newTable(false, false)
	route := NewRoute("test", IngressTraffic, "", false, NewTimeoutConfig(true, 0, time.Millisecond*1500), NewRateLimiterConfig(true, 503, 100, 10, ""))

	ctrl, errs := newController(route, t)
	fmt.Printf("test: newController() -> [errs:%v]\n", errs)

	err := ctrl.Signal(behaviorValues(TimeoutBehavior))
	fmt.Printf("test: Signal(timeout) -> [err:%v]\n", err)

	err = ctrl.Signal(behaviorValues("invalid-behavior"))
	fmt.Printf("test: Signal(timeout) -> [err:%v]\n", err)

	//Output:
	//test: newController() -> [errs:[]]
	//test: Signal(timeout) -> [err:<nil>]
	//test: Signal(timeout) -> [err:invalid argument: behavior [invalid-behavior] is not supported]

}
