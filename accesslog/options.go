package accesslog

// SetIngressLogStatus - enable/disable ingress logging
func SetIngressLogStatus(enabled bool) {
	opt.ingress = enabled
}

// SetEgressLogStatus - enable/disable egress logging
func SetEgressLogStatus(enabled bool) {
	opt.egress = enabled
}

// SetPingLogStatus - enable/disable ping logging
func SetPingLogStatus(enabled bool) {
	opt.ping = enabled
}

type options struct {
	ingress bool
	egress  bool
	ping    bool
}

var opt options

func init() {
	opt.ingress = true
	opt.egress = true
	opt.ping = true
}
