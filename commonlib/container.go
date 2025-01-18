package commonlib

import (
	"sync"
	"time"
)

// Container with Hashmap and Expiration
type Container struct {
	InternalMap map[any]time.Time
	Mu          sync.Mutex
	Ttl         time.Duration
}

// Constructor
func NewContainer(ttl time.Duration) *Container {
	var c Container
	c.InternalMap = make(map[any]time.Time)
	c.Ttl = ttl
	return &c
}

// Add - add to container
func (c *Container) Add(value any) {
	c.Mu.Lock()
	c.InternalMap[value] = time.Now().Add(c.Ttl)
	c.Mu.Unlock()
}

// Delete - remove from container
func (c *Container) Delete(value any) {
	c.Mu.Lock()
	delete(c.InternalMap, value)
	c.Mu.Unlock()
}

// CheckExpired - wipe all expired
func (c *Container) CheckExpired() {
	now := time.Now()
	c.Mu.Lock()
	for key, timeOfExpire := range c.InternalMap {
		if now.After(timeOfExpire) {
			delete(c.InternalMap, key)
		}
	}
	c.Mu.Unlock()
}

// Check - check value inside map
func (c *Container) CheckInside(value any) bool {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	_, ok := c.InternalMap[value]
	return ok
}
