package main

import (
	"math/rand"
	"sync"
)

// LossOfInterest represents the strategy that a Node uses to decide
// the status of an Update
// A LossOfInterest instance should be thread safe
type LossOfInterest interface {
	// Status returns the status of an update
	Status(Update) Status
	// Feedback updates the internal status of the strategy
	// given the update and the feedback received from the other peer
	Feedback(Update, bool)
}

func NewCounterFeedback(K uint) LossOfInterest {
	return &CounterFeedback{
		K: K,

		seen: make(map[Update]uint),
	}
}

// CounterFeedback implements the counter and feedback strategy
type CounterFeedback struct {
	// K is the maximum number of times a node can send an update to peers which already have it.
	// If the number is greater the update is considered removed.
	K uint

	seen map[Update]uint
	mu   sync.RWMutex
}

func (c *CounterFeedback) Status(u Update) Status {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if val, ok := c.seen[u]; ok {
		if val >= c.K {
			return Removed
		}
	}
	return Infective
}

func (c *CounterFeedback) Feedback(u Update, feedback bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.seen[u]; !ok {
		c.seen[u] = 0
	}
	if feedback {
		c.seen[u]++
	}
}

func NewBlindRandom(K uint) LossOfInterest {
	return &BlindRandom{
		K: K,

		removed: make(map[Update]bool),
	}
}

// BlindRandom implements the blind and random strategy
type BlindRandom struct {
	// 1/K is the probability that an update will be considered removed.
	K uint

	removed map[Update]bool
	mu      sync.RWMutex
}

func (b *BlindRandom) Status(u Update) Status {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if removed, ok := b.removed[u]; ok {
		if removed {
			return Removed
		}
	}
	return Infective
}

func (b *BlindRandom) Feedback(u Update, _ bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.removed[u]; !ok {
		b.removed[u] = false
	}

	// dice toss
	threshold := 1 / float32(b.K)
	result := rand.Float32() < threshold

	b.removed[u] = result
}
