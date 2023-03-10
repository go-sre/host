package controller

import (
	"context"
	"time"
)

// Shadowed from : https://grpc.github.io/grpc/core/md_doc_statuscodes.html

const (
	StatusDeadlineExceeded = 4
	StatusRateLimited      = 94
)

// EgressApply - function to be used by non Http egress traffic to apply an controller
func EgressApply(ctx context.Context, statusCode func() int, uri, requestId, method string) (func(), context.Context, bool) {
	statusFlags := ""
	limited := false
	start := time.Now()
	newCtx := ctx
	var cancelCtx context.CancelFunc

	act := EgressTable.LookupUri(uri, method)
	if rlc, ok := act.RateLimiter(); ok && !rlc.Allow() {
		limited = true
		statusFlags = RateLimitFlag
	}
	if !limited {
		if toc, ok := act.Timeout(); ok {
			newCtx, cancelCtx = context.WithTimeout(ctx, toc.Duration())
		}
	}
	return func() {
		if cancelCtx != nil {
			cancelCtx()
		}
		//code := (*status).Code()
		code := statusCode()
		if code == StatusDeadlineExceeded {
			statusFlags = UpstreamTimeoutFlag
		}
		act.LogEgress(start, time.Since(start), code, uri, requestId, method, statusFlags)
	}, newCtx, limited
}
