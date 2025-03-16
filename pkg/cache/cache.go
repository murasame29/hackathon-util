package cache

import (
	"errors"
	"sync"
)

var (
	ErrNotFound = errors.New("not found")
)

type Cache struct {
	value sync.Map
}

func New() *Cache {
	return &Cache{
		value: sync.Map{},
	}
}

func (c *Cache) Get(key string) (any, bool) {
	return c.value.Load(key)
}

func (c *Cache) Set(key string, value any) {
	c.value.Store(key, value)
}

func (c *Cache) Delete(key string) {
	c.value.Delete(key)
}
