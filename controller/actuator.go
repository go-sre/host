package controller

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

const (
	EnableOpCode = "enable"
	SetOpCode    = "set"
	IncOpCode    = "inc"
	DecOpCode    = "dec"
)

type Actuator interface {
	Signal(opCode, value string) error
}

func boolValue(value string) (bool, error) {
	if len(value) == 0 {
		return false, errors.New("value is empty")
	}
	if value == "0" {
		return false, nil
	}
	if value == "1" {
		return true, nil
	}
	return false, errors.New(fmt.Sprintf("value is invalid: %v", value))
}

func intValue(value string) (int, error) {
	if len(value) == 0 {
		return -1, errors.New("value is empty")
	}
	return strconv.Atoi(value)
}

func urlValue(value string) (*url.URL, error) {
	if len(value) == 0 {
		return nil, errors.New("value is empty")
	}
	return url.Parse(value)
}
