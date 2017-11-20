package ratelimit

import (
	"sync"
	"time"
)

// singleLimit is a rate limit corresponding to a specific time interval.
type singleLimit struct {
	lock sync.Mutex

	// riotMatcher is a timer that fires a reconciliation event. For example, if
	// Riot believes we have used 100 and we believe we have used 90, then we
	// will immediately up our estimate us usage by 10, and then decrease by 10
	// after some time has elapsed. Only one reconciliation event is needed at
	// any given point, since a future reconciliation should immediately replace
	// the existing one.
	riotMatcher *time.Timer

	// riotOffset is the quantity that will be added back following the
	// expiration of riotMatcher.
	riotOffset int64

	capacity int64
	quantity int64
}

// Acquire attempts to reserve one unit and returns true on success.
func (s *singleLimit) Acquire() (ok bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.quantity <= 0 {
		return false
	}
	s.quantity--
	return true
}

// Cancel adds one to the available quantity. This must only be called
// following a successful Acquire(), and is intended to be used to signify that
// an acquired resource was not used. The function does not check whether this
// is the case.
func (s *singleLimit) Cancel() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.quantity++
}

// AddQuantity adds or subtracts the resource after the given duration.
func (s *singleLimit) AddQuantity(q int64, d time.Duration) {
	time.AfterFunc(d, func() {
		s.lock.Lock()
		defer s.lock.Unlock()
		s.quantity += q
	})
}

// SetCapacity sets the new limit capacity. The available resources are shrunk
// to the new capacity if it is smaller. The resources are increased by the
// capacity increase if the new capacity is larger.
func (s *singleLimit) SetCapacity(c int64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	old := s.capacity
	s.capacity = c
	if c > old {
		s.quantity += (c - old)
	}
	if s.quantity > c {
		s.quantity = c
	}
}

// MatchRiotCounts reconciles the currently tracked quantity to the given
// counts from Riot. If the Riot counts are higher than the implied counts we
// are tracking, we will immediately lower the current quantity by the
// difference, and add an eventual offest after the given duration.
func (s *singleLimit) MatchRiotCounts(counts int64, reverseAfter time.Duration) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.riotMatcher != nil {
		// If the reconciliation already happened, there's no need for special
		// logic, since the reconciliation we are about to schedule can be added
		// right on top. If the reconciliation was pending, then we need to reverse
		// it first.
		if stopped := s.riotMatcher.Stop(); stopped {
			s.quantity += s.riotOffset
		}
	}

	// Example: capacity of 10, Riot shows 4 counts => 6 implied quantity
	impliedQuantity := s.capacity - counts

	if impliedQuantity != s.quantity {
		s.riotOffset = s.quantity - impliedQuantity
		s.quantity -= s.riotOffset
		s.riotMatcher = time.AfterFunc(reverseAfter, func() {
			s.lock.Lock()
			defer s.lock.Unlock()
			s.quantity += s.riotOffset
		})
	}
}
