package commonlib

import (
	"sync"
	"time"
)

type Container struct {
	Storage map[any]time.Time
	mu      sync.Mutex
}

// Constructor
func NewContainer() *Container {
	var c Container
	c.Storage = make(map[any]time.Time)
	return &c
}

// Add - add to container
func (c *Container) Add(value any) {
	c.mu.Lock()
	c.Storage[value] = time.Now()
	c.mu.Unlock()
}

// Delete - remove from container
func (c *Container) Delete(value any) {
	c.mu.Lock()
	delete(c.Storage, value)
	c.mu.Unlock()
}

// Check - check value after timestamp
func (c *Container) CheckExpired(value any, timeOfExpiration time.Time) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	timeToCheck, ok := c.Storage[value]
	if !ok {
		return false
	}
	if timeToCheck.After(timeOfExpiration) {
		return true
	}
	return false
}

// Check - check value inside map
func (c *Container) CheckInside(value any) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.Storage[value]
	return ok
}
