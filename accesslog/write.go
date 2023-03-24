package accesslog

import (
	"fmt"
	"github.com/go-sre/host/accessdata"
)

const (
	errorName     = "error"
	errorNilEntry = "access data entry is nil"
	errorEmptyFmt = "%v accesslog entries are empty"
)

var ingressOperators []accessdata.Operator
var egressOperators []accessdata.Operator

// Write - templated function handling writing the access data utilizing the OutputHandler and Formatter
func Write[O OutputHandler, F accessdata.Formatter](entry *accessdata.Entry) {
	var o O
	var f F
	if entry == nil {
		o.Write([]accessdata.Operator{{errorName, errorNilEntry}}, accessdata.NewEmptyEntry(), f)
		return
	}
	var operators []accessdata.Operator
	switch entry.Traffic {
	case accessdata.IngressTraffic, accessdata.PingTraffic:
		if entry.Traffic == accessdata.IngressTraffic && !opt.ingress {
			return
		}
		if entry.Traffic == accessdata.PingTraffic && !opt.ping {
			return
		}
		operators = ingressOperators
	case accessdata.EgressTraffic:
		if !opt.egress {
			return
		}
		operators = egressOperators
	}
	if len(operators) == 0 {
		operators = emptyOperators(entry)
	}
	o.Write(operators, entry, f)
}

func emptyOperators(entry *accessdata.Entry) []accessdata.Operator {
	return []accessdata.Operator{{errorName, fmt.Sprintf(errorEmptyFmt, entry.Traffic)}}
}
