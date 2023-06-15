package controller

import (
	"context"
	"github.com/go-sre/core/runtime"
	"time"
)

// Shadowed from : https://grpc.github.io/grpc/core/md_doc_statuscodes.html

const (
	StatusDeadlineExceeded = 4
	StatusRateLimited      = 94
)

type EgressController interface {
	Apply(ctx context.Context, statusCode func() int, uri, requestId, method string) (fn func(), newCtx context.Context, rateLimited bool)
}

// DebugEgressController - debug egress controller
type DebugEgressController struct{}

// Apply - function to be used by non Http egress traffic to apply a controller
func (e *DebugEgressController) Apply(ctx context.Context, statusCode func() int, uri, requestId, method string) (func(), context.Context, bool) {
	return nil, nil, false
}

// EgressApply - function to be used by non Http egress traffic to apply an controller
func Apply(ctx context.Context, statusCode func() int, uri, requestId, method string) (func(), context.Context, bool) {
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

func NewStatusCode(status **runtime.Status) func() int {
	return func() int {
		return int((*(status)).Code())
	}
}
