package accessdata

import (
	"errors"
	"fmt"
)

var operators = map[string]*Operator{
	TrafficOperator:        {"traffic", TrafficOperator},
	StartTimeOperator:      {"start-time", StartTimeOperator},
	DurationOperator:       {"duration-ms", DurationOperator},
	DurationStringOperator: {"duration", DurationStringOperator},

	OriginRegionOperator:     {"region", OriginRegionOperator},
	OriginZoneOperator:       {"zone", OriginZoneOperator},
	OriginSubZoneOperator:    {"sub-zone", OriginSubZoneOperator},
	OriginServiceOperator:    {"service", OriginServiceOperator},
	OriginInstanceIdOperator: {"instance-id", OriginInstanceIdOperator},

	// Route
	RouteNameOperator:       {"route-name", RouteNameOperator},
	TimeoutDurationOperator: {"timeout-ms", TimeoutDurationOperator},
	RateLimitOperator:       {"rate-limit", RateLimitOperator},
	RateBurstOperator:       {"rate-burst", RateBurstOperator},
	RateThresholdOperator:   {"rate-threshold", RateThresholdOperator},

	RetryOperator: {"retry", RetryOperator},
	//RetryRateLimitOperator:  {"retry_rate_limit", RetryRateLimitOperator},
	//RetryRateBurstOperator:  {"retry_rate_burst", RetryRateBurstOperator},
	ProxyOperator:          {"proxy", ProxyOperator},
	ProxyThresholdOperator: {"proxy-threshold", ProxyThresholdOperator},

	// Response
	ResponseStatusCodeOperator:    {"status-code", ResponseStatusCodeOperator},
	ResponseBytesReceivedOperator: {"bytes-received", ResponseBytesReceivedOperator},
	ResponseBytesSentOperator:     {"bytes-sent", ResponseBytesSentOperator},
	StatusFlagsOperator:           {"status-flags", StatusFlagsOperator},
	//UpstreamHostOperator:  {"upstream_host", UpstreamHostOperator},

	// Request
	RequestProtocolOperator: {"protocol", RequestProtocolOperator},
	RequestUrlOperator:      {"url", RequestUrlOperator},
	RequestMethodOperator:   {"method", RequestMethodOperator},
	RequestPathOperator:     {"path", RequestPathOperator},
	RequestHostOperator:     {"host", RequestHostOperator},

	RequestIdOperator:           {"request-id", RequestIdOperator},
	RequestFromRouteOperator:    {"request-from", RequestFromRouteOperator},
	RequestUserAgentOperator:    {"user-agent", RequestUserAgentOperator},
	RequestAuthorityOperator:    {"authority", RequestAuthorityOperator},
	RequestForwardedForOperator: {"forwarded", RequestForwardedForOperator},

	// gRPC
	GRPCStatusOperator:       {"grpc-status", GRPCStatusOperator},
	GRPCStatusNumberOperator: {"grpc-number", GRPCStatusNumberOperator},
}

func CreateOperators(operators []string) ([]Operator, error) {
	var items []Operator
	for _, op := range operators {
		items = append(items, Operator{
			Name:  "",
			Value: op,
		})
	}
	return InitOperators(items)
}

func InitOperators(operators []Operator) ([]Operator, error) {
	var items []Operator

	if len(operators) == 0 {
		return nil, errors.New("invalid configuration: configuration slice is empty")
	}
	dup := make(map[string]string)
	for _, op := range operators {
		op2, err := createOperator(op)
		if err != nil {
			return nil, err
		}
		if _, ok := dup[op2.Name]; ok {
			return nil, errors.New(fmt.Sprintf("invalid operator: name is a duplicate [%v]", op2.Name))
		}
		dup[op2.Name] = op2.Name
		items = append(items, op2)
	}
	return items, nil
}

func createOperator(op Operator) (Operator, error) {
	if IsEmpty(op.Value) {
		return Operator{}, errors.New(fmt.Sprintf("invalid operator: value is empty %v", op.Name))
	}
	if IsDirectOperator(op) {
		if IsEmpty(op.Name) {
			return Operator{}, errors.New(fmt.Sprintf("invalid operator: name is empty [%v]", op.Value))
		}
		return Operator{Name: op.Name, Value: op.Value}, nil
	}
	if op2, ok := operators[op.Value]; ok {
		newOp := Operator{Name: op2.Name, Value: op.Value}
		if !IsEmpty(op.Name) {
			newOp.Name = op.Name
		}
		return newOp, nil
	}
	if IsRequestOperator(op) {
		return Operator{Name: RequestOperatorHeaderName(op), Value: op.Value}, nil
	}
	return Operator{}, errors.New(fmt.Sprintf("invalid operator: value not found or invalid %v", op.Value))
}
