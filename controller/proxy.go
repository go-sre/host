package controller

import (
	"errors"
	"net/url"
	"strconv"
)

// Proxy - interface for proxy
type Proxy interface {
	IsEnabled() bool
	Enable()
	Disable()
	SetPattern(pattern string)
	BuildUrl(uri *url.URL) *url.URL
}

type ProxyConfig struct {
	Enabled bool
	Pattern string
}

func NewProxyConfig(pattern string) *ProxyConfig {
	return &ProxyConfig{Enabled: false, Pattern: pattern}
}

type proxy struct {
	table   *table
	name    string
	enabled bool
	pattern string
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
		t.pattern = config.Pattern
	}
	t.enabled = false
	return t
}

func (p *proxy) validate() error {
	if len(p.pattern) == 0 && p.enabled {
		return errors.New("invalid configuration: Proxy pattern is empty for enabled proxy")
	}
	return nil
}

func proxyState(m map[string]string, p *proxy) {
	if p == nil {
		m[ProxyName] = ""
	} else {
		m[ProxyName] = strconv.FormatBool(p.IsEnabled())
	}
}

func (p *proxy) IsEnabled() bool { return p.enabled }

func (p *proxy) Disable() {
	if !p.IsEnabled() {
		return
	}
	p.table.enableProxy(p.name, false)
}

func (p *proxy) Enable() {
	if p.IsEnabled() {
		return
	}
	p.table.enableProxy(p.name, true)
}

func (p *proxy) SetPattern(pattern string) {
	if len(pattern) != 0 {
		p.table.setProxyPattern(p.name, pattern, false)
	}
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
