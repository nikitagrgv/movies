package cache

import "sync"

type LRUCache[K comparable, V any] struct {
	mu      sync.RWMutex
	m       map[K]*node[K, V]
	size    int
	maxSize int
	back    *node[K, V]
	front   *node[K, V]
}

type node[K comparable, V any] struct {
	prev  *node[K, V]
	next  *node[K, V]
	key   K
	value V
}

func NewLRUCache[K comparable, V any](maxSize int) *LRUCache[K, V] {
	if maxSize <= 0 {
		panic("maxSize must be greater than zero")
	}
	return &LRUCache[K, V]{
		m:       make(map[K]*node[K, V]),
		maxSize: maxSize,
	}
}

func (c *LRUCache[K, V]) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.size
}

func (c *LRUCache[K, V]) MaxSize() int {
	return c.maxSize
}

func (c *LRUCache[K, V]) Put(k K, v V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	n, ok := c.m[k]
	if ok {
		n.value = v
		c.moveToFront(n)
		return
	}

	if c.size >= c.maxSize {
		// reuse back node
		n = c.back
		delete(c.m, n.key)
		c.moveToFront(n)
		c.m[k] = n
		n.key = k
		n.value = v
		return
	}

	n = &node[K, V]{
		key:   k,
		value: v,
	}

	if c.size == 0 {
		c.back = n
		c.front = n
	} else {
		c.moveToFront(n)
	}
	c.m[k] = n
	c.size++
}

func (c *LRUCache[K, V]) Get(k K) (v V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	n, ok := c.m[k]
	if !ok {
		var empty V
		return empty, false
	}

	c.moveToFront(n)
	return n.value, true
}

func (c *LRUCache[K, V]) moveToFront(n *node[K, V]) {
	if n == c.front {
		return
	}

	if n.prev != nil {
		n.prev.next = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	}

	if c.back == n {
		c.back = n.prev
	}

	if c.front != nil {
		c.front.prev = n
	}

	n.next = c.front
	n.prev = nil
	c.front = n
}
