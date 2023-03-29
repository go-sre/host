package controller

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Timeout - interface for timeouts
type Timeout interface {
	State
	Actuator
	StatusCode() int
	Duration() time.Duration
	SetTimeout(timeout time.Duration)
}

type TimeoutConfig struct {
	Disabled   bool
	StatusCode int
	Duration   time.Duration
}

func NewTimeoutConfig(enabled bool, statusCode int, duration time.Duration) *TimeoutConfig {
	if statusCode <= 0 {
		statusCode = http.StatusGatewayTimeout
	}
	return &TimeoutConfig{Disabled: !enabled, StatusCode: statusCode, Duration: duration}
}

type timeout struct {
	table  *table
	name   string
	config TimeoutConfig
}

func cloneTimeout(curr *timeout) *timeout {
	t := new(timeout)
	*t = *curr
	return t
}

func newTimeout(name string, table *table, config *TimeoutConfig) *timeout {
	t := new(timeout)
	t.table = table
	t.name = name
	if config != nil {
		t.config = *config
	}
	return t
}

func (t *timeout) validate() error {
	if t.config.Duration <= 0 {
		return errors.New("invalid configuration: Timeout duration is <= 0")
	}
	return nil
}

func timeoutState(m map[string]string, t *timeout) {
	var val int64 = -1
	//var statusCode = -1
	if t != nil {
		val = int64(t.Duration() / time.Millisecond)
		//	statusCode = t.StatusCode()
	}
	m[TimeoutName] = strconv.Itoa(int(val))
}

func (t *timeout) Signal(_ url.Values) error { return errors.New("timeout Actuator not supported") }

func (t *timeout) IsEnabled() bool { return !t.config.Disabled }

func (t *timeout) Enable() {
	if t.IsEnabled() {
		return
	}
	t.config.Disabled = false
	// Need to update table
}

func (t *timeout) Disable() {
	if !t.IsEnabled() {
		return
	}
	t.config.Disabled = true
	// Need to update table
}

func (t *timeout) StatusCode() int {
	return t.config.StatusCode
}

func (t *timeout) Duration() time.Duration {
	if t.config.Duration <= 0 {
		return 0
	}
	return t.config.Duration
}

func (t *timeout) SetTimeout(duration time.Duration) {
	if t.config.Duration == duration || duration <= 0 {
		return
	}
	t.table.setTimeout(t.name, duration)
}
