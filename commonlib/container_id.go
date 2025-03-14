package commonlib

import (
	"sync"
)

// Container with Id
type ContainerId struct {
	InternalContainer map[string]any
	Mut               sync.Mutex
}

// Constructor
func NewContainerId() *ContainerId {
	var c ContainerId
	c.InternalContainer = make(map[string]any)
	return &c
}

// Add - add to container
func (c *ContainerId) Add(key string, value any) {
	c.Mut.Lock()
	defer c.Mut.Unlock()
	c.InternalContainer[key] = value
}

// Get - get value by key inside container
func (c *ContainerId) Get(key string) any {
	c.Mut.Lock()
	defer c.Mut.Unlock()
	res, ok := c.InternalContainer[key]
	if !ok {
		return nil
	}
	return res
}

// GetOrDefault - get value by key inside container or get default value
func (c *ContainerId) GetOrDefault(key string, defaultValue any) any {
	c.Mut.Lock()
	defer c.Mut.Unlock()
	res, ok := c.InternalContainer[key]
	if !ok {
		return defaultValue
	}
	return res
}

// Delete - remove from container
func (c *ContainerId) Delete(key string) {
	c.Mut.Lock()
	defer c.Mut.Unlock()
	delete(c.InternalContainer, key)
}

// Size - size of container
func (c *ContainerId) Size() int {
	c.Mut.Lock()
	defer c.Mut.Unlock()
	return len(c.InternalContainer)
}

// Clear - clear container
func (c *ContainerId) Clear() {
	c.Mut.Lock()
	defer c.Mut.Unlock()
	c.InternalContainer = make(map[string]any)
}

// Check - check key inside container
func (c *ContainerId) CheckInside(key string) bool {
	c.Mut.Lock()
	defer c.Mut.Unlock()
	_, ok := c.InternalContainer[key]
	return ok
}

// GetAllValues - GetAllValues
func (c *ContainerId) GetAllValues() []any {
	c.Mut.Lock()
	defer c.Mut.Unlock()
	var result []any
	for _, val := range c.InternalContainer {
		result = append(result, val)
	}
	return result
}

// GetAllKeys - GetAllKeys
func (c *ContainerId) GetAllKeys() []string {
	c.Mut.Lock()
	defer c.Mut.Unlock()
	var result []string
	for key := range c.InternalContainer {
		result = append(result, key)
	}
	return result
}
