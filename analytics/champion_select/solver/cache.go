package solver

import (
	"errors"
	"sync"
)

// Cache stores payoffs that have already been calculated for starting state.
type Cache interface {
	// Put associates the payoff with the given normalized state.
	Put(state State, p Payoff) error

	// Get returns the payoff associated with the given normalized state. Returns
	// error if the state is not currently cached.
	Get(state State) (Payoff, error)
}

// mapCache implements Cache via in-memory map.
type mapCache struct {
	lock  sync.Mutex
	cache map[State]Payoff
}

// NewMapCache returns a Cache implemented via in-memory map.
func NewMapCache() Cache {
	return &mapCache{
		cache: make(map[State]Payoff),
	}
}

func (m *mapCache) Put(state State, p Payoff) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.cache[state] = p
	return nil
}

func (m *mapCache) Get(state State) (Payoff, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	got, ok := m.cache[state]
	if !ok {
		return Payoff{}, errors.New("key does not exist")
	}
	return got, nil
}
