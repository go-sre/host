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

	ctrl := EgressTable().LookupUri(uri, method)
	if rlc := ctrl.RateLimiter(); rlc.IsEnabled() && !rlc.Allow() {
		limited = true
		statusFlags = RateLimitFlag
	}
	if !limited {
		if to := ctrl.Timeout(); to.IsEnabled() {
			newCtx, cancelCtx = context.WithTimeout(ctx, to.Duration())
		}
	}
	return func() {
		if cancelCtx != nil {
			cancelCtx()
		}
		code := statusCode()
		if code == StatusDeadlineExceeded {
			statusFlags = UpstreamTimeoutFlag
		}
		ctrl.LogEgress(start, time.Since(start), code, uri, requestId, method, statusFlags)
	}, newCtx, limited
}
