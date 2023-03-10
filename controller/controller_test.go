package controller

import (
	"fmt"
	"github.com/gotemplates/host/accessdata"
	"net/http"
	"time"
)

/*
func ExampleActuator_newActuator() {
	t := newTable(true)

	a, _ := newActuator("test", t, controllerFn, NewTimeoutConfig(time.Millisecond*1500, 0), NewRateLimiterConfig(100, 10, 503))

	_, toOk := a.Timeout()
	_, rateOk := a.RateLimiter()
	_, retryOk := a.Retry()
	_, failOk := a.Failover()
	fmt.Printf("test: newActuator() -> [logger:%v] [timeout:%v] [rateLimit:%v] [retry:%v] [failover:%v]\n", a.Logger() != nil, toOk, rateOk, retryOk, failOk)

	d := a.timeout.Duration()
	a1 := cloneActuator[*timeout](a, newTimeout("new-timeout", t, NewTimeoutConfig(time.Millisecond*500, http.StatusGatewayTimeout)))

	d1 := a1.timeout.Duration()
	fmt.Printf("test: cloneActuator() -> [prev-duration:%v] [curr-duration:%v]\n", d, d1)

	a.Actuate(nil)

	//Output:
	//test: newActuator() -> [logger:true] [timeout:true] [rateLimit:true] [retry:false] [failover:false]
	//test: cloneActuator() -> [prev-duration:1.5s] [curr-duration:500ms]
	//test: Actuate() -> [test]

}
*/

func ExampleController_newController() {
	t := newTable(true, false)
	route := NewRoute("test", accessdata.EgressTraffic, "", false, NewTimeoutConfig(time.Millisecond*1500, 0), NewRateLimiterConfig(100, 10, 503))

	ctrl, _ := newController(route, t)

	_, toOk := ctrl.Timeout()
	_, rateOk := ctrl.RateLimiter()
	_, retryOk := ctrl.Retry()
	_, failOk := ctrl.Failover()
	fmt.Printf("test: newController() -> [timeout:%v] [rateLimit:%v] [retry:%v] [failover:%v]\n", toOk, rateOk, retryOk, failOk)

	d := ctrl.timeout.Duration()
	a1 := cloneController[*timeout](ctrl, newTimeout("new-timeout", t, NewTimeoutConfig(time.Millisecond*500, http.StatusGatewayTimeout)))

	d1 := a1.timeout.Duration()
	fmt.Printf("test: cloneController() -> [prev-duration:%v] [curr-duration:%v]\n", d, d1)

	//Output:
	//test: newController() -> [timeout:true] [rateLimit:true] [retry:false] [failover:false]
	//test: cloneController() -> [prev-duration:1.5s] [curr-duration:500ms]

}

func ExampleController_newController_config() {
	t := newTable(true, false)
	route := NewRoute("test", accessdata.EgressTraffic, "", false, NewTimeoutConfig(time.Millisecond*1500, 0), nil, NewRateLimiterConfig(100, 10, 503), nil)

	ctrl, _ := newController(route, t)

	_, toOk := ctrl.Timeout()
	_, rateOk := ctrl.RateLimiter()
	_, retryOk := ctrl.Retry()
	_, failOk := ctrl.Failover()
	fmt.Printf("test: newController() -> [timeout:%v] [rateLimit:%v] [retry:%v] [failover:%v]\n", toOk, rateOk, retryOk, failOk)

	//d := ctrl.timeout.Duration()
	//ctrl1 := cloneController[*timeout](ctrl, newTimeout("new-timeout", t, NewTimeoutConfig(time.Millisecond*500, http.StatusGatewayTimeout)))

	//d1 := ctrl1.timeout.Duration()
	//fmt.Printf("test: cloneController() -> [prev-duration:%v] [curr-duration:%v]\n", d, d1)

	//ctrl.Actuate(nil)

	//Output:
	//test: newController() -> [timeout:true] [rateLimit:true] [retry:false] [failover:false]

}

func ExampleController_newController_Error() {
	t := newTable(false, false)
	route := NewRoute("test", accessdata.IngressTraffic, "", false, NewTimeoutConfig(time.Millisecond*1500, 0), NewRateLimiterConfig(100, 10, 503))

	_, errs := newController(route, t)
	fmt.Printf("test: newController() -> [errs:%v]\n", errs)

	route = NewRoute("test", accessdata.IngressTraffic, "", false, NewTimeoutConfig(time.Millisecond*1500, 0), NewRetryConfig(nil, 100, 10, 0))
	_, errs = newController(route, t)
	fmt.Printf("test: newController() -> [errs:%v]\n", errs)

	route = NewRoute("test", accessdata.IngressTraffic, "", false, NewTimeoutConfig(0, 0))
	_, errs = newController(route, t)
	fmt.Printf("test: newController() -> [errs:%v]\n", errs)

	route = NewRoute("test", accessdata.IngressTraffic, "", false, NewTimeoutConfig(10, 0), NewFailoverConfig(nil))
	_, errs = newController(route, t)
	fmt.Printf("test: newController() -> [errs:%v]\n", errs)

	route = newRoute("test", NewRateLimiterConfig(-1, 10, 504))
	_, errs = newController(route, t)
	fmt.Printf("test: newController() -> [errs:%v]\n", errs)

	//Output:
	//test: newController() -> [errs:[]]
	//test: newController() -> [errs:[invalid configuration: Retry status codes are empty]]
	//test: newController() -> [errs:[invalid configuration: Timeout duration is <= 0]]
	//test: newController() -> [errs:[invalid configuration: Failover FailureInvoke function is nil]]
	//test: newController() -> [errs:[invalid configuration: RateLimiter limit is < 0]]

}
