// Package ratelimit implements rate limiting for the public Riot API.
//
// This package defines client-side rate limiting. For centralized rate
// limiting, see the service sub-package.
package ratelimit

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const sleepBeforeRetryAcquire = 25 * time.Millisecond

// Invocation represents a specific application's invocation of the Riot API.
type Invocation struct {
	// ApplicationKey is any unique application identifier, typically the Riot
	// API key provided by the Riot developer portal.
	ApplicationKey string

	// Region is the region for which the method is called. Limits are enforced
	// on a by-region basis.
	Region string

	// Method is the relative method path with all options stripped. For example,
	// a valid method is "/lol/match/v4/matches".
	Method string

	// Uniquifier is an optional token that helps make the Invocation unique. In
	// almost all instances, the Method is sufficient as a uniquifier, so this can
	// be left as the empty string. However, for the match API, getting
	// matchlists by account and getting recent matchlists by account have the same
	// underlying method with different path arguments. The Uniquifier field
	// allows the methods to be considered as separate invocations.
	Uniquifier string

	// NoAppQuota is true if the invocation doesn't take up appliation quota. The
	// default false value is typical for most invocations, which do in fact use
	// app quota.
	NoAppQuota bool
}

// App returns an invocation that is application-level as opposed to
// method-level. It is used to track global quota.
func (i Invocation) App() Invocation {
	return Invocation{
		ApplicationKey: i.ApplicationKey,
		Region:         i.Region,
	}
}

// Done is a callback returned by Acquire() that signals the end of an API
// method call. Calling Done will schedule the rate to be added back to the
// pool at the appropriate time.
//
// If the given response is non-nil, then Done() will parse the response
// headers for rate-limiting information and update configured limits. If the
// response indicates a rate violation, then Done will require the indicated
// sleep time until resources can be reserved again.
type Done func(res *http.Response) error

// Cancel is a callback returned by Acquire() that signals the immediate
// release of rate resources. This function must only be called when no riot
// API method was invoked following Acquire().
type Cancel func() error

// Limiter brokers access to rate resources.
type Limiter interface {
	// Acquire blocks until all configured limits for the invocation are
	// satisfied, or until the context is cancelled. Once acquired, the rate
	// resource is reserved until Done() or Cancel() are called and return nil.
	Acquire(ctx context.Context, inv Invocation) (Done, Cancel, error)
}

// invocationLimit represents a rate limit for a specific type of invocation.
type invocationLimit struct {
	// limits maps interval length in seconds to the *singleLimit.
	limits sync.Map
}

// Get returns the singleLimit for the interval in seconds, or nil if no limit is
// configured for that time interval.
func (i *invocationLimit) Get(ts int64) *singleLimit {
	obj, ok := i.limits.Load(ts)
	if !ok {
		return nil
	}
	return obj.(*singleLimit)
}

// ForEachLimit applies the function f to each limit associated with this
// invocation.
func (i *invocationLimit) ForEachLimit(f func(seconds int64, limit *singleLimit) (next bool)) {
	i.limits.Range(func(key, value interface{}) bool {
		return f(key.(int64), value.(*singleLimit))
	})
}

// SetLimitCapacity either modifies the stored limit or creates one with the
// given capacity.
func (i *invocationLimit) SetLimitCapacity(seconds, capacity int64) {
	o, _ := i.limits.LoadOrStore(seconds, &singleLimit{
		capacity: capacity,
		quantity: capacity,
	})
	o.(*singleLimit).SetCapacity(capacity)
}

// NewLimiter returns an in-proecss limiter.
func NewLimiter() Limiter {
	return &limiter{
		methodWake: make(map[Invocation]time.Time),
	}
}

type limiter struct {
	// limits maps from Invocation to an *invocationLimit. The Invocation with
	// empty Method field corresponds to the application-level limits.
	limits sync.Map

	// lock protects methodWake. The empty method corresponds to application
	// limits. Service limits are also included as application limits, since they
	// have the same underlying effect.
	lock       sync.RWMutex
	methodWake map[Invocation]time.Time
}

// getInvocationLimit returns the limit corresponding to the given invocation.
// If no limit is configured, then this returns nil.
func (l *limiter) getInvocationLimit(inv Invocation) *invocationLimit {
	got, ok := l.limits.Load(inv)
	if !ok {
		return nil
	}
	return got.(*invocationLimit)
}

// getOrCreateInvocationLimit returns the limit corresponding to the given
// invocation. If it does not yet exist, then create one and return it.
func (l *limiter) getOrCreateInvocationLimit(inv Invocation) *invocationLimit {
	val, _ := l.limits.LoadOrStore(inv, &invocationLimit{})
	return val.(*invocationLimit)
}

// maybeSleep returns after either all required sleeps for the invocation are
// complete, or the context is cancelled.
func (l *limiter) maybeSleep(ctx context.Context, inv Invocation) error {
	l.lock.RLock()
	appWake := l.methodWake[inv.App()]
	methodWake := l.methodWake[inv]
	l.lock.RUnlock()

	if !appWake.IsZero() {
		select {
		case <-time.NewTimer(time.Until(appWake)).C:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	if !methodWake.IsZero() {
		select {
		case <-time.NewTimer(time.Until(methodWake)).C:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

// setCapacityForInvocation takes an HTTP header containing rate capacities and
// stores these capacitites in the limits structure corresponding to the given
// invocation.
func (l *limiter) setCapacityForInvocation(header string, inv Invocation) error {
	limits, err := headerIntMap(header)
	if err != nil {
		return err
	}
	if len(limits) != 0 {
		il := l.getOrCreateInvocationLimit(inv)
		for seconds, capacity := range limits {
			il.SetLimitCapacity(seconds, capacity)
		}
	}
	return nil
}

// matchRiotCounts parses the header containing counts, and reconciles them to
// the invocationLimit corresponding to the given Invocation.
func (l *limiter) matchRiotCounts(header string, inv Invocation) error {
	counts, err := headerIntMap(header)
	if err != nil {
		return err
	}
	if len(counts) != 0 {
		il := l.getOrCreateInvocationLimit(inv)
		for seconds, q := range counts {
			got := il.Get(seconds)
			if got != nil {
				got.MatchRiotCounts(q, time.Duration(seconds)*time.Second)
			}
		}
	}
	return nil
}

// cancelAllAcquired cancels all acquired limits in the given map.
func cancelAllAcquired(acquired map[int64]*singleLimit) {
	for _, lim := range acquired {
		lim.Cancel()
	}
}

// acquireAllOrCancel acquires all time interval quota for the given
// invocation, returning the acquired limits and true on success. If any
// interval quota cannot be acquired, then return nil and false.
func (l *limiter) acquireAllOrCancel(inv Invocation) (map[int64]*singleLimit, bool) {
	limits := l.getInvocationLimit(inv)
	if limits == nil {
		return nil, true
	}

	acquired := make(map[int64]*singleLimit)
	allAcquired := true
	limits.ForEachLimit(func(seconds int64, lim *singleLimit) bool {
		if lim.Acquire() {
			acquired[seconds] = lim
			return true
		}
		allAcquired = false
		return false
	})

	if allAcquired {
		return acquired, true
	}

	cancelAllAcquired(acquired)

	return nil, false
}

// Acquire blocks until all configured limits for the invocation are satisfied,
// or until the context is cancelled. Once acquired, the rate resource is
// reserved until Done() or Cancel() are called and return nil.
func (l *limiter) Acquire(ctx context.Context, inv Invocation) (Done, Cancel, error) {
	err := l.maybeSleep(ctx, inv)
	if err != nil {
		return nil, nil, err
	}

	var (
		appAcquired, acquired       map[int64]*singleLimit
		appAllAcquired, allAcquired bool
	)

	for {
		if inv.NoAppQuota {
			appAllAcquired = true
		} else {
			appAcquired, appAllAcquired = l.acquireAllOrCancel(inv.App())
		}

		if appAllAcquired {
			acquired, allAcquired = l.acquireAllOrCancel(inv)
			if allAcquired {
				break
			}
			cancelAllAcquired(appAcquired)
		}
		// Sleep before retrying, up until cancellation.
		select {
		case <-time.NewTimer(sleepBeforeRetryAcquire).C:
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		}
	}

	var refundOnce, cancelOnce sync.Once

	done := func(res *http.Response) error {
		refundOnce.Do(func() {
			for seconds, lim := range appAcquired {
				lim.AddQuantity(1, time.Duration(seconds)*time.Second)
			}
			for seconds, lim := range acquired {
				lim.AddQuantity(1, time.Duration(seconds)*time.Second)
			}
		})

		if res != nil {
			appLimit := strings.TrimSpace(res.Header.Get("X-App-Rate-Limit"))
			appCount := strings.TrimSpace(res.Header.Get("X-App-Rate-Limit-Count"))
			methodLimit := strings.TrimSpace(res.Header.Get("X-Method-Rate-Limit"))
			methodCount := strings.TrimSpace(res.Header.Get("X-Method-Rate-Limit-Count"))
			retryAfter := strings.TrimSpace(res.Header.Get("Retry-After"))
			retryType := strings.TrimSpace(res.Header.Get("X-Rate-Limit-Type"))

			appKey := inv.App()

			if appLimit != "" {
				err = l.setCapacityForInvocation(appLimit, appKey)
				if err != nil {
					return err
				}
				err = l.matchRiotCounts(appCount, appKey)
				if err != nil {
					return err
				}
			}
			if methodLimit != "" {
				err = l.setCapacityForInvocation(methodLimit, inv)
				if err != nil {
					return err
				}
				err = l.matchRiotCounts(methodCount, inv)
				if err != nil {
					return err
				}
			}
			if retryAfter != "" {
				retrySeconds, err := strconv.ParseInt(retryAfter, 10, 64)
				if err != nil {
					return err
				}
				until := time.Now().Add(time.Duration(retrySeconds) * time.Second)
				var sleepKey Invocation
				// Method sleeps are tied to this specific invocation.
				if retryType == "method" {
					sleepKey = inv
				} else {
					sleepKey = inv.App()
				}
				l.lock.Lock()
				if until.After(l.methodWake[sleepKey]) {
					l.methodWake[sleepKey] = until
				}
				l.lock.Unlock()
			}
		}
		return nil
	}

	cancel := func() error {
		cancelOnce.Do(func() {
			cancelAllAcquired(appAcquired)
			cancelAllAcquired(acquired)
		})
		return nil
	}

	return done, cancel, nil
}

// headerIntMap takes a string representing rates for seconds intervals like
// "100:20,200:30" and returns the value-key map {20: 100, 30: 200}. Returns
// error if the format is incorrect.
func headerIntMap(header string) (map[int64]int64, error) {
	secondToCount := make(map[int64]int64)
	pieces := strings.Split(header, ",")
	for _, piece := range pieces {
		kv := strings.Split(piece, ":")
		if len(kv) != 2 {
			return nil, fmt.Errorf("expected K:V in %q", header)
		}
		count, err := strconv.ParseInt(kv[0], 10, 64)
		if err != nil {
			return nil, err
		}
		seconds, err := strconv.ParseInt(kv[1], 10, 64)
		if err != nil {
			return nil, err
		}
		secondToCount[seconds] = count
	}
	return secondToCount, nil
}
