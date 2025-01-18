package commonlib

import (
	"sync"
	"time"
)

type Container struct {
	Storage map[any]time.Time
	mu      sync.Mutex
	ttl     time.Duration
}

// Constructor
func NewContainer(ttl time.Duration) *Container {
	var c Container
	c.Storage = make(map[any]time.Time)
	c.ttl = ttl
	return &c
}

// Add - add to container
func (c *Container) Add(value any) {
	c.mu.Lock()
	c.Storage[value] = time.Now().Add(c.ttl)
	c.mu.Unlock()
}

// Delete - remove from container
func (c *Container) Delete(value any) {
	c.mu.Lock()
	delete(c.Storage, value)
	c.mu.Unlock()
}

// CheckExpired - wipe all expired
func (c *Container) CheckExpired() {
	now := time.Now()
	c.mu.Lock()
	for key, timeOfExpire := range c.Storage {
		if now.After(timeOfExpire) {
			delete(c.Storage, key)
		}
	}
	c.mu.Unlock()
}

// Check - check value inside map
func (c *Container) CheckInside(value any) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.Storage[value]
	return ok
}
