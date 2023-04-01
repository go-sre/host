package controller

import (
	"errors"
	"net/url"
	"strconv"
)

type FailoverInvoke func(name string, failover bool)

// Failover - interface for failover
type Failover interface {
	State
	Actuator
	Invoke(failover bool)
}

type FailoverConfig struct {
	Enabled bool
	invoke  FailoverInvoke
}

var disabledFailover = newFailover("[disabled]", nil, NewFailoverConfig(false, nil))

func NewFailoverConfig(enabled bool, invoke FailoverInvoke) *FailoverConfig {
	return &FailoverConfig{Enabled: enabled, invoke: invoke}
}

type failover struct {
	table   *table
	name    string
	enabled bool
	invoke  FailoverInvoke
}

func cloneFailover(curr *failover) *failover {
	t := new(failover)
	*t = *curr
	return t
}

func newFailover(name string, table *table, config *FailoverConfig) *failover {
	t := new(failover)
	t.table = table
	t.name = name
	if config != nil {
		t.invoke = config.invoke
	}
	t.enabled = false
	return t
}

func (f *failover) validate() error {
	if f.invoke == nil {
		return errors.New("invalid configuration: Failover FailureInvoke function is nil")
	}
	return nil
}

func failoverState(m map[string]string, f *failover) {
	if f != nil {
		m[FailoverName] = strconv.FormatBool(f.IsEnabled())
	} else {
		m[FailoverName] = strconv.FormatBool(false)
	}
}

func (f *failover) Signal(values url.Values) error {
	UpdateEnable(f, values)
	return nil
}

func (f *failover) IsEnabled() bool { return f.enabled }

func (f *failover) Enable() {
	if f.IsEnabled() {
		return
	}
	f.enableFailover(true)
}

func (f *failover) Disable() {
	if !f.IsEnabled() {
		return
	}
	f.enableFailover(false)
}

func (f *failover) Invoke(failover bool) {
	if f.invoke == nil {
		return
	}
	f.invoke(f.name, failover)
}

func (f *failover) enableFailover(enabled bool) {
	if f.table == nil {
		return
	}
	f.table.mu.Lock()
	defer f.table.mu.Unlock()
	if ctrl, ok := f.table.controllers[f.name]; ok {
		c := cloneFailover(ctrl.failover)
		c.enabled = enabled
		f.table.update(f.name, cloneController[*failover](ctrl, c))
	}
}

/*
func (t *table) setFailoverInvoke(name string, fn FailoverInvoke, enable bool) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if ctrl, ok := t.controllers[name]; ok {
		fc := cloneFailover(ctrl.failover)
		fc.enabled = true
		fc.invoke = fn
		fc.enabled = enable
		t.update(name, cloneController[*failover](ctrl, fc))
	}
}


*/
