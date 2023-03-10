package controller

import (
	"net/url"
	"strings"
)

// ParseUri - urn safe Uri parser
func ParseUri(uri string) (scheme, host, path string) {
	if uri == "" {
		return
	}
	u, err := url.Parse(uri)
	if err != nil {
		return err.Error(), "", ""
	}
	if u.Scheme == "urn" && u.Host == "" {
		t := strings.Split(u.Opaque, ":")
		if len(t) == 1 {
			return u.Scheme, t[0], ""
		}
		return u.Scheme, t[0], t[1]
	}
	return u.Scheme, u.Host, u.Path
}
