package controller

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Timeout - interface for timeouts
type Timeout interface {
	State
	Actuator
	StatusCode() int
	Duration() time.Duration
}

type TimeoutConfig struct {
	Enabled    bool
	StatusCode int
	Duration   time.Duration
}

var nilTimeout = newTimeout(NilBehaviorName, nil, NewTimeoutConfig(false, 0, 1))

func NewTimeoutConfig(enabled bool, statusCode int, duration time.Duration) *TimeoutConfig {
	if statusCode <= 0 {
		statusCode = http.StatusGatewayTimeout
	}
	return &TimeoutConfig{Enabled: enabled, StatusCode: statusCode, Duration: duration}
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
	if t.config.Duration < 0 {
		return errors.New(fmt.Sprintf("invalid configuration: Timeout duration is < 0 [%v]", t.name))
	}
	return nil
}

func timeoutState(t *timeout) int {
	var val int = -1

	if t != nil && t.IsEnabled() {
		val = int(t.Duration() / time.Millisecond)
	}
	return val
}

func (t *timeout) Signal(values url.Values) error {
	if t.IsNil() {
		return errors.New("invalid signal: timeout is not configured")
	}
	if values == nil {
		return errors.New("invalid argument: values are nil for timeout signal")
	}
	UpdateEnable(t, values)
	if values.Has(DurationKey) {
		duration, err := ParseDuration(values.Get(DurationKey))
		if err != nil {
			return err
		}
		if duration <= 0 {
			return errors.New("invalid configuration: timeout duration is <= 0")
		}
		if duration != t.Duration() {
			t.setTimeout(duration)
		}
	}
	pct := ParsePercentage(values)
	if pct != NilPercentageValue {
		val := t.Duration() + time.Duration(pct*float64(t.Duration()))
		t.setTimeout(val)
	}
	return nil
}

func (t *timeout) IsEnabled() bool { return t.config.Enabled }

func (t *timeout) IsNil() bool { return t.name == NilBehaviorName }

func (t *timeout) Enable() {
	if t.IsEnabled() {
		return
	}
	t.enableTimeout(true)
}

func (t *timeout) Disable() {
	if !t.IsEnabled() {
		return
	}
	t.enableTimeout(false)
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

func (t *timeout) enableTimeout(enable bool) {
	if t.table == nil || t.IsNil() {
		return
	}
	t.table.mu.Lock()
	defer t.table.mu.Unlock()
	if ctrl, ok := t.table.controllers[t.name]; ok {
		c := cloneTimeout(ctrl.timeout)
		c.config.Enabled = enable
		t.table.update(t.name, cloneController[*timeout](ctrl, c))
	}
}

func (t *timeout) setTimeout(duration time.Duration) {
	if t.table == nil || t.IsNil() || t.config.Duration == duration {
		return
	}
	t.table.mu.Lock()
	defer t.table.mu.Unlock()
	if ctrl, ok := t.table.controllers[t.name]; ok {
		c := cloneTimeout(ctrl.timeout)
		c.config.Duration = duration
		t.table.update(t.name, cloneController[*timeout](ctrl, c))
	}
}
