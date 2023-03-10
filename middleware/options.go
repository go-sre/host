package middleware

import (
	"github.com/gotemplates/host/accessdata"
	"github.com/gotemplates/host/accesslog"
)

// SetAccessLogFn - allows setting an application configured logging function
func SetAccessLogFn(fn func(e *accessdata.Entry)) {
	if fn != nil {
		logFn = fn
	}
}

var logFn = defaultLogFn

var defaultLogFn = func(e *accessdata.Entry) {
	accesslog.Write[accesslog.LogOutputHandler, accessdata.JsonFormatter](e)
}
