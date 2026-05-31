package cache

import (
	"sync"
	"testing"
)

func TestLRUCache_PutGet(t *testing.T) {
	lru := NewLRUCache[string, int](1)
	lru.Put("test", 111)

	AssertEq(t, lru.Size(), 1, "size should be 1")

	v, ok := lru.Get("test")
	AssertTrue(t, ok, "get should return true")
	AssertEq(t, v, 111, "value should be 111")
}

func TestLRUCache_NoPutGet(t *testing.T) {
	lru := NewLRUCache[string, int](1)

	AssertEq(t, lru.Size(), 0, "size should be 0")
	_, ok := lru.Get("test")
	AssertFalse(t, ok, "key found unexpectedly")
}

func TestLRUCache_DoesntEvictIfHaveSpace(t *testing.T) {
	lru := NewLRUCache[string, int](2)
	lru.Put("test 1", 111)
	lru.Put("test 2", 222)

	v, ok := lru.Get("test 1")
	AssertTrue(t, ok, "key not found")
	AssertEq(t, v, 111, "value should be 111")

	v, ok = lru.Get("test 2")
	AssertTrue(t, ok, "key not found")
	AssertEq(t, v, 222, "value should be 222")

}

func TestLRUCache_EvictIfMaxSizeExceeds(t *testing.T) {
	lru := NewLRUCache[string, int](2)
	lru.Put("test 1", 111)
	lru.Put("test 2", 222)
	lru.Put("test 3", 333)

	v, ok := lru.Get("test 2")
	AssertTrue(t, ok, "key not found")
	AssertEq(t, v, 222, "value should be 222")

	v, ok = lru.Get("test 3")
	AssertTrue(t, ok, "key not found")
	AssertEq(t, v, 333, "value should be 333")

	v, ok = lru.Get("test 1")
	AssertFalse(t, ok, "key not evicted")
}

func TestLRUCache_GetMovesToFront(t *testing.T) {
	lru := NewLRUCache[string, int](2)
	lru.Put("test 1", 111)
	lru.Put("test 2", 222)
	_, _ = lru.Get("test 1")
	lru.Put("test 3", 333)

	v, ok := lru.Get("test 1")
	AssertTrue(t, ok, "key not found")
	AssertEq(t, v, 111, "value should be 111")

	v, ok = lru.Get("test 3")
	AssertTrue(t, ok, "key not found")
	AssertEq(t, v, 333, "value should be 333")

	v, ok = lru.Get("test 2")
	AssertFalse(t, ok, "key not evicted")
}

func TestLRUCache_UpdateExistingKey(t *testing.T) {
	lru := NewLRUCache[string, int](1)
	lru.Put("test 1", 111)
	lru.Put("test 1", 1111)

	v, ok := lru.Get("test 1")
	AssertTrue(t, ok, "key not found")
	AssertEq(t, v, 1111, "value should be updated")
}

func TestLRUCache_Concurrent(t *testing.T) {
	const N = 10000
	lru := NewLRUCache[int, int](N)

	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			lru.Put(i, i)
			v, ok := lru.Get(i)
			AssertTrue(t, ok, "key not found")
			AssertEq(t, v, i, "invalid value")
		}(i)
	}

	wg.Wait()
	for i := 0; i < N; i++ {
		v, ok := lru.Get(i)
		AssertTrue(t, ok, "key not found")
		AssertEq(t, v, i, "invalid value")
	}
}

func AssertTrue(t *testing.T, v bool, message string) {
	if !v {
		t.Fatal(message)
	}
}

func AssertFalse(t *testing.T, v bool, message string) {
	if v {
		t.Fatal(message)
	}
}

func AssertEq[T comparable](t *testing.T, v1 T, v2 T, message string) {
	if v1 != v2 {
		t.Fatal(message)
	}
}
