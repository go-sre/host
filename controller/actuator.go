package controller

import (
	"errors"
	"fmt"
	"golang.org/x/time/rate"
	"net/url"
	"strconv"
)

const (
	TrafficKey  = "traffic"
	RouteKey    = "route"
	BehaviorKey = "behavior"

	RateLimitKey = "limit"
	RateBurstKey = "burst"
	DurationKey  = "duration"
	EnableKey    = "enable"
	PatternKey   = "pattern"
	WaitKey      = "wait"
	PercentKey   = "pct"

	FalseValue = "false"
	TrueValue  = "true"

	TimeoutBehavior   = "timeout"
	RetryBehavior     = "retry"
	RateLimitBehavior = "rate-limit"
	ProxyBehavior     = "proxy"

	NilPercentageValue = float64(-1)
)

type Actuator interface {
	Signal(values url.Values) error
}

func IsEgressTraffic(values url.Values) bool {
	if values == nil {
		return false
	}
	if values.Get(TrafficKey) == EgressTraffic {
		return true
	}
	return false
}

func IsIngressTraffic(values url.Values) bool {
	if values == nil {
		return false
	}
	if values.Get(TrafficKey) == IngressTraffic {
		return true
	}
	return false
}

func IsDisable(values url.Values) bool {
	if values == nil {
		return false
	}
	if values.Get(EnableKey) == FalseValue {
		return true
	}
	return false
}

func IsEnable(values url.Values) bool {
	return !IsDisable(values)
}

func NewValues(key, value string) url.Values {
	values := url.Values{}
	values.Set(key, value)
	return values
}

func UpdateEnable(s State, values url.Values) (stateChange bool, err error) {
	if s == nil {
		return false, errors.New("invalid argument: state is nil")
	}
	if !values.Has(EnableKey) {
		return false, nil
	}
	v := values.Get(EnableKey)
	if v == TrueValue && !s.IsEnabled() {
		s.Enable()
		return true, nil
	}
	if v == FalseValue && s.IsEnabled() {
		s.Disable()
		return true, nil
	}
	return false, errors.New(fmt.Sprintf("invalid argument: enable value is invalid : [%v]", v))
}

func ParseLimitAndBurst(values url.Values) (rate.Limit, int, error) {
	var limit = rate.Limit(-1)
	var burst = -1

	if values == nil {
		return limit, burst, nil
	}
	if values.Has(RateLimitKey) {
		s := values.Get(RateLimitKey)
		if len(s) > 0 {
			temp, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return limit, burst, err
			}
			if temp <= 0 {
				return limit, burst, errors.New(fmt.Sprintf("invalid argument: limit value is <= 0 [%v]", temp))
			}
			limit = rate.Limit(temp)
		}
	}
	if values.Has(RateBurstKey) {
		s := values.Get(RateBurstKey)
		if len(s) > 0 {
			temp, err := strconv.Atoi(s)
			if err != nil {
				return limit, burst, err
			}
			if temp <= 0 {
				return limit, burst, errors.New(fmt.Sprintf("invalid argument: burst value is <= 0 [%v]", temp))
			}
			burst = temp
		}
	}
	return limit, burst, nil
}

func ParsePercentage(values url.Values) float64 {
	if values == nil {
		return NilPercentageValue
	}
	if values.Has(PercentKey) {
		s := values.Get(PercentKey)
		if len(s) > 0 {
			temp, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return NilPercentageValue
			}
			if temp <= 0 || temp > 1000 {
				return NilPercentageValue
			}
			return temp / float64(100)
		}
	}
	return NilPercentageValue
}
