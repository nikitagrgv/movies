package cache

import (
	"sync"
	"testing"
)

func TestLRUCache_PutAndGetWorks(t *testing.T) {
	lru := NewLRUCache[string, int](1)
	lru.Put("test", 111)

	if lru.Size() != 1 {
		t.Fatal("size should be 1")
	}
	AssertExists(t, lru, "test", 111)
}

func TestLRUCache_GetNonExisting(t *testing.T) {
	lru := NewLRUCache[string, int](1)
	if lru.Size() != 0 {
		t.Fatal("size should be 0")
	}
	AssertNotExists(t, lru, "test")
}

func TestLRUCache_DoesntEvictIfHaveSpace(t *testing.T) {
	lru := NewLRUCache[string, int](2)
	lru.Put("test 1", 111)
	lru.Put("test 2", 222)

	AssertExists(t, lru, "test 1", 111)
	AssertExists(t, lru, "test 2", 222)
}

func TestLRUCache_EvictIfMaxSizeExceeds(t *testing.T) {
	lru := NewLRUCache[string, int](2)
	lru.Put("test 1", 111)
	lru.Put("test 2", 222)
	lru.Put("test 3", 333)

	AssertNotExists(t, lru, "test 1")
	AssertExists(t, lru, "test 2", 222)
	AssertExists(t, lru, "test 3", 333)
}

func TestLRUCache_EvictNormallyIfSize1(t *testing.T) {
	lru := NewLRUCache[string, int](1)
	lru.Put("test 1", 111)
	lru.Put("test 2", 222)

	AssertNotExists(t, lru, "test 1")
	AssertExists(t, lru, "test 2", 222)
}

func TestLRUCache_GetMovesToFrontLastUsed(t *testing.T) {
	lru := NewLRUCache[string, int](3)
	lru.Put("test 1", 111)
	lru.Put("test 2", 222)
	lru.Put("test 3", 333)
	_, _ = lru.Get("test 1")
	lru.Put("test 4", 444)

	AssertExists(t, lru, "test 1", 111)
	AssertNotExists(t, lru, "test 2")
	AssertExists(t, lru, "test 3", 333)
	AssertExists(t, lru, "test 4", 444)
}

func TestLRUCache_GetMovesToFrontMiddleUsed(t *testing.T) {
	lru := NewLRUCache[string, int](3)
	lru.Put("test 1", 111)
	lru.Put("test 2", 222)
	lru.Put("test 3", 333)
	_, _ = lru.Get("test 2")
	lru.Put("test 4", 444)

	AssertNotExists(t, lru, "test 1")
	AssertExists(t, lru, "test 2", 222)
	AssertExists(t, lru, "test 3", 333)
	AssertExists(t, lru, "test 4", 444)
}

func TestLRUCache_UpdateExistingKey(t *testing.T) {
	lru := NewLRUCache[string, int](1)
	lru.Put("test 1", 111)
	lru.Put("test 1", 1111)

	AssertExists(t, lru, "test 1", 1111)
}

func TestLRUCache_UpdateExistingKeyDoesntChangeSize(t *testing.T) {
	lru := NewLRUCache[string, int](1)
	lru.Put("test 1", 111)
	lru.Put("test 1", 1111)

	if lru.Size() != 1 {
		t.Fatal("size should be 1")
	}
}

func TestLRUCache_PanicIfSizeZero(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("should have panicked")
		}
	}()
	_ = NewLRUCache[string, int](0)
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
			_, _ = lru.Get(i)
		}(i)
	}

	wg.Wait()
	for i := 0; i < N; i++ {
		v, ok := lru.Get(i)
		if !ok {
			t.Fatal("not found")
		}
		if v != i {
			t.Fatal("invalid value")
		}
	}
}

func AssertExists[K, V comparable](t *testing.T, lru *LRUCache[K, V], key K, value V) {
	v, ok := lru.Get(key)
	if !ok {
		t.Fatal("not found: ", key)
	}
	if v != value {
		t.Fatalf("expected %v, got %v", value, v)
	}
}

func AssertNotExists[K comparable, V any](t *testing.T, lru *LRUCache[K, V], key K) {
	_, ok := lru.Get(key)
	if ok {
		t.Fatal("mustn't exist: ", key)
	}
}
