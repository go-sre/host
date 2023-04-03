package middleware

import (
	"github.com/felixge/httpsnoop"
	"github.com/go-sre/host/controller"
	"net/http"
	"time"
)

// ControllerHttpHostMetricsHandler - handler that applies controller controllers
func ControllerHttpHostMetricsHandler(appHandler http.Handler, msg string) http.Handler {
	wrappedH := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UTC()
		ctrl := controller.IngressTable().Host()
		var m httpsnoop.Metrics

		if rlc := ctrl.RateLimiter(); rlc.IsEnabled() && !rlc.Allow() {
			ctrl.LogHttpIngress(start, time.Since(start), r, rlc.StatusCode(), 0, controller.RateLimitFlag)
			return
		}
		ctrl = controller.IngressTable().LookupHttp(r)
		if toc := ctrl.Timeout(); toc.IsEnabled() {
			m = httpsnoop.CaptureMetrics(http.TimeoutHandler(appHandler, toc.Duration(), msg), w, r)
		} else {
			m = httpsnoop.CaptureMetrics(appHandler, w, r)
		}
		//	log.Printf("%s %s (code=%d dt=%s written=%d)", r.Method, r.URL, m.Code, m.Duration, m.Written)
		ctrl.LogHttpIngress(start, time.Since(start), r, m.Code, m.Written, "")
	})
	return wrappedH
}
