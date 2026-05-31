package cache

import "sync"

type LRUCache[K comparable, V any] struct {
	mu      sync.RWMutex
	m       map[K]*node[K, V]
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
	return len(c.m)
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
		c.moveFront(n)
		return
	}

	if c.Size() >= c.maxSize {
		// reuse back node
		n = c.back
		delete(c.m, n.key)
		c.m[k] = n
		n.key = k
		n.value = v
		c.removeNode(n)
		c.pushFront(n)
		return
	}

	n = &node[K, V]{
		key:   k,
		value: v,
	}

	c.pushFront(n)
	c.m[k] = n
}

func (c *LRUCache[K, V]) Get(k K) (v V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	n, ok := c.m[k]
	if !ok {
		var empty V
		return empty, false
	}

	c.moveFront(n)
	return n.value, true
}

func (c *LRUCache[K, V]) removeNode(n *node[K, V]) {
	if c.front == n {
		c.front = n.next
	}
	if c.back == n {
		c.back = n.prev
	}

	if n.prev != nil {
		n.prev.next = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	}

	n.prev = nil
	n.next = nil
}

func (c *LRUCache[K, V]) pushFront(n *node[K, V]) {
	if c.front != nil {
		c.front.prev = n
	}
	c.front = n
	if c.back == nil {
		c.back = n
	}
}

func (c *LRUCache[K, V]) moveFront(n *node[K, V]) {
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
