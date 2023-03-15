# host



## accessdata 

[AccessData][datapkg] provides the Entry type, which contains all of the data needed for access logging. Also provided are functions and types that define command operators which 
allow the extraction and formatting of Entry data. The formatting of Entry data is implemented as a template parameter: 
~~~
// Formatter - template parameter for formatting
type Formatter interface {
	Format(items []Operator, data *Entry) string
}
~~~
Configurable items, specific to a package, are defined in an options.go file.

## accesslog

[AccessLog][logpkg] encompasses access logging functionality. Seperate operators, and runtime initialization of those operators, are provided for ingress and egress traffic. An output template parameter allows redirection of the access logging: 
~~~
// OutputHandler - template parameter for log output
type OutputHandler interface {
	Write(items []accessdata.Operator, data *accessdata.Entry, formatter accessdata.Formatter)
}
~~~
The log.Write function is a templated function, allowing for selection of output and formatting:
~~~
func Write[O OutputHandler, F accessdata.Formatter](entry *accessdata.Entry) {
    // implementation details
}
~~~

## controller

[Controller][controllerpkg] provides resiliency through the implementation of configurable timeouts, rate limiting, retries, failover, and proxy controllers.
The controllers can be applied to any ingress or egress http traffic, and support initialization through external configuration files. All attributes 
related to the application of the controllers to traffic are logged via AccessLog. Non-http calls, like database client calls, can also 
be configured for resiliency.

## messaging
[Messaging][messagingpkg] provides a way for a hosting process to communicate with packages. Packages that register themselves can then be started and pinged by the 
host via the templated functions:
~~~
// Ping - templated function to "ping" a registered resource
func Ping[E template.ErrorHandler](ctx context.Context, uri string) (status *runtime.Status) {
    // Implementation details
}

// Startup - templated function to startup all registered resources.
func Startup[E template.ErrorHandler, O template.OutputHandler](duration time.Duration, content ContentMap) (status *runtime.Status) {
    // Implementation details
}
~~~



## middleware

[Middleware][middlewarepkg] provides implementations of a http.Handler and http.RoundTripper that support ingress and egress logging. Options
available allow configuring a logging function.

Ingress logging implementations: 

~~~
// AccessHttpHostMetricsHandler - http handler that captures metrics about an ingress request, also logs an access entry.
func AccessHttpHostMetricsHandler(appHandler http.Handler, msg string) http.Handler {
    // implementation details
}

// ControllerHttpHostMetricsHandler - handler that applies controller controllers
func ControllerHttpHostMetricsHandler(appHandler http.Handler, msg string) http.Handler {
    // implementation details
}
~~~

Egress logging implementation:

~~~

// AccessLogWrapTransport - provides a RoundTrip wrapper that applies controller controllers
func AccessLogWrapTransport(client *http.Client) {
   // implementation details
}

// RoundTrip - implementation of the RoundTrip interface for a transport, also logs an access entry
func (w *accessWrapper) RoundTrip(req *http.Request) (*http.Response, error) {
   // implementation details
}

// ControllerWrapTransport - provides a RoundTrip wrapper that applies controller controllers
func ControllerWrapTransport(client *http.Client) {
   // implementation details
}

// RoundTrip - implementation of the RoundTrip interface for a transport, also logs an access entry
func (w *controllerWrapper) RoundTrip(req *http.Request) (*http.Response, error) {
   // implementation details
}

~~~

Configuration of a logging function is supported via an option, which can be used to change the default:

~~~

// AccessSetLogFn - allows setting an application configured logging function
func AccessSetLogFn(fn func(e *data.Entry)) {
    // implementation details
}

var defaultLogFn = func(e *accessdata.Entry) {
	accesslog.Write[accesslog.LogOutputHandler, accessdata.JsonFormatter](e)
}
~~~



[datapkg]: <https://pkg.go.dev/github.com/gotemplates/host/accessdata>
[logpkg]: <https://pkg.go.dev/github.com/gotemplates/host/accesslog>
[controllerpkg]: <https://pkg.go.dev/github.com/gotemplates/host/controller>
[messagingpkg]: <https://pkg.go.dev/github.com/gotemplates/host/messaging>
[middlewarepkg]: <https://pkg.go.dev/github.com/gotemplates/host/middleware>

