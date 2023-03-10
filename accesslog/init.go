package accesslog

import (
	"encoding/json"
	"errors"
	"github.com/gotemplates/host/accessdata"
)

// InitIngressOperators - allows configuration of access accesslog attributes for ingress traffic
func InitIngressOperators(config []accessdata.Operator) error {
	var err error
	ingressOperators, err = accessdata.InitOperators(config)
	return err
}

// InitEgressOperators - allows configuration of access accesslog attributes for egress traffic
func InitEgressOperators(config []accessdata.Operator) error {
	var err error
	egressOperators, err = accessdata.InitOperators(config)
	return err
}

// CreateEgressOperators - provides creation of egress operators
func CreateEgressOperators(read func() ([]byte, error)) error {
	if read == nil {
		return errors.New("invalid argument: ReadConfig function is nil")
	}
	buf, err0 := read()
	if err0 != nil {
		return err0
	}
	operators, err := ReadOperators(buf)
	if err != nil {
		return err
	}
	return InitEgressOperators(operators)
}

// CreateIngressOperators - provides creation of ingress operators from a []byte
func CreateIngressOperators(read func() ([]byte, error)) error {
	if read == nil {
		return errors.New("invalid argument: ReadConfig function is nil")
	}
	buf, err0 := read()
	if err0 != nil {
		return err0
	}
	operators, err := ReadOperators(buf)
	if err != nil {
		return err
	}
	return InitIngressOperators(operators)
}

// ReadOperators - read the operators from a []byte
func ReadOperators(buf []byte) ([]accessdata.Operator, error) {
	var operators []accessdata.Operator

	if buf == nil {
		return nil, errors.New("invalid argument: buffer is nil")
	}
	err1 := json.Unmarshal(buf, &operators)
	return operators, err1
}
