package lru

import "errors"

type Cache[K comparable, V any] struct {
	vals     map[K]V
	used     []K
	capacity int
}

func (c *Cache[K, V]) Set(k K, v V) {
	if len(c.vals) == c.capacity {
		delete(c.vals, c.used[0])
		c.used = c.used[1:]
	}
	c.vals[k] = v
	c.used = append(c.used, k)
}

func (c *Cache[K, V]) Get(k K) (V, bool) {
	v, ok := c.vals[k]
	return v, ok
}

func New[K comparable, V any](capacity int) (*Cache[K, V], error) {
	if capacity <= 1 {
		return nil, errors.New("cache capacity must be greater than 1")
	}
	return &Cache[K, V]{
		vals:     make(map[K]V),
		used:     make([]K, 0, capacity),
		capacity: capacity,
	}, nil
}
