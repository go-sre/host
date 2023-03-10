package controller

import (
	"errors"
	"fmt"
	"golang.org/x/time/rate"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
// https://github.com/keikoproj/inverse-exp-backoff

// Retry - interface for retries
type Retry interface {
	IsEnabled() bool
	Enable()
	Disable()
	IsRetryable(statusCode int) (ok bool, status string)
	SetRateLimiter(limit rate.Limit, burst int)
	AdjustRateLimiter(percentage int) bool
	LimitAndBurst() (rate.Limit, int)
}

type RetryConfig struct {
	Limit rate.Limit
	Burst int
	Wait  time.Duration
	Codes []int
}

func NewRetryConfig(validCodes []int, limit rate.Limit, burst int, wait time.Duration) *RetryConfig {
	c := new(RetryConfig)
	c.Wait = wait
	c.Limit = limit
	c.Burst = burst
	c.Codes = validCodes
	return c
}

type retry struct {
	name        string
	table       *table
	enabled     bool
	rand        *rand.Rand
	config      RetryConfig
	rateLimiter *rate.Limiter
}

func cloneRetry(curr *retry) *retry {
	t := new(retry)
	*t = *curr
	return t
}

func newRetry(name string, table *table, config *RetryConfig) *retry {
	t := new(retry)
	t.name = name
	t.table = table
	t.enabled = true
	t.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	if config != nil {
		t.config = *config
	}
	t.rateLimiter = rate.NewLimiter(t.config.Limit, t.config.Burst)
	return t
}

func (r *retry) validate() error {
	if len(r.config.Codes) == 0 {
		return errors.New("invalid configuration: Retry status codes are empty")
	}
	if r.config.Limit < 0 {
		return errors.New("invalid configuration: Retry limit is < 0")
	}
	if r.config.Burst < 0 {
		return errors.New("invalid configuration: Retry burst is < 0")
	}
	return nil
}

func retryState(m map[string]string, r *retry, retried bool) map[string]string {
	var limit rate.Limit = -1
	var burst = -1
	var name = ""
	if r != nil {
		name = strconv.FormatBool(retried)
		limit = r.config.Limit
		if limit == rate.Inf {
			limit = RateLimitInfValue
		}
		burst = r.config.Burst
	}
	if m == nil {
		m = make(map[string]string, 16)
	}
	m[RetryName] = name
	m[RetryRateLimitName] = fmt.Sprintf("%v", limit)
	m[RetryRateBurstName] = strconv.Itoa(burst)
	return m

}

func (r *retry) IsEnabled() bool { return r.enabled }

func (r *retry) Disable() {
	if !r.IsEnabled() {
		return
	}
	r.table.enableRetry(r.name, false)
}

func (r *retry) Enable() {
	if r.IsEnabled() {
		return
	}
	r.table.enableRetry(r.name, true)
}

func (r *retry) SetRateLimiter(limit rate.Limit, burst int) {
	if r.config.Limit == limit && r.config.Burst == burst {
		return
	}
	r.table.setRetryRateLimit(r.name, limit, burst)
}

func (r *retry) IsRetryable(statusCode int) (bool, string) {
	if !r.IsEnabled() {
		return false, NotEnabledFlag
	}
	if statusCode < http.StatusInternalServerError {
		return false, ""
	}
	if !r.rateLimiter.Allow() {
		return false, RateLimitFlag
	}
	for _, code := range r.config.Codes {
		if code == statusCode {
			jitter := time.Duration(r.rand.Int31n(1000))
			time.Sleep(r.config.Wait + jitter)
			return true, ""
		}
	}
	return false, ""
}

func (r *retry) AdjustRateLimiter(percentage int) bool {
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

func (r *retry) LimitAndBurst() (rate.Limit, int) {
	return r.config.Limit, r.config.Burst
}
