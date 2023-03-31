package controller

import (
	"errors"
	"fmt"
	"golang.org/x/time/rate"
	"net/http"
	"net/url"
	"strconv"
)

const (
	InfValue     = "-1"
	DefaultBurst = 1
)

// RateLimiter - interface for rate limiting
type RateLimiter interface {
	State
	Actuator
	Allow() bool
	StatusCode() int
	LimitAndBurst() (rate.Limit, int)
	SetLimit(limit rate.Limit)
	SetBurst(burst int)
}

type RateLimiterConfig struct {
	Enabled    bool
	StatusCode int
	Limit      rate.Limit
	Burst      int
}

var disabledRateLimiter = newRateLimiter("[disabled]", nil, NewRateLimiterConfig(false, 0, 0, 0))

func NewRateLimiterConfig(enabled bool, statusCode int, limit rate.Limit, burst int) *RateLimiterConfig {
	//validateLimiter(&limit, &burst)
	c := new(RateLimiterConfig)
	c.Limit = limit
	c.Burst = burst
	if statusCode <= 0 {
		statusCode = http.StatusTooManyRequests
	}
	c.StatusCode = statusCode
	c.Enabled = enabled
	return c
}

type rateLimiter struct {
	name        string
	table       *table
	config      RateLimiterConfig
	rateLimiter *rate.Limiter
}

func cloneRateLimiter(curr *rateLimiter) *rateLimiter {
	newLimiter := new(rateLimiter)
	*newLimiter = *curr
	return newLimiter
}

func newRateLimiter(name string, table *table, config *RateLimiterConfig) *rateLimiter {
	t := new(rateLimiter)
	t.name = name
	t.table = table
	t.config = RateLimiterConfig{Limit: rate.Inf, Burst: DefaultBurst}
	if config != nil {
		t.config = *config
	}
	t.rateLimiter = rate.NewLimiter(t.config.Limit, t.config.Burst)
	return t
}

func validateLimiter(max *rate.Limit, burst *int) {
	if max != nil && *max < 0 {
		*max = rate.Inf
	}
	if burst != nil && *burst < 0 {
		*burst = DefaultBurst
	}
}

func (r *rateLimiter) validate() error {
	if r.config.Limit < 0 {
		return errors.New(fmt.Sprintf("invalid configuration: RateLimiter limit is < 0"))
	}
	if r.config.Burst < 0 {
		return errors.New(fmt.Sprintf("invalid configuration: RateLimiter burst is < 0"))
	}
	return nil
}

func rateLimiterState(m map[string]string, r *rateLimiter) map[string]string {
	var limit rate.Limit = -1
	var burst = -1

	if r != nil {
		limit = r.config.Limit
		if limit == rate.Inf {
			limit = RateLimitInfValue
		}
		burst = r.config.Burst
	}
	if m == nil {
		m = make(map[string]string, 16)
	}
	m[RateLimitName] = fmt.Sprintf("%v", limit)
	m[RateBurstName] = strconv.Itoa(burst)
	return m
}

func (r *rateLimiter) IsEnabled() bool { return r.config.Enabled }

func (r *rateLimiter) Enable() {
	if r.IsEnabled() {
		return
	}
	r.enableRateLimiter(true)
}

func (r *rateLimiter) Disable() {
	if !r.IsEnabled() {
		return
	}
	r.enableRateLimiter(false)
}

func (r *rateLimiter) Signal(values url.Values) error {
	if values == nil {
		return nil
	}
	UpdateEnable(r, values)
	return nil
}

func (r *rateLimiter) Allow() bool {
	if r.config.Limit == rate.Inf {
		return true
	}
	return r.rateLimiter.Allow()
}

func (r *rateLimiter) StatusCode() int {
	return r.config.StatusCode
}

func (r *rateLimiter) LimitAndBurst() (rate.Limit, int) {
	return r.config.Limit, r.config.Burst
}

func (r *rateLimiter) SetLimit(limit rate.Limit) {
	if r.config.Limit == limit {
		return
	}
	r.setRateLimit(limit)
}

func (r *rateLimiter) SetBurst(burst int) {
	if r.config.Burst == burst {
		return
	}
	r.setRateBurst(burst)
}

/*
func (r *rateLimiter) SetRateLimiter(limit rate.Limit, burst int) {
	validateLimiter(&limit, &burst)
	if r.config.Limit == limit && r.config.Burst == burst {
		return
	}
	r.table.setRateLimiter(r.name, RateLimiterConfig{Limit: limit, Burst: burst})
}

func (r *rateLimiter) AdjustRateLimiter(percentage int) bool {
	newLimit, ok := limitAdjust(float64(r.config.Limit), percentage)
	if !ok {
		return false
	}
	newBurst, ok1 := burstAdjust(r.config.Burst, percentage)
	if !ok1 {
		return false
	}
	r.table.setRateLimiter(r.name, RateLimiterConfig{Limit: rate.Limit(newLimit), Burst: newBurst})
	return true
}

func limitAdjust(val float64, percentage int) (float64, bool) {
	change := (math.Abs(float64(percentage)) / 100.0) * val
	if change >= val {
		return val, false
	}
	if percentage > 0 {
		return val + change, true
	}
	return val - change, true
}

func burstAdjust(val int, percentage int) (int, bool) {
	floatChange := (math.Abs(float64(percentage)) / 100.0) * float64(val)
	change := int(math.Round(floatChange))
	if change == 0 || change >= val {
		return val, false
	}
	if percentage > 0 {
		return val + change, true
	}
	return val - change, true
}

*/

func (r *rateLimiter) enableRateLimiter(enabled bool) {
	if r.table == nil {
		return
	}
	r.table.mu.Lock()
	defer r.table.mu.Unlock()
	if ctrl, ok := r.table.controllers[r.name]; ok {
		c := cloneRateLimiter(ctrl.rateLimiter)
		c.config.Enabled = enabled
		r.table.update(r.name, cloneController[*rateLimiter](ctrl, c))
	}
}

func (r *rateLimiter) setRateLimit(limit rate.Limit) {
	if r.table == nil {
		return
	}
	r.table.mu.Lock()
	defer r.table.mu.Unlock()
	if ctrl, ok := r.table.controllers[r.name]; ok {
		c := cloneRateLimiter(ctrl.rateLimiter)
		c.config.Limit = limit
		// Not cloning the limiter as an old reference will not cause stale data when logging
		c.rateLimiter.SetLimit(limit)
		r.table.update(r.name, cloneController[*rateLimiter](ctrl, c))
	}
}

func (r *rateLimiter) setRateBurst(burst int) {
	if r.table == nil {
		return
	}
	r.table.mu.Lock()
	defer r.table.mu.Unlock()
	if ctrl, ok := r.table.controllers[r.name]; ok {
		c := cloneRateLimiter(ctrl.rateLimiter)
		c.config.Burst = burst
		// Not cloning the limiter as an old reference will not cause stale data when logging
		c.rateLimiter.SetBurst(burst)
		r.table.update(r.name, cloneController[*rateLimiter](ctrl, c))
	}
}

/*
func (t *table) setRateLimiter(name string, config RateLimiterConfig) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if ctrl, ok := t.controllers[name]; ok {
		c := cloneRateLimiter(ctrl.rateLimiter)
		c.config.Limit = config.Limit
		c.config.Burst = config.Burst
		c.rateLimiter = rate.NewLimiter(c.config.Limit, c.config.Burst)
		t.update(name, cloneController[*rateLimiter](ctrl, c))
	}
}

*/
