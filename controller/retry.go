package controller

import (
	"errors"
	"fmt"
	"golang.org/x/time/rate"
	"math/rand"
	"net/url"
	"time"
)

// https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
// https://github.com/keikoproj/inverse-exp-backoff

// Retry - interface for retries
type Retry interface {
	State
	Actuator
	IsValidStatusCode(statusCode int) bool
	IsRetryable(statusCode int) (ok bool, status string)
	Limit() rate.Limit
	Burst() int
	Wait() time.Duration
}

type RetryConfig struct {
	Enabled     bool
	Limit       rate.Limit
	Burst       int
	Wait        time.Duration
	StatusCodes []int
}

var nilRetry = newRetry(NilBehaviorName, nil, NewRetryConfig(false, 0, 0, 0, nil))

func NewRetryConfig(enabled bool, limit rate.Limit, burst int, wait time.Duration, validCodes []int) *RetryConfig {
	c := new(RetryConfig)
	c.Wait = wait
	c.Limit = limit
	c.Burst = burst
	c.StatusCodes = validCodes
	c.Enabled = enabled
	return c
}

type retry struct {
	name        string
	table       *table
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
	t.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	if config != nil {
		t.config = *config
	}
	t.rateLimiter = rate.NewLimiter(t.config.Limit, t.config.Burst)
	return t
}

func (r *retry) validate() error {
	if len(r.config.StatusCodes) == 0 {
		return errors.New(fmt.Sprintf("invalid configuration: retry status codes are empty [%v]", r.name))
	}
	if r.config.Limit < 0 {
		return errors.New(fmt.Sprintf("invalid configuration: retry limit is < 0 [%v]", r.name))
	}
	if r.config.Burst < 0 {
		return errors.New(fmt.Sprintf("invalid configuration: retry burst is < 0 [%v]", r.name))
	}
	if r.config.Wait < 0 {
		return errors.New(fmt.Sprintf("invalid configuration: wait duration is < 0 [%v]", r.name))
	}
	return nil
}

func retryState(r *retry) (rate.Limit, int) {
	var limit rate.Limit = -1
	var burst = -1
	//var name = "false"
	if r != nil && r.IsEnabled() {
		//	name = strconv.FormatBool(retried)
		limit = r.config.Limit
		if limit == rate.Inf {
			limit = RateLimitInfValue
		}
		burst = r.config.Burst
	}
	return limit, burst

	//if m == nil {
	//	m = make(map[string]string, 16)
	//}
	//m[RetryName] = name
	//m[RetryRateLimitName] = fmt.Sprintf("%v", limit)
	//m[RetryRateBurstName] = strconv.Itoa(burst)
	//return m

}

func (r *retry) IsEnabled() bool { return r.config.Enabled }

func (r *retry) IsNil() bool { return r.name == NilBehaviorName }

func (r *retry) Enable() {
	if r.IsEnabled() {
		return
	}
	r.enableRetry(true)
}

func (r *retry) Disable() {
	if !r.IsEnabled() {
		return
	}
	r.enableRetry(false)
}

func (r *retry) Signal(values url.Values) error {
	if r.IsNil() {
		return errors.New("invalid signal: retry is not configured")
	}
	if values == nil {
		return errors.New("invalid argument: values are nil for retry signal")
	}
	UpdateEnable(r, values)
	limit, burst, err := ParseLimitAndBurst(values)
	if err != nil {
		return err
	}
	if limit != -1 || burst != -1 {
		if limit == -1 {
			limit = r.config.Limit
		}
		if burst == -1 {
			burst = r.config.Burst
		}
		if r.config.Limit != limit || r.config.Burst != burst {
			r.setRetryRateLimiter(limit, burst)
		}
	}
	if values.Has(WaitKey) {
		duration, err1 := ParseDuration(values.Get(WaitKey))
		if err1 != nil {
			return err1
		}
		if duration < 0 {
			return errors.New("invalid configuration: wait duration is < 0")
		}
		if duration != r.Wait() {
			r.setWait(duration)
		}
	}

	return nil
}

func (r *retry) IsValidStatusCode(statusCode int) bool {
	for _, code := range r.config.StatusCodes {
		if code == statusCode {
			return true
		}
	}
	return false
}

func (r *retry) IsRetryable(statusCode int) (bool, string) {
	if !r.rateLimiter.Allow() {
		return false, RateLimitFlag
	}
	for _, code := range r.config.StatusCodes {
		if code == statusCode {
			if r.config.Wait == 0 {
				return true, ""
			}
			jitter := time.Duration(r.rand.Int31n(1000))
			time.Sleep(r.config.Wait + jitter)
			return true, ""
		}
	}
	return false, ""
}

func (r *retry) Limit() rate.Limit {
	return r.config.Limit
}

func (r *retry) Burst() int {
	return r.config.Burst
}

func (r *retry) Wait() time.Duration {
	return r.config.Wait
}

/*
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

*/

//func (r *retry) LimitAndBurst() (rate.Limit, int) {
//	return r.config.Limit, r.config.Burst
//}

/*
func (r *retry) SetLimit(limit rate.Limit) {
	if r.config.Limit == limit {
		return
	}
	r.setRetryRateLimit(limit)
}

func (r *retry) SetBurst(burst int) {
	if r.config.Burst == burst {
		return
	}
	r.setRetryRateBurst(burst)
}


*/

func (r *retry) enableRetry(enable bool) {
	if r.table == nil || r.IsNil() {
		return
	}
	r.table.mu.Lock()
	defer r.table.mu.Unlock()
	if ctrl, ok := r.table.controllers[r.name]; ok {
		c := cloneRetry(ctrl.retry)
		c.config.Enabled = enable
		r.table.update(r.name, cloneController[*retry](ctrl, c))
	}
}

func (r *retry) setRetryRateLimiter(limit rate.Limit, burst int) {
	if r.table == nil || r.IsNil() {
		return
	}
	r.table.mu.Lock()
	defer r.table.mu.Unlock()
	if ctrl, ok := r.table.controllers[r.name]; ok {
		c := cloneRetry(ctrl.retry)
		c.config.Limit = limit
		c.config.Burst = burst
		// Not cloning the limiter as an old reference will not cause stale data when logging
		c.rateLimiter = rate.NewLimiter(limit, burst)
		r.table.update(r.name, cloneController[*retry](ctrl, c))
	}
}

func (r *retry) setWait(duration time.Duration) {
	if r.table == nil || r.IsNil() || r.config.Wait == duration {
		return
	}
	r.table.mu.Lock()
	defer r.table.mu.Unlock()
	if ctrl, ok := r.table.controllers[r.name]; ok {
		c := cloneRetry(ctrl.retry)
		c.config.Wait = duration
		r.table.update(r.name, cloneController[*retry](ctrl, c))
	}
}

/*
func (r *retry) setRetryRateBurst(burst int) {
	if r.table == nil {
		return
	}
	r.table.mu.Lock()
	defer r.table.mu.Unlock()
	if ctrl, ok := r.table.controllers[r.name]; ok {
		c := cloneRetry(ctrl.retry)
		c.config.Burst = burst
		// Not cloning the limiter as an old reference will not cause stale data when logging
		c.rateLimiter = rate.NewLimiter(ctrl.retry.config.Limit, burst)
		r.table.update(r.name, cloneController[*retry](ctrl, c))
	}
}


*/
