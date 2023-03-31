package controller

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const (
	TrafficKey    = "traffic"
	ControllerKey = "controller"
	PatternKey    = "pattern"
	MethodKey     = "method"

	ActionKey = "action"

	EnableAction = "enable"
	SetAction    = "set"
	IncAction    = "inc"
	DecAction    = "dec"
)

type Actuator interface {
	Signal(values url.Values) error
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

func parseTest(req *http.Request) {
	req.URL.Query()
}

func UpdateEnable(s State, values url.Values) error {
	if s == nil || values == nil {
		return errors.New("invalid argument: state or values is nil")
	}
	if values.Has("enable") {
		v := values.Get("enable")
		if v == "true" {
			s.Enable()
		} else {
			if v == "false" {
				s.Disable()
			} else {
				return errors.New(fmt.Sprintf("invalid argument: enable value is invalid : [%v]", v))
			}
		}
	}
	return nil
}