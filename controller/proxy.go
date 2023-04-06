package controller

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

type Header struct {
	Name  string
	Value string
}

// Proxy - interface for proxy
type Proxy interface {
	State
	Actuator
	Action() Actuator
	SetAction(action Actuator) error
	Pattern() string
	Headers() []Header
	BuildUrl(uri *url.URL) *url.URL
}

type ProxyConfig struct {
	Enabled bool
	Pattern string
	Headers []Header
	Action  Actuator
}

var nilProxy = newProxy(NilBehaviorName, nil, NewProxyConfig(false, "", nil, nil))

func NewProxyConfig(enabled bool, pattern string, headers []Header, action Actuator) *ProxyConfig {
	p := new(ProxyConfig)
	p.Enabled = enabled
	p.Pattern = pattern
	p.Headers = headers
	p.Action = action
	return p
}

type proxy struct {
	table  *table
	name   string
	config ProxyConfig
}

func cloneProxy(curr *proxy) *proxy {
	t := new(proxy)
	*t = *curr
	return t
}

func newProxy(name string, table *table, config *ProxyConfig) *proxy {
	t := new(proxy)
	t.table = table
	t.name = name
	if config != nil {
		t.config = *config
	}
	return t
}

func (p *proxy) validate() error {
	if p.config.Enabled {
		return p.validatePattern(p.config.Pattern)
	}
	return nil
}

func proxyState(m map[string]string, p *proxy) {
	if p != nil {
		m[ProxyName] = strconv.FormatBool(p.IsEnabled())
	} else {
		m[ProxyName] = strconv.FormatBool(false)
	}
}

func (p *proxy) Signal(values url.Values) error {
	if p.name == NilBehaviorName {
		return errors.New("invalid signal: proxy is not configured")
	}
	if values == nil {
		return errors.New("invalid argument: values are nil for proxy signal")
	}
	if IsDisable(values) {
		p.Disable()
	}
	if values.Has(PatternKey) {
		v := values.Get(PatternKey)
		if len(v) == 0 {
			return errors.New("invalid configuration: proxy pattern is empty")
		}
		if v != p.config.Pattern {
			err := p.setPattern(v)
			if err != nil {
				return err
			}
			if p.config.Action != nil {
				return p.config.Action.Signal(NewValues(PatternKey, v))
			}
		}
	}
	if IsEnable(values) {
		p.Enable()
	}
	return nil
}

func (p *proxy) IsEnabled() bool { return p.config.Enabled }

func (p *proxy) Enable() {
	if p.IsEnabled() {
		return
	}
	p.enableProxy(true)
}

func (p *proxy) Disable() {
	if !p.IsEnabled() {
		return
	}
	p.enableProxy(false)
}

func (p *proxy) Action() Actuator {
	return p.config.Action
}

func (p *proxy) SetAction(action Actuator) error {
	if action == nil {
		return errors.New("invalid configuration: Proxy action is nil")
	}
	p.setAction(action)
	return nil
}

func (p *proxy) Pattern() string {
	return p.config.Pattern
}

func (p *proxy) Headers() []Header {
	return p.config.Headers
}

func (p *proxy) BuildUrl(uri *url.URL) *url.URL {
	if uri == nil || len(p.config.Pattern) == 0 {
		return uri
	}
	uri2, err := url.Parse(p.config.Pattern)
	if err != nil {
		return uri
	}
	var newUri = uri2.Scheme + "://"
	if len(uri2.Host) > 0 {
		newUri += uri2.Host
	} else {
		newUri += uri.Host
	}
	if len(uri2.Path) > 0 {
		newUri += uri2.Path
	} else {
		newUri += uri.Path
	}
	if len(uri2.RawQuery) > 0 {
		newUri += "?" + uri2.RawQuery
	} else {
		if len(uri.RawQuery) > 0 {
			newUri += "?" + uri.RawQuery
		}
	}
	u, err1 := url.Parse(newUri)
	if err1 != nil {
		return uri
	}
	return u
}

func (p *proxy) enableProxy(enabled bool) {
	if p.table == nil {
		return
	}
	p.table.mu.Lock()
	defer p.table.mu.Unlock()
	if ctrl, ok := p.table.controllers[p.name]; ok {
		c := cloneProxy(ctrl.proxy)
		c.config.Enabled = enabled
		p.table.update(p.name, cloneController[*proxy](ctrl, c))
	}
}

func (p *proxy) setPattern(pattern string) error {
	if p.table == nil {
		return nil
	}
	err := p.validatePattern(pattern)
	if err != nil {
		return err
	}
	p.table.mu.Lock()
	defer p.table.mu.Unlock()
	if ctrl, ok := p.table.controllers[p.name]; ok {
		fc := cloneProxy(ctrl.proxy)
		fc.config.Pattern = pattern
		p.table.update(p.name, cloneController[*proxy](ctrl, fc))
	}
	return nil
}

func (p *proxy) setAction(action Actuator) {
	if p.table == nil {
		return
	}
	p.table.mu.Lock()
	defer p.table.mu.Unlock()
	if ctrl, ok := p.table.controllers[p.name]; ok {
		fc := cloneProxy(ctrl.proxy)
		fc.config.Action = action
		p.table.update(p.name, cloneController[*proxy](ctrl, fc))
	}
}

func (p *proxy) validatePattern(pattern string) error {
	if len(pattern) == 0 {
		return errors.New(fmt.Sprintf("invalid argument: proxy pattern is empty [%v]", p.name))
	}
	_, err := url.Parse(pattern)
	if err != nil {
		return errors.New(fmt.Sprintf("%v [%v]", err.Error(), p.name))
	}
	return nil
}
