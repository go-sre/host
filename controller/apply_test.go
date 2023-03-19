package controller

import (
	"context"
	"time"
)

const (
	applyTestUri = "urn:postgres:query.access-log"
)

type testStatus struct {
	code uint32
}

func newStatusOK() *testStatus {
	return &testStatus{code: 0}
}

func newStatusCode(code uint32) *testStatus {
	return &testStatus{code: code}
}

func (t *testStatus) Code() uint32 {
	return t.code
}

func init() {
	/*
		defaultLogFn = func(start time.Time, duration time.Duration, statusCode int, uri, requestId, method, statusFlags string, actuatorState map[string]string) {
			s := fmt.Sprintf("traffic:%v ,"+
				"route:%v ,"+
				"request-id:%v, "+
				"status-code:%v, "+
				"method:%v, "+
				"url:%v, "+
				"host:%v, "+
				"path:%v, "+
				"timeout:%v, "+
				"rate-limit:%v, "+
				"rate-burst:%v, "+
				"retry:%v, "+
				"retry-rate-limit:%v, "+
				"retry-rate-burst:%v, "+
				"status-flags:%v",
				"egress", actuatorState[ActName], requestId, statusCode, method, uri, "postgres", "query.access-log",
				actuatorState[TimeoutName],
				actuatorState[RateLimitName], actuatorState[RateBurstName],
				actuatorState[RetryName], actuatorState[RetryRateLimitName], actuatorState[RetryRateBurstName],
				statusFlags)
			fmt.Printf("{%v}\n", s)
		}

	*/
}

func ExampleEgressApply() {
	function(context.Background())

	//Output:
	//{traffic:egress ,route:* ,request-id:123-456-7890, status-code:0, method:GET, url:urn:postgres:query.access-log, host:postgres, path:query.access-log, timeout:-1, rate-limit:-1, rate-burst:-1, retry:, retry-rate-limit:-1, retry-rate-burst:-1, status-flags:}

}

func ExampleEgressApply_RateLimit() {
	name := "rate-limit-route"
	egressTable = NewEgressTable()

	route := NewRoute(name, EgressTraffic, "", false, NewRateLimiterConfig(1, 0, 503))
	EgressTable().AddController(route)
	EgressTable().SetUriMatcher(func(uri string, method string) (string, bool) {
		return name, true
	})

	functionRateLimited(context.Background())

	//Output:
	//{traffic:egress ,route:rate-limit-route ,request-id:123-456-7890, status-code:94, method:GET, url:urn:postgres:query.access-log, host:postgres, path:query.access-log, timeout:-1, rate-limit:1, rate-burst:0, retry:, retry-rate-limit:-1, retry-rate-burst:-1, status-flags:RL}

}

func ExampleEgressApply_Timeout() {
	name := "timeout-route"
	egressTable = NewEgressTable()

	route := NewRoute(name, EgressTraffic, "", false, NewTimeoutConfig(time.Second, 504))
	EgressTable().AddController(route)
	EgressTable().SetUriMatcher(func(uri string, method string) (string, bool) {
		return name, true
	})

	functionTimeout(context.Background())

	//Output:
	//{traffic:egress ,route:timeout-route ,request-id:123-456-7890, status-code:4, method:GET, url:urn:postgres:query.access-log, host:postgres, path:query.access-log, timeout:1000, rate-limit:-1, rate-burst:-1, retry:, retry-rate-limit:-1, retry-rate-burst:-1, status-flags:UT}

}

func function(ctx context.Context) (status *testStatus) {
	var fn func()

	fn, ctx, _ = EgressApply(ctx, func() int { return int((*(&status)).Code()) }, applyTestUri, "123-456-7890", "GET")
	defer fn()
	return newStatusOK()
}

func functionRateLimited(ctx context.Context) (status *testStatus) {
	var fn func()
	var limited = false

	fn, ctx, limited = EgressApply(ctx, func() int { return int((*(&status)).Code()) }, applyTestUri, "123-456-7890", "GET")
	defer fn()
	if limited {
		return newStatusCode(StatusRateLimited)
	}
	return newStatusOK()
}

func functionTimeout(ctx context.Context) (status *testStatus) {
	var fn func()
	var limited = false

	fn, ctx, limited = EgressApply(ctx, func() int { return int((*(&status)).Code()) }, applyTestUri, "123-456-7890", "GET")
	defer fn()
	if limited {
		return newStatusCode(StatusRateLimited)
	}
	done := make(chan struct{})
	panicChan := make(chan any, 1)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()
		workFunction()
		close(done)
	}()
	// Waiting for events
	select {
	case p := <-panicChan:
		panic(p)
	case <-done:
		break
	case <-ctx.Done():
		switch err := ctx.Err(); err {
		case context.DeadlineExceeded:
			return newStatusCode(StatusDeadlineExceeded)
		default:
		}
	}
	return newStatusOK()
}

func workFunction() {
	time.Sleep(time.Second * 2)
	return
}
