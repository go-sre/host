package controller

import (
	"errors"
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
	Pattern() string
	Headers() []Header
	BuildUrl(uri *url.URL) *url.URL
}

type ProxyConfig struct {
	Enabled bool
	Pattern string
	Headers []Header
}

var disabledProxy = newProxy("[disabled]", nil, NewProxyConfig(false, "", nil))

func NewProxyConfig(enabled bool, pattern string, headers []Header) *ProxyConfig {
	p := new(ProxyConfig)
	p.Enabled = enabled
	p.Pattern = pattern
	p.Headers = headers
	return p
}

type proxy struct {
	table   *table
	name    string
	enabled bool
	pattern string
	headers []Header
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
		t.enabled = config.Enabled
		t.pattern = config.Pattern
		t.headers = config.Headers
	}
	return t
}

func (p *proxy) validate() error {
	if p.enabled {
		return validatePattern(p.pattern)
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
	if values == nil {
		return nil
	}
	UpdateEnable(p, values)
	if values.Has(PatternKey) {
		v := values.Get(PatternKey)
		if v != p.pattern {
			return p.setPattern(v)
		}
	}
	return nil
}

func (p *proxy) IsEnabled() bool { return p.enabled }

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

func (p *proxy) Pattern() string {
	return p.pattern
}

func (p *proxy) Headers() []Header {
	return p.headers
}

func (p *proxy) BuildUrl(uri *url.URL) *url.URL {
	if uri == nil || len(p.pattern) == 0 {
		return uri
	}
	uri2, err := url.Parse(p.pattern)
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
		c.enabled = enabled
		p.table.update(p.name, cloneController[*proxy](ctrl, c))
	}
}

func (p *proxy) setPattern(pattern string) error {
	if p.table == nil || p.pattern == pattern {
		return nil
	}
	err := validatePattern(pattern)
	if err != nil {
		return err
	}
	p.table.mu.Lock()
	defer p.table.mu.Unlock()
	if ctrl, ok := p.table.controllers[p.name]; ok {
		fc := cloneProxy(ctrl.proxy)
		fc.pattern = pattern
		p.table.update(p.name, cloneController[*proxy](ctrl, fc))
	}
	return nil
}

func validatePattern(pattern string) error {
	if len(pattern) == 0 {
		return errors.New("invalid argument: proxy pattern is empty")
	}
	_, err := url.Parse(pattern)
	if err != nil {
		return err
	}
	return nil
}
