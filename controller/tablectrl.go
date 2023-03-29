package controller

import (
	"golang.org/x/time/rate"
	"time"
)

func (t *table) enableFailover(name string, enabled bool) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if ctrl, ok := t.controllers[name]; ok {
		c := cloneFailover(ctrl.failover)
		c.enabled = enabled
		t.update(name, cloneController[*failover](ctrl, c))
	}
}

func (t *table) enableProxy(name string, enabled bool) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if ctrl, ok := t.controllers[name]; ok {
		c := cloneProxy(ctrl.proxy)
		c.enabled = enabled
		t.update(name, cloneController[*proxy](ctrl, c))
	}
}

func (t *table) setProxyPattern(name string, pattern string, enable bool) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if ctrl, ok := t.controllers[name]; ok {
		fc := cloneProxy(ctrl.proxy)
		fc.pattern = pattern
		fc.enabled = enable
		t.update(name, cloneController[*proxy](ctrl, fc))
	}
}

/*
func (t *table) setFailoverInvoke(name string, fn FailoverInvoke, enable bool) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if ctrl, ok := t.controllers[name]; ok {
		fc := cloneFailover(ctrl.failover)
		fc.enabled = true
		fc.invoke = fn
		fc.enabled = enable
		t.update(name, cloneController[*failover](ctrl, fc))
	}
}


*/
/*
func (t *table) enableTimeout(name string, enabled bool) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if ctrl, ok := t.controllers[name]; ok {
		c := cloneTimeout(ctrl.timeout)
		c.enabled = enabled
		t.update(name, cloneController[*timeout](ctrl, c))
		//ctrl.timeout.enabled = enabled
		//t.controllers[name] = ctrl
	}
}


*/
func (t *table) setTimeout(name string, duration time.Duration) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if ctrl, ok := t.controllers[name]; ok {
		c := cloneTimeout(ctrl.timeout)
		c.config.Duration = duration
		t.update(name, cloneController[*timeout](ctrl, c))
	}
}

/*
func (t *table) enableRateLimiter(name string, enabled bool) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if ctrl, ok := t.controllers[name]; ok {
		c := cloneRateLimiter(ctrl.rateLimiter)
		c.enabled = enabled
		t.update(name, cloneController[*rateLimiter](ctrl, c))
	}
}


*/
func (t *table) setRateLimit(name string, limit rate.Limit) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if ctrl, ok := t.controllers[name]; ok {
		c := cloneRateLimiter(ctrl.rateLimiter)
		c.config.Limit = limit
		// Not cloning the limiter as an old reference will not cause stale data when logging
		c.rateLimiter.SetLimit(limit)
		t.update(name, cloneController[*rateLimiter](ctrl, c))
	}
}

func (t *table) setRateBurst(name string, burst int) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if ctrl, ok := t.controllers[name]; ok {
		c := cloneRateLimiter(ctrl.rateLimiter)
		c.config.Burst = burst
		// Not cloning the limiter as an old reference will not cause stale data when logging
		c.rateLimiter.SetBurst(burst)
		t.update(name, cloneController[*rateLimiter](ctrl, c))
	}
}

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

func (t *table) enableRetry(name string, enabled bool) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if ctrl, ok := t.controllers[name]; ok {
		c := cloneRetry(ctrl.retry)
		if enabled {
			c.Enable()
		} else {
			c.Disable()
		}
		t.update(name, cloneController[*retry](ctrl, c))
	}
}

func (t *table) setRetryRateLimit(name string, limit rate.Limit, burst int) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if ctrl, ok := t.controllers[name]; ok {
		c := cloneRetry(ctrl.retry)
		c.config.Limit = limit
		c.config.Burst = burst
		// Not cloning the limiter as an old reference will not cause stale data when logging
		c.rateLimiter = rate.NewLimiter(limit, burst)
		t.update(name, cloneController[*retry](ctrl, c))
	}
}
