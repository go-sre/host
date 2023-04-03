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

	FalseValue = "false"
	TrueValue  = "true"

	BehaviorTimeout   = "timeout"
	BehaviorRetry     = "retry"
	BehaviorRateLimit = "rate-limit"
	BehaviorProxy     = "proxy"
	BehaviorFailover  = "failover"
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

func UpdateEnable(s State, values url.Values) error {
	if s == nil || values == nil {
		return errors.New("invalid argument: state or values is nil")
	}
	if values.Has(EnableKey) {
		v := values.Get(EnableKey)
		if v == TrueValue {
			s.Enable()
		} else {
			if v == FalseValue {
				s.Disable()
			} else {
				return errors.New(fmt.Sprintf("invalid argument: enable value is invalid : [%v]", v))
			}
		}
	}
	return nil
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
