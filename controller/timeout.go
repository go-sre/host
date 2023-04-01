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
}

type TimeoutConfig struct {
	Enabled    bool
	StatusCode int
	Duration   time.Duration
}

var disabledTimeout = newTimeout("[disabled]", nil, NewTimeoutConfig(false, 0, -1))

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
	if t.config.Duration <= 0 {
		return errors.New("invalid configuration: Timeout duration is <= 0")
	}
	return nil
}

func timeoutState(m map[string]string, t *timeout) {
	var val int64 = -1

	if t != nil && t.IsEnabled() {
		val = int64(t.Duration() / time.Millisecond)
	}
	m[TimeoutName] = strconv.Itoa(int(val))
}

func (t *timeout) Signal(values url.Values) error {
	if values == nil {
		return nil
	}
	UpdateEnable(t, values)
	if values.Has(DurationKey) {
		v := values.Get(DurationKey)
		duration, err := ParseDuration(v)
		if err != nil {
			return err
		}
		t.setTimeout(duration)
	}
	return nil
}

func (t *timeout) IsEnabled() bool { return t.config.Enabled }

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
	if t.table == nil {
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
	if t.table == nil || t.config.Duration == duration || duration <= 0 {
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
