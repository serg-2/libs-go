package commonlib

import (
	"sync"
)

// Container with hashset
type Hashset struct {
	InternalHashset map[any]bool
	Mut             sync.Mutex
}

// Constructor
func NewHashset() *Hashset {
	var c Hashset
	c.InternalHashset = make(map[any]bool)
	return &c
}

// Add - add to hashset
func (c *Hashset) Add(value any) {
	c.Mut.Lock()
	c.InternalHashset[value] = true
	c.Mut.Unlock()
}

// Delete - remove from hashset
func (c *Hashset) Delete(value any) {
	c.Mut.Lock()
	delete(c.InternalHashset, value)
	c.Mut.Unlock()
}

// Check - check value inside hashset
func (c *Hashset) CheckInside(value any) bool {
	c.Mut.Lock()
	defer c.Mut.Unlock()
	_, ok := c.InternalHashset[value]
	return ok
}

// GetAll - GetAll
func (c *Hashset) GetAll(value any) []any {
	c.Mut.Lock()
	defer c.Mut.Unlock()
	var result []any
	for val := range c.InternalHashset {
		result = append(result, val)
	}
	return result
}
