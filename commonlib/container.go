package commonlib

import (
	"sync"
	"time"
)

// Container with Hashmap and Expiration
type Container struct {
	storage map[any]time.Time
	mu      sync.Mutex
	ttl     time.Duration
}

// Constructor
func NewContainer(ttl time.Duration) *Container {
	var c Container
	c.storage = make(map[any]time.Time)
	c.ttl = ttl
	return &c
}

// Add - add to container
func (c *Container) Add(value any) {
	c.mu.Lock()
	c.storage[value] = time.Now().Add(c.ttl)
	c.mu.Unlock()
}

// Delete - remove from container
func (c *Container) Delete(value any) {
	c.mu.Lock()
	delete(c.storage, value)
	c.mu.Unlock()
}

// CheckExpired - wipe all expired
func (c *Container) CheckExpired() {
	now := time.Now()
	c.mu.Lock()
	for key, timeOfExpire := range c.storage {
		if now.After(timeOfExpire) {
			delete(c.storage, key)
		}
	}
	c.mu.Unlock()
}

// Check - check value inside map
func (c *Container) CheckInside(value any) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.storage[value]
	return ok
}
