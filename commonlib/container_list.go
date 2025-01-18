package commonlib

import (
	"sync"
)

// Container with hashset
type Hashset struct {
	storage map[any]bool
	mu      sync.Mutex
}

// Constructor
func NewHashset() *Hashset {
	var c Hashset
	c.storage = make(map[any]bool)
	return &c
}

// Add - add to hashset
func (c *Hashset) Add(value any) {
	c.mu.Lock()
	c.storage[value] = true
	c.mu.Unlock()
}

// Delete - remove from hashset
func (c *Hashset) Delete(value any) {
	c.mu.Lock()
	delete(c.storage, value)
	c.mu.Unlock()
}

// Check - check value inside hashset
func (c *Hashset) CheckInside(value any) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.storage[value]
	return ok
}
